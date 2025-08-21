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
			"service": "location-service",
			"status":  "healthy",
		})
	})

	// 位置相关路由
	api := r.Group("/api/v1/location")
	{
		api.POST("/update", updateLocation)
		api.GET("/nearby", getNearbyUsers)
		api.POST("/privacy", setLocationPrivacy)
	}

	// 启动服务
	port := "8006"
	log.Printf("Location Service starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}

// 更新位置
func updateLocation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Update location - TODO: implement",
	})
}

// 获取附近用户
func getNearbyUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get nearby users - TODO: implement",
	})
}

// 设置位置隐私
func setLocationPrivacy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Set location privacy - TODO: implement",
	})
}