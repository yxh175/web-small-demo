package main

import (
	"context"
	"fmt"
	"gin-mall/soldout_demo/db/dao"
	"gin-mall/soldout_demo/db/model"
	"gin-mall/soldout_demo/redis_lock"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func init() {
	dao.InitMySQL()
	go getOrder()
}

func main() {
	http.HandleFunc("/", addOrder)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func addOrder(w http.ResponseWriter, r *http.Request) {
	rl := redis_lock.NewRedisLock()

	defer r.Body.Close()

	// r.ParseForm()
	// uid := r.FormValue("uid")
	uid := "123"

	rl.EvalScript(uid)
}

func getOrder() {
	var ctx = context.Background()
	var client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 如果有密码
		DB:       0,
	})

	// 消费者组名称
	subject := "my_orders_stream"
	consumerGroup := "my_orders_consumer_group"

	// 创建stream
	err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: subject,
		ID:     "*",
		Values: map[string]interface{}{
			"user_id": "hh",
		},
	}).Err()
	if err != nil {
		fmt.Printf("创建stream异常: %v\n", err)
		return
	}

	// 创建 Redis Stream 操作器
	err = client.XGroupCreate(ctx, subject, consumerGroup, "$").Err()
	if err != nil {
		fmt.Printf("创建消费者组出错: %v\n", err)
		return
	}
	orderDb := dao.NewOrderDao(ctx)
	// 订阅 Stream
	for {
		messages, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumerGroup,
			Consumer: "my_consumer",
			Streams:  []string{subject, ">"},
			Count:    1,
			Block:    0,
			NoAck:    false,
		}).Result()
		if err != nil {
			// 处理错误
			log.Fatal(err)
		}

		for _, message := range messages {
			fmt.Println("获取到信息", message)
			for _, xMessage := range message.Messages {
				userId := xMessage.Values["user_id"]
				if userId == nil {
					fmt.Println("无效消息")
					break
				}
				order := &model.Order{
					Name: userId.(string),
				}
				err = orderDb.AddOrder(order)
				if err != nil {
					log.Fatal(err)
					break
				}
				_, err = client.XAck(ctx, "my_orders_stream", consumerGroup, xMessage.ID).Result()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
