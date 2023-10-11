package main

import (
	"fmt"
	"gin-mall/cache_demo/cache"
	"gin-mall/cache_demo/db/dao"
	"gin-mall/cache_demo/service"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	cache.InitLocal()
	cache.InitRedis()
	dao.InitMySQL()

}

func main() {
	r := gin.Default()

	r.GET("/getProduct/:pid", func(c *gin.Context) {
		code := 200
		pid, err := strconv.Atoi(c.Param("pid"))
		if err != nil {
			code = 401
			c.JSON(code, gin.H{
				"status": 401,
				"msg":    "参数错误",
			})
			return
		}
		l := service.GetProductSrv()
		start := time.Now()
		product, err := l.GetData(c, uint(pid))
		end := time.Now()
		fmt.Printf("本次执行花费:%v\n", end.Sub(start))
		if err != nil {
			code = 500
			c.JSON(code, gin.H{
				"status": code,
				"msg":    "内部异常",
				"err":    err.Error(),
			})
			return
		}
		c.JSON(code, gin.H{
			"status": code,
			"data":   product,
		})

	})

	r.POST("/create", func(c *gin.Context) {
		code := 200
		l := service.GetProductSrv()
		start := time.Now()
		err := l.CreateData(c)
		end := time.Now()
		fmt.Printf("本次执行花费:%v\n", end.Sub(start))
		if err != nil {
			code = 500
			c.JSON(code, gin.H{
				"status": code,
				"msg":    "内部异常",
				"err":    err.Error(),
			})
			return
		}
		c.JSON(code, gin.H{
			"status": code,
			"data":   "创建成功",
		})
	})

	r.DELETE("/delete/:pid", func(c *gin.Context) {
		code := 200
		pid, err := strconv.Atoi(c.Param("pid"))
		if err != nil {
			code = 401
			c.JSON(code, gin.H{
				"status": 401,
				"msg":    "参数错误",
			})
			return
		}
		l := service.GetProductSrv()
		start := time.Now()
		err = l.DeleteData(c, uint(pid))
		end := time.Now()
		fmt.Printf("本次执行花费:%v\n", end.Sub(start))
		if err != nil {
			code = 500
			c.JSON(code, gin.H{
				"status": code,
				"msg":    "内部异常",
				"err":    err.Error(),
			})
			return
		}
		c.JSON(code, gin.H{
			"status": code,
			"data":   "删除成功",
		})

	})

	r.PUT("/update/:pid", func(c *gin.Context) {
		code := 200
		pid, err := strconv.Atoi(c.Param("pid"))
		if err != nil {
			code = 401
			c.JSON(code, gin.H{
				"status": 401,
				"msg":    "参数错误",
			})
			return
		}
		l := service.GetProductSrv()
		start := time.Now()
		err = l.UpdateData(c, uint(pid), rand.Float64())
		end := time.Now()
		fmt.Printf("本次执行花费:%v\n", end.Sub(start))
		if err != nil {
			code = 500
			c.JSON(code, gin.H{
				"status": code,
				"msg":    "内部异常",
				"err":    err.Error(),
			})
			return
		}
		c.JSON(code, gin.H{
			"status": code,
			"data":   "创建成功",
		})
	})
	r.Run(":8000")
}
