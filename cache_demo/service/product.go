package service

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gin-mall/cache_demo/cache"
	"gin-mall/cache_demo/db/dao"
	"gin-mall/cache_demo/db/model"
	"math/rand"
	"sync"
	"time"

	"github.com/bits-and-blooms/bloom"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var ProductSrvIns *ProductSrv
var ProductSrvOnce sync.Once
var key string = "product:"
var lockKey string = "lock:product:"

type ProductSrv struct {
	localCache *cache.Cache
	redisCache *cache.RedisCache
	filter     *bloom.BloomFilter
}

func GetProductSrv() *ProductSrv {
	ProductSrvOnce.Do(func() {
		ProductSrvIns = &ProductSrv{
			localCache: cache.LocalCache,
			redisCache: cache.RDCache,
			filter:     bloom.NewWithEstimates(1000000, 0.01),
		}
	})
	return ProductSrvIns
}

func (ps *ProductSrv) GetData(c *gin.Context, pId uint) (product *model.Product, err error) {
	// 创建布隆过滤器，预计容纳1000000个元素，误差率设置为0.01
	n1 := make([]byte, 8)
	fmt.Println(n1)
	binary.BigEndian.PutUint64(n1, uint64(pId))
	if !ps.filter.Test(n1) {
		fmt.Println("bloom filter")
		return
	}
	uniqueKey := key + fmt.Sprint(pId)
	// 模拟requestId
	rand.Seed(time.Now().Unix())
	requestId := fmt.Sprint(rand.Intn(10000000))
	// 查本地缓存
	if value, ok := ps.localCache.Get(uniqueKey); ok {
		// 一级缓存查询有值
		product = &model.Product{}
		err = json.Unmarshal(value.([]byte), product)
		if err != nil {
			fmt.Println("反序列化失败:", err)
			return
		}
		return
	}

	// 否则查二级缓存
	value, err := ps.redisCache.Get(c, uniqueKey)

	// 查询异常
	if err != nil {
		if err != redis.Nil {
			return
		}
		// 查询无果
		// 查询数据库中的数据
		// redis 上锁
		var locked bool
		locked, err = ps.redisCache.SetNx(c, lockKey, requestId)
		defer ps.redisCache.Unlock(c, lockKey, requestId)
		if err != nil {
			return
		}
		if !locked {
			time.Sleep(50 * time.Millisecond)
			return ps.GetData(c, pId)
		}
		product, err = dao.NewProductDao(c).GetProduct(pId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// 查询无果
				// 防止缓存穿透
				ps.localCache.Set(uniqueKey, "", 10*time.Second)
				ps.redisCache.Set(c, uniqueKey, "", 60*time.Second)
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
		ps.localCache.Set(uniqueKey, string(jsonData), 10*time.Second)
		ps.redisCache.Set(c, uniqueKey, string(jsonData), 60*time.Second)
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

func (ps *ProductSrv) UpdateData(c *gin.Context, pId uint, newPrice float64) (err error) {
	// 先操作数据库，在删除缓存是比较好的选择
	// 后删除防止别的线程进入数据库查询
	uniqueKey := key + fmt.Sprint(pId)
	// 更新数据库
	if err = dao.NewProductDao(c).UpdateProduct(pId, &model.Product{Price: newPrice}); err != nil {
		return
	}

	// 更新缓存, 从下到上
	err = ps.redisCache.Delete(c, uniqueKey)
	if err == redis.Nil {
		return nil
	}
	ps.localCache.Delete(uniqueKey)
	return
}

func (ps *ProductSrv) DeleteData(c *gin.Context, pId uint) (err error) {
	uniqueKey := key + fmt.Sprint(pId)

	if err = dao.NewProductDao(c).DeleteProduct(pId); err != nil {
		return
	}

	// 更新缓存, 从下到上
	err = ps.redisCache.Delete(c, uniqueKey)
	if err == redis.Nil {
		return nil
	}
	ps.localCache.Delete(uniqueKey)
	return
}

func (ps *ProductSrv) CreateData(c *gin.Context) (err error) {
	rand.Seed(time.Now().Unix())
	pid := rand.Intn(10000)
	uniqueKey := key + fmt.Sprint(pid)
	newProduct := &model.Product{
		PID:  uint(pid),
		Name: "测试",
	}
	if err = dao.NewProductDao(c).CreateProduct(newProduct); err != nil {
		return
	}
	data, err := json.Marshal(newProduct)
	if err != nil {
		return err
	}
	err = ps.redisCache.Set(c, uniqueKey, string(data), 60*time.Second)
	if err != nil {
		return
	}
	ps.localCache.Set(uniqueKey, data, 10*time.Second)
	n1 := make([]byte, 8)
	binary.BigEndian.PutUint64(n1, uint64(pid))
	ps.filter.Add(n1)
	return
}
