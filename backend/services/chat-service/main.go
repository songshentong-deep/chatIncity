package main

import (
	"fmt"
	"log"
	"os"
	"social-app/services/chat-service/handlers"
	"social-app/services/chat-service/repository"
	"social-app/services/chat-service/service"
	"social-app/shared/config"
	"social-app/shared/database"
	"social-app/services/chat-service/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	godotenv.Load()

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("配置加载失败:", err)
	}

	// 构建数据库连接字符串
	databaseURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.PostgreSQL.Host,
		cfg.Database.PostgreSQL.Port,
		cfg.Database.PostgreSQL.User,
		cfg.Database.PostgreSQL.Password,
		cfg.Database.PostgreSQL.DBName,
		cfg.Database.PostgreSQL.SSLMode,
	)

	// 连接数据库
	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 初始化仓库
	chatRepo := repository.NewChatRepository(db)

	// 初始化服务
	chatService := service.NewChatService(chatRepo, cfg)

	// 初始化处理器
	chatHandler := handlers.NewChatHandler(chatService)

	// 设置路由
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// API路由组
	api := r.Group("/api/v1")
	{
		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			// 聊天室管理
			auth.GET("/chat/rooms", chatHandler.GetChatRooms)
			auth.POST("/chat/rooms", chatHandler.CreateChatRoom)
			auth.POST("/chat/private", chatHandler.GetOrCreatePrivateChat)

			// 消息管理
			auth.GET("/chat/rooms/:chatId/messages", chatHandler.GetMessages)
			auth.POST("/chat/messages", chatHandler.SendMessage)
			auth.PUT("/chat/rooms/:chatId/read", chatHandler.MarkAsRead)
			auth.DELETE("/chat/messages/:messageId", chatHandler.DeleteMessage)

			// 媒体上传
			auth.POST("/chat/upload", chatHandler.UploadMedia)
			
			// 群成员管理
			auth.GET("/chat/rooms/:chatId/members", chatHandler.GetGroupMembers)
		}
	}

	// 启动服务器
	port := os.Getenv("CHAT_SERVICE_PORT")
	if port == "" {
		port = "8004" // 使用8004端口
	}

	log.Printf("🚀 聊天服务启动在端口 %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}