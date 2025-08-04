package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 简单的测试认证中间件
func TestAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header中获取用户ID（测试用）
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			c.Abort()
			return
		}

		// 将用户ID存储到上下文中
		c.Set("user_id", userID)
		c.Next()
	}
}

// CORS中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-User-ID")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}