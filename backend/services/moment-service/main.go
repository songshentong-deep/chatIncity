package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"social-app/services/moment-service/handlers"
	"social-app/services/moment-service/middleware"
	"social-app/services/moment-service/repository"
	"social-app/services/moment-service/service"
	"social-app/shared/config"
	"social-app/shared/models"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 连接数据库
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.PostgreSQL.Host,
		cfg.Database.PostgreSQL.Port,
		cfg.Database.PostgreSQL.User,
		cfg.Database.PostgreSQL.Password,
		cfg.Database.PostgreSQL.DBName,
		cfg.Database.PostgreSQL.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(&models.Moment{}, &models.MomentLike{}, &models.MomentComment{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化仓库层
	momentRepo := repository.NewMomentRepository(db)
	interactionRepo := repository.NewInteractionRepository(db)

	// 初始化服务层
	momentService := service.NewMomentService(momentRepo, interactionRepo)

	// 初始化处理器层
	momentHandler := handlers.NewMomentHandler(momentService)
	uploadHandler := handlers.NewUploadHandler()
	interactionHandler := handlers.NewInteractionHandler(interactionRepo)

	// 设置Gin路由
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())

	// API路由组
	api := r.Group("/api/v1")
	{
		// 动态相关路由
		moments := api.Group("/moments")
		moments.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			moments.POST("", momentHandler.CreateMoment)
			moments.GET("/my", momentHandler.GetMyMoments)
			moments.GET("/friends", momentHandler.GetFriendsMoments)
			moments.DELETE("/:moment_id", momentHandler.DeleteMoment)
		}

		// 上传相关路由
		upload := api.Group("/upload")
		upload.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			upload.POST("/image", uploadHandler.UploadImage)
			upload.POST("/video", uploadHandler.UploadVideo)
		}

		// 交互相关路由
		interaction := api.Group("/moments/:moment_id")
		interaction.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			interaction.POST("/like", interactionHandler.ToggleLike)
			interaction.GET("/likes", interactionHandler.GetLikes)
			interaction.POST("/comment", interactionHandler.AddComment)
			interaction.GET("/comments", interactionHandler.GetComments)
		}
	}

	// 启动服务器
	port := "8010"
	fmt.Printf("🚀 动态服务启动在端口 %s\n", port)
	log.Fatal(r.Run(":" + port))
}
