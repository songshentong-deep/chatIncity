package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT认证中间件
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			c.Abort()
			return
		}

		// 提取token
		tokenString := authHeader[7:]
		fmt.Printf("JWT Debug: Token前10位: %s...\n", tokenString[:10])

		// 解析token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		fmt.Printf("JWT Debug: Token解析结果 - Valid: %v, Error: %v\n", token != nil && token.Valid, err)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			c.Abort()
			return
		}

		// 提取claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Printf("JWT Debug: Claims提取失败\n")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			c.Abort()
			return
		}
		fmt.Printf("JWT Debug: Claims: %+v\n", claims)

		// 提取用户ID
		userID, ok := claims["user_id"].(string)
		if !ok {
			fmt.Printf("JWT Debug: user_id提取失败, claims中的user_id: %v\n", claims["user_id"])
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			c.Abort()
			return
		}
		fmt.Printf("JWT Debug: 成功提取user_id: %s\n", userID)

		// 将用户ID存储到上下文中
		c.Set("user_id", userID)
		c.Next()
	}
}