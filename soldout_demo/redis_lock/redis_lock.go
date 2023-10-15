package redis_lock

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var rL *RedisLock
var once sync.Once

const orderSet = "orderSet"     //用户id的集合
const goodsTotal = "goodsTotal" //商品库存的key
const orderList = "orderList"   //订单队列

type RedisLock struct {
	client *redis.Client
}

func NewRedisLock() *RedisLock {
	once.Do(func() {
		rL = &RedisLock{
			client: redis.NewClient(&redis.Options{
				Addr:     "localhost:6379", // Redis 服务器地址
				Password: "",               // 如果有密码，设置密码
				DB:       0,                // 默认数据库
				PoolSize: 1000,
			}),
		}
	})
	return rL
}

// CreateScript 加载脚本
func (rl *RedisLock) CreateScript() (*redis.Script, error) {
	str, err := os.ReadFile("./soldout_demo/lua-case/luaScript.lua")
	if err != nil {
		return nil, err
	}
	scriptStr := string(str)
	script := redis.NewScript(scriptStr)
	return script, nil
}

// EvalScript 执行脚本
func (rl *RedisLock) EvalScript(userId string) {
	ctx := context.Background()
	script, err := rl.CreateScript()
	if err != nil {
		return
	}

	sha, err := script.Load(ctx, rl.client).Result()
	if err != nil {
		log.Fatalln(err)
	}
	ret := rl.client.EvalSha(ctx, sha, []string{
		userId,
		orderSet,
	}, []string{
		goodsTotal,
		orderList,
	})
	if result, err := ret.Result(); err != nil {
		log.Fatal("Redis Error")
	} else {
		total := result.(int64)
		if total == -1 {
			fmt.Printf("userid: %s, 已经抢过了 \n", userId)
		} else if total == 0 {
			fmt.Printf("userid: %s, 什么都没抢到 \n", userId)
		} else {
			fmt.Printf("userid: %s 抢到了, 库存: %d \n", userId, total)
		}
	}

}
