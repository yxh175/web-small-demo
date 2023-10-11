package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-mall/cache_demo/db/dao"
	"gin-mall/cache_demo/db/model"
	"testing"
	"time"

	"gorm.io/gorm"
)

var cache = NewCache()
var key = "product:"

func TestLocalCache(t *testing.T) {
	dao.InitMySQL()
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
	product, err = getProduct(pid, productDao)
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
	product, err = getProduct(pid, productDao)
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

func getProduct(pid uint, dao *dao.ProductDao) (product *model.Product, err error) {
	localKey := key + fmt.Sprint(pid)
	// 从缓存中获取数据
	if value, ok := cache.Get(localKey); ok {
		product = &model.Product{}
		err = json.Unmarshal(value.([]byte), product)
		if err != nil {
			fmt.Println("反序列化失败:", err)
			return
		}
		return
	} else {
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
		cache.Set(localKey, jsonData, 10*time.Second)
	}
	return
}
