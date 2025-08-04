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
			"service": "notification-service",
			"status":  "healthy",
		})
	})

	// 通知相关路由
	api := r.Group("/api/v1/notification")
	{
		api.POST("/send", sendNotification)
		api.GET("/history", getNotificationHistory)
		api.PUT("/settings", updateNotificationSettings)
	}

	// 启动服务
	port := "8007"
	log.Printf("Notification Service starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}

// 发送通知
func sendNotification(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Send notification - TODO: implement",
	})
}

// 获取通知历史
func getNotificationHistory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get notification history - TODO: implement",
	})
}

// 更新通知设置
func updateNotificationSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Update notification settings - TODO: implement",
	})
}