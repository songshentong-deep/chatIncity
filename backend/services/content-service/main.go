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
			"service": "content-service",
			"status":  "healthy",
		})
	})

	// 内容相关路由
	api := r.Group("/api/v1/content")
	{
		api.GET("/posts", getPosts)
		api.POST("/posts", createPost)
		api.POST("/posts/:id/like", likePost)
		api.POST("/posts/:id/comment", commentPost)
	}

	// 启动服务
	port := "8005"
	log.Printf("Content Service starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}

// 获取动态列表
func getPosts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get posts - TODO: implement",
	})
}

// 创建动态
func createPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Create post - TODO: implement",
	})
}

// 点赞动态
func likePost(c *gin.Context) {
	postId := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Like post " + postId + " - TODO: implement",
	})
}

// 评论动态
func commentPost(c *gin.Context) {
	postId := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Comment on post " + postId + " - TODO: implement",
	})
}