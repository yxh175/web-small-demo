package main

import (
	"gin-mall/auth_demo/middleware"
	"gin-mall/auth_demo/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/login", func(c *gin.Context) {

		aToken, rToken, err := util.GenerateToken(123, "guest")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"aToken": aToken,
			"rToken": rToken,
		})
	})
	r.Use(middleware.JWTMiddleware())
	r.GET("/hello", func(c *gin.Context) {
		userId, idOk := c.Get("id")
		userName, nameOk := c.Get("username")
		aToken, aTokenOK := c.Get("access-Token")
		rToken, rTokenOk := c.Get("refresh-Token")
		if !idOk || !nameOk {
			c.JSON(http.StatusOK, gin.H{"msg": "用户信息不存在"})
			return
		}

		if !aTokenOK || !rTokenOk {
			c.JSON(http.StatusOK, gin.H{
				"id":       userId,
				"userName": userName,
				"msg":      "token未更新",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":       userId,
			"username": userName,
			"aToken":   aToken,
			"rToken":   rToken,
		})
	})
	r.Run(":8000")
}
