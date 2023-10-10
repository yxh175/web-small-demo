package auth_demo

import (
	"gin-mall/auth_demo/middleware"
	"gin-mall/auth_demo/util"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUserMiddleware(t *testing.T) {
	r := gin.Default()
	r.Use(middleware.JWTMiddleware())
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
	r.GET("/hello", func(c *gin.Context) {
		userId, idOk := c.Get("id")
		userName, nameOk := c.Get("user_name")
		if !idOk || !nameOk {
			c.JSON(http.StatusOK, gin.H{"msg": "hello world"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":       userId,
			"username": userName,
		})
	})

}
