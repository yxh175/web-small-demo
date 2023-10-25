package main

import (
	"context"
	"errors"
	"gin-mall/soldout_demo/db/dao"
	"gin-mall/soldout_demo/db/model"
	"gin-mall/soldout_demo/redis_lock"
	"log"
	"net/http"
	"strconv"
	"time"

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

	// 消费队列
	orderList := "orderList"

	orderDb := dao.NewOrderDao(ctx)
	// 订阅 Stream
	for {
		str, err := client.LIndex(ctx, orderList, -1).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			// 处理错误
			log.Fatal(err)
			return
		}

		if errors.Is(err, redis.Nil) {
			time.Sleep(time.Second)
			continue
		}
		id, err := strconv.Atoi(str)
		if err != nil {
			log.Fatal(err)
			return
		}
		err = orderDb.AddOrder(&model.Order{
			Uid: id,
		})
		if err != nil {
			log.Fatal(err)
			return
		}
		if err := client.RPop(ctx, orderList).Err(); err != nil {
			log.Fatal(err)
			return
		}
	}
}
