package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"social-app/shared/config"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "match-service",
			"status":  "healthy",
		})
	})

	// 匹配相关路由
	api := r.Group("/api/v1/match")
	{
		api.GET("/recommendations", getRecommendations)
		api.POST("/like", likeUser)
		api.POST("/pass", passUser)
		api.GET("/matches", getMatches)
	}

	// 启动服务
	port := "8003"
	log.Printf("Match Service starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}

// 获取推荐用户
func getRecommendations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get recommendations - TODO: implement",
	})
}

// 喜欢用户
func likeUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Like user - TODO: implement",
	})
}

// 跳过用户
func passUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Pass user - TODO: implement",
	})
}

// 获取匹配列表
func getMatches(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get matches - TODO: implement",
	})
}