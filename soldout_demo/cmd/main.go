package main

import (
	"context"
	"fmt"
	"gin-mall/soldout_demo/db/dao"
	"gin-mall/soldout_demo/db/model"
	"gin-mall/soldout_demo/redis_lock"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func init() {
	dao.InitMySQL()
	getOrder()
}

func main() {
	http.HandleFunc("/", addOrder)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func addOrder(w http.ResponseWriter, r *http.Request) {
	rl := redis_lock.NewRedisLock()

	defer r.Body.Close()

	r.ParseForm()
	uid := r.FormValue("uid")

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

	// 创建 Redis Stream 操作器
	err := client.XGroupCreate(ctx, subject, consumerGroup, "$")
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
			for _, xMessage := range message.Messages {
				userId := xMessage.Values["user_id"]
				order := &model.Order{
					Gid: userId.(int),
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
