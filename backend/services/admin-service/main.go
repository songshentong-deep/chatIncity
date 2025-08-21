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
			"service": "admin-service",
			"status":  "healthy",
		})
	})

	// 管理相关路由
	api := r.Group("/api/v1/admin")
	{
		api.GET("/users", getUsers)
		api.POST("/users/:id/ban", banUser)
		api.GET("/reports", getReports)
		api.POST("/reports/:id/process", processReport)
		api.GET("/dashboard", getDashboard)
	}

	// 启动服务
	port := "8008"
	log.Printf("Admin Service starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}

// 获取用户列表
func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get users - TODO: implement",
	})
}

// 封禁用户
func banUser(c *gin.Context) {
	userId := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Ban user " + userId + " - TODO: implement",
	})
}

// 获取举报列表
func getReports(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get reports - TODO: implement",
	})
}

// 处理举报
func processReport(c *gin.Context) {
	reportId := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Process report " + reportId + " - TODO: implement",
	})
}

// 获取数据看板
func getDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get dashboard data - TODO: implement",
	})
}