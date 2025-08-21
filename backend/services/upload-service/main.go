package main

import (
	"fmt"
	"log"
	"os"
	"social-app/services/upload-service/handlers"
	"social-app/services/upload-service/repository"
	"social-app/services/upload-service/service"
	"social-app/shared/config"
	"social-app/shared/database"
	"social-app/services/upload-service/middleware"
	"social-app/shared/models"
	"social-app/shared/storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载环境变量
	godotenv.Load()

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
	photoRepo := repository.NewPhotoRepository(mongoClient.Database)

	// Redis客户端（暂时设为nil）
	var redisClient *redis.Client = nil

	// 初始化OBS客户端
	obsAccessKey := os.Getenv("OBS_ACCESS_KEY")
	obsSecretKey := os.Getenv("OBS_SECRET_KEY")
	obsEndpoint := os.Getenv("OBS_ENDPOINT")
	obsBucket := os.Getenv("OBS_BUCKET")
	
	if obsAccessKey == "" || obsSecretKey == "" {
		log.Fatal("OBS配置缺失，请设置OBS_ACCESS_KEY和OBS_SECRET_KEY环境变量")
	}
	
	obsClient := storage.NewOBSClientWithAuth(obsEndpoint, obsBucket, obsAccessKey, obsSecretKey)

	// 连接PostgreSQL数据库
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "chat", "password", "social_app")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("连接PostgreSQL失败: %v", err)
		db = nil
	} else {
		// 自动迁移表结构
		db.AutoMigrate(&models.UserAvatar{})
		log.Println("PostgreSQL连接成功")
	}

	// 初始化服务
	photoService := service.NewPhotoService(photoRepo, cfg.Upload.StoragePath, redisClient, obsClient, db)

	// 初始化处理器
	photoHandler := handlers.NewPhotoHandler(photoService)
	uploadHandler := handlers.NewUploadHandler(photoService)

	// 设置路由
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// 静态文件服务
	r.Static("/uploads", cfg.Upload.StoragePath)

	// API路由组
	api := r.Group("/api/v1")
	{
		// 文件上传（需要认证）
		upload := api.Group("/upload")
		upload.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			upload.POST("/image", uploadHandler.UploadImage)
			upload.POST("/video", uploadHandler.UploadVideo)
			upload.POST("/file", uploadHandler.UploadFile) // 通用文件上传
		}
		
		// 照片管理（需要认证）
		auth := api.Group("/")
		auth.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			auth.POST("/photos", photoHandler.UploadPhoto)
			auth.POST("/photos/upload-bytes", photoHandler.UploadPhotoBytes) // Web平台字节上传
		}
		
		// 公开接口（不需要认证）
		api.GET("/photos", photoHandler.GetUserPhotos)
		api.GET("/photos/:id", photoHandler.GetPhoto)
		api.GET("/photos/:id/data", photoHandler.GetPhotoData)
		api.DELETE("/photos/:id", photoHandler.DeletePhoto)
		api.PUT("/photos/:id/avatar", photoHandler.SetActiveAvatar)
		api.GET("/users/:user_id/avatar", photoHandler.GetUserAvatar)
		api.GET("/proxy/image", photoHandler.ProxyImage) // 图片代理
		api.GET("/proxy/upload-image", uploadHandler.ProxyImage) // 上传服务图片代理
		api.GET("/files/:id", uploadHandler.GetFile) // 获取上传的文件
	}

	// 启动服务器
	port := os.Getenv("PHOTO_SERVICE_PORT")
	if port == "" {
		port = "8009" // 使用8009端口
	}

	log.Printf("🚀 文件上传服务启动在端口 %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}