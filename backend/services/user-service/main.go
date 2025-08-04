package main

import (
	"fmt"
	"log"
	"social-app/shared/config"
	"social-app/services/user-service/handlers"
	"social-app/services/user-service/middleware"
	"social-app/services/user-service/repository"
	"social-app/services/user-service/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	// 初始化仓库层
	userRepo := repository.NewUserRepository(db)
	interestRepo := repository.NewInterestRepository(db)

	// 初始化服务层
	userService := service.NewUserService(userRepo, interestRepo, cfg)

	// 初始化处理器层
	userHandler := handlers.NewUserHandler(userService)

	// 设置Gin路由
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// API路由组
	api := r.Group("/api/v1")
	{
		// 公开路由（无需认证）
		public := api.Group("/auth")
		{
			public.POST("/register", userHandler.Register)
			public.POST("/login", userHandler.Login)
			public.POST("/refresh", userHandler.RefreshToken)
		}

		// 需要认证的路由
		protected := api.Group("/user")
		protected.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			protected.GET("/profile", userHandler.GetProfile)
			protected.PUT("/profile", userHandler.UpdateProfile)
			protected.POST("/avatar", userHandler.UploadAvatar)
			protected.POST("/photos", userHandler.UploadPhotos)
			protected.DELETE("/photos/:photo_id", userHandler.DeletePhoto)
			protected.POST("/interests", userHandler.UpdateInterests)
			protected.GET("/profile/completeness", userHandler.GetProfileCompleteness)
			protected.POST("/verify/identity", userHandler.VerifyIdentity)
			protected.POST("/verify/face", userHandler.VerifyFace)
		}

		// 兴趣标签路由
		interests := api.Group("/interests")
		{
			interests.GET("", userHandler.GetInterestTags)
		}
	}

	// 启动服务器
	port := cfg.Server.Port
	if port == "" {
		port = "8001"
	}

	fmt.Printf("🚀 用户服务启动在端口 %s\n", port)
	log.Fatal(r.Run(":" + port))
}