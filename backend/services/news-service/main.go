package main

import (
	"log"
	"os"
	"time"
	"social-app/services/news-service/handlers"
	"social-app/services/news-service/repository"
	"social-app/services/news-service/service"
	"social-app/services/news-service/scheduler"
	"social-app/services/news-service/middleware"
	"social-app/shared/config"
	"social-app/shared/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "***" + key[len(key)-4:]
}

func main() {
	// 加载环境变量 - 尝试多个路径
	envPaths := []string{".env", "../.env", "../../.env", "../../../.env"}
	envLoaded := false
	
	for _, envPath := range envPaths {
		if err := godotenv.Load(envPath); err == nil {
			log.Printf("✅ 成功加载环境变量文件: %s", envPath)
			envLoaded = true
			break
		}
	}
	
	if !envLoaded {
		log.Printf("警告: 无法加载 .env 文件，尝试使用系统环境变量")
	}

	// 检查 NEWS_API_KEY
	if apiKey := os.Getenv("NEWS_API_KEY"); apiKey == "" {
		log.Printf("警告: NEWS_API_KEY 环境变量未设置，将使用测试数据")
		log.Printf("当前工作目录: %s", getCurrentDir())
	} else {
		log.Printf("✅ NEWS_API_KEY 已配置: %s", maskAPIKey(apiKey))
	}

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("配置加载失败:", err)
	}

	// 连接MongoDB
	mongoClient, err := database.ConnectMongoDB(cfg.Database.MongoDB.URI, cfg.Database.MongoDB.Database)
	if err != nil {
		log.Fatal("MongoDB连接失败:", err)
	}
	defer mongoClient.Close()

	// 初始化仓库
	newsRepo := repository.NewNewsRepository(mongoClient.Database)

	// 初始化服务
	newsService := service.NewNewsService(newsRepo)

	// 启动定时任务
	newsScheduler := scheduler.NewNewsScheduler(newsService)
	newsScheduler.Start()
	
	// 立即刷新一次新闻（用于测试）
	go func() {
		time.Sleep(2 * time.Second) // 等待服务完全启动
		log.Println("🔄 手动触发新闻刷新...")
		if err := newsService.RefreshNews(); err != nil {
			log.Printf("手动刷新新闻失败: %v", err)
		} else {
			log.Println("✅ 手动刷新新闻成功")
		}
	}()

	// 初始化处理器
	newsHandler := handlers.NewNewsHandler(newsService)

	// 设置路由
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// API路由组
	api := r.Group("/api/v1")
	{
		// 新闻相关接口
		api.GET("/news", newsHandler.GetNews)
		api.GET("/news/categories", newsHandler.GetCategories)
		api.GET("/news/category/:category", newsHandler.GetNewsByCategory)
		api.GET("/news/:id", newsHandler.GetNewsDetail)
		
		// 管理接口（需要认证）
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			admin.POST("/news/refresh", newsHandler.RefreshNews)
		}
	}

	// 启动服务器
	port := os.Getenv("NEWS_SERVICE_PORT")
	if port == "" {
		port = "8011"
	}

	log.Printf("🚀 新闻服务启动在端口 %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}