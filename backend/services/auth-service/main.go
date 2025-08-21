package main

import (
	"fmt"
	"log"
	"social-app/shared/config"
	"social-app/services/auth-service/handlers"
	"social-app/services/auth-service/middleware"
	"social-app/services/auth-service/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化服务层
	authService := service.NewAuthService(cfg)

	// 初始化处理器层
	authHandler := handlers.NewAuthHandler(authService)

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
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/validate", authHandler.ValidateToken)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 认证验证路由
		verify := api.Group("/verify")
		{
			verify.POST("/identity", authHandler.VerifyIdentity)
			verify.POST("/face", authHandler.VerifyFace)
			verify.POST("/liveness", authHandler.VerifyLiveness)
			verify.POST("/ocr/idcard", authHandler.OCRIDCard)
			verify.GET("/status/:user_id", authHandler.GetVerificationStatus)
			verify.POST("/retry/:user_id", authHandler.RetryVerification)
		}

		// OAuth路由
		oauth := api.Group("/oauth")
		{
			// 微信登录
			oauth.GET("/wechat/login", authHandler.WeChatLogin)
			oauth.GET("/wechat/callback", authHandler.WeChatCallback)
			
			// QQ登录
			oauth.GET("/qq/login", authHandler.QQLogin)
			oauth.GET("/qq/callback", authHandler.QQCallback)
			
			// Apple登录
			oauth.GET("/apple/login", authHandler.AppleLogin)
			oauth.GET("/apple/callback", authHandler.AppleCallback)
		}

		// 文件上传路由
		upload := api.Group("/upload")
		{
			upload.POST("/avatar", authHandler.UploadAvatar)
			upload.POST("/photo", authHandler.UploadPhoto)
		}
	}

	// 启动服务器
	port := "8002" // 认证服务使用8012端口
	fmt.Printf("🔐 认证服务启动在端口 %s\n", port)
	log.Fatal(r.Run(":" + port))
}