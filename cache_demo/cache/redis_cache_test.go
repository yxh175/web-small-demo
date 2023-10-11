package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-mall/cache_demo/db/dao"
	"gin-mall/cache_demo/db/model"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func TestPing(t *testing.T) {
	redisCache := NewRedisCache()
	redisCache.Ping()
}

func TestGetByRedis(t *testing.T) {
	dao.InitMySQL()
	redisCache := NewRedisCache()
	productDao := dao.NewProductDao(context.Background())

	var pid uint = 1
	// DB创建两个，防止bufferPool的影响
	productDao.CreateProduct(&model.Product{
		Name:  "DB测试",
		PID:   2,
		Price: 123,
		Count: 33,
	})
	productDao.CreateProduct(&model.Product{
		Name:  "缓存测试",
		PID:   pid,
		Price: 123,
		Count: 33,
	})

	// 普通DB查询
	startTime := time.Now()
	product, err := productDao.GetProduct(2)
	if err != nil {
		t.Fatal(err)
	}
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Println("--------普通DB查询-------")
	fmt.Printf("运行时长：%v\n", duration)
	fmt.Println(product)
	fmt.Printf("-----------------------\n\n\n")

	// 测试有缓存DB第一次查询
	startTime = time.Now()
	product, err = getByRedis(redisCache, pid, productDao)
	if err != nil {
		t.Fatal(err)
	}
	endTime = time.Now()
	duration = endTime.Sub(startTime)
	fmt.Println("-----有缓存DB第一查询-----")
	fmt.Printf("运行时长：%v\n", duration)
	fmt.Println(product)
	fmt.Printf("-----------------------\n\n\n")

	startTime = time.Now()
	product, err = getByRedis(redisCache, pid, productDao)
	if err != nil {
		t.Fatal(err)
	}
	endTime = time.Now()
	duration = endTime.Sub(startTime)
	fmt.Println("-----有缓存缓存查询-----")
	fmt.Printf("运行时长：%v\n", duration)
	fmt.Println(product)
	fmt.Printf("-----------------------\n\n\n")
}

func getByRedis(redisCache *RedisCache, pid uint, dao *dao.ProductDao) (product *model.Product, err error) {
	localKey := key + fmt.Sprint(pid)
	// 从缓存中获取数据
	value, err := redisCache.Get(localKey)

	if err != nil {
		if err != redis.Nil {
			return
		}
		// 查询数据库中的数据
		product, err = dao.GetProduct(pid)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// 查询无果
				return &model.Product{}, nil
			} else {
				return
			}
		}
		var jsonData []byte
		jsonData, err = json.Marshal(*product)
		if err != nil {
			fmt.Println("序列化失败:", err)
			return
		}
		// 更新缓存
		redisCache.Set(localKey, string(jsonData), 10*time.Second)
		return
	}
	product = &model.Product{}
	jsonData := []byte(value)
	err = json.Unmarshal(jsonData, product)
	if err != nil {
		fmt.Println("反序列化失败:", err)
		return
	}
	return
}
