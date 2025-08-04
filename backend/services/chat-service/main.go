package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"social-app/shared/config"
	"social-app/services/chat-service/handlers"
	"social-app/services/chat-service/repository"
	"social-app/services/chat-service/service"
	"social-app/services/chat-service/websocket"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 为了演示，使用内存数据库
	db, err := gorm.Open(postgres.Open("host=localhost user=chat password=password dbname=social_app port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		// 如果PostgreSQL连接失败，使用SQLite内存数据库
		log.Printf("PostgreSQL连接失败，使用内存数据库: %v", err)
		// 这里可以使用SQLite作为备选方案，但为了简化，我们先跳过数据库
		db = nil
	}

	// 连接MongoDB数据库（用于消息）
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Printf("MongoDB连接失败，使用内存存储: %v", err)
		mongoClient = nil
	}

	var mongoDB *mongo.Database
	if mongoClient != nil {
		// 测试MongoDB连接
		err = mongoClient.Ping(context.TODO(), nil)
		if err != nil {
			log.Printf("MongoDB ping失败: %v", err)
			mongoClient = nil
		} else {
			mongoDB = mongoClient.Database("social_app")
			fmt.Println("✅ 成功连接到MongoDB")
		}
	}

	// 初始化仓库层
	var chatRepo repository.ChatRepository
	if db != nil && mongoDB != nil {
		chatRepo = repository.NewChatRepository(db, mongoDB)
	} else {
		// 使用内存仓库进行测试
		fmt.Println("⚠️  使用内存数据库进行测试")
		chatRepo = repository.NewMemoryChatRepository()
	}

	// 创建WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// 初始化服务层
	chatService := service.NewChatService(chatRepo, cfg, hub)

	// 初始化处理器层
	chatHandler := handlers.NewChatHandler(chatService)

	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	r := gin.Default()

	// 中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "chat-service",
			"status":  "healthy",
		})
	})

	// WebSocket连接
	r.GET("/ws", hub.HandleWebSocket)

	// 聊天相关路由
	api := r.Group("/api/v1/chat")
	{
		// 需要认证的路由（这里简化处理，实际应该添加JWT中间件）
		api.GET("/rooms", chatHandler.GetChatRooms)
		api.POST("/rooms", chatHandler.CreateOrGetChatRoom)
		api.GET("/:chatId/messages", chatHandler.GetMessages)
		api.POST("/messages", chatHandler.SendMessage)
		api.POST("/messages/read", chatHandler.MarkAsRead)
		api.DELETE("/messages/:messageId", chatHandler.DeleteMessage)
		api.POST("/upload", chatHandler.UploadMedia)
	}

	// 测试路由（模拟用户认证）
	test := r.Group("/api/v1/test/chat")
	test.Use(func(c *gin.Context) {
		// 模拟用户认证，设置测试用户ID
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "test-user-1" // 默认测试用户
		}
		c.Set("user_id", userID)
		c.Next()
	})
	{
		test.GET("/rooms", chatHandler.GetChatRooms)
		test.POST("/rooms", chatHandler.CreateOrGetChatRoom)
		test.GET("/:chatId/messages", chatHandler.GetMessages)
		test.POST("/messages", chatHandler.SendMessage)
		test.POST("/messages/read", chatHandler.MarkAsRead)
		test.DELETE("/messages/:messageId", chatHandler.DeleteMessage)
		test.POST("/upload", chatHandler.UploadMedia)
	}

	// 启动服务
	port := "8004"
	fmt.Printf("💬 聊天服务启动在端口 %s\n", port)
	log.Fatal(r.Run(":" + port))
}