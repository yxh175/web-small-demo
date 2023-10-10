package middleware

import (
	"gin-mall/auth_demo/util"
	"time"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := 200
		accessToken := c.GetHeader("access_token")
		refreshToken := c.GetHeader("refresh_token")

		if accessToken == "" {
			code = 401
			c.JSON(200, gin.H{
				"status": code,
				"data":   "token不能为空",
			})
			c.Abort()
			return
		}
		// 先验证accessToken
		accessClaims, err := util.ParseToken(accessToken)
		if err != nil {
			code = 401
			c.JSON(200, gin.H{
				"status": code,
				"data":   "token异常",
			})
			c.Abort()
			return
		}
		// accessToken未过期
		if accessClaims.ExpiresAt < time.Now().Unix() {
			c.Header("access-Token", accessToken)
			c.Header("refresh-Token", refreshToken)
			c.Set("id", accessClaims.ID)
			c.Set("username", accessClaims.UserName)
			c.Next()
			return
		}

		// 判断refresh是否过期
		newAccessToken, newFreshToken, err := util.ParseRefreshToken(accessToken, refreshToken)
		if err != nil {
			code = 403
		}

		if code != 200 {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    "鉴权失败",
				"err":    err.Error(),
			})
			c.Abort()
			return
		}

		// 更新成功
		claims, err := util.ParseToken(newAccessToken)
		if err != nil {
			code = 401
			c.JSON(200, gin.H{
				"status": code,
				"msg":    "解析错误",
				"data":   err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("id", claims.ID)
		c.Set("username", claims.UserName)
		c.Header("access-Token", newAccessToken)
		c.Header("refresh-Token", newFreshToken)
		c.Next()
	}
}
