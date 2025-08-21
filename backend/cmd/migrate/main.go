package main

import (
	"fmt"
	"log"
	"social-app/shared/config"
	"social-app/shared/models"

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

	// 自动迁移数据库表
	err = db.AutoMigrate(
		&models.User{},
		&models.InterestTag{},
		&models.Friend{},
		&models.FriendRequest{},
		&models.ChatRoom{},
		&models.ChatRoomMember{},
		&models.Message{},
		&models.Post{},
		&models.PostLike{},
		&models.PostComment{},
		&models.PostShare{},
		&models.Report{},
		&models.CallSession{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("✅ 数据库迁移完成")

	// 初始化兴趣标签数据
	if err := initInterestTags(db); err != nil {
		log.Printf("Warning: Failed to initialize interest tags: %v", err)
	} else {
		fmt.Println("✅ 兴趣标签数据初始化完成")
	}
}

func initInterestTags(db *gorm.DB) error {
	// 检查是否已有数据
	var count int64
	db.Model(&models.InterestTag{}).Count(&count)
	if count > 0 {
		fmt.Println("兴趣标签数据已存在，跳过初始化")
		return nil
	}

	// 初始化兴趣标签数据
	tags := []models.InterestTag{
		// 运动健身
		{ID: "sport-001", Name: "健身", Category: "运动", Description: "喜欢健身锻炼", Icon: "💪", Color: "#FF6B6B", IsActive: true},
		{ID: "sport-002", Name: "跑步", Category: "运动", Description: "热爱跑步运动", Icon: "🏃", Color: "#4ECDC4", IsActive: true},
		{ID: "sport-003", Name: "游泳", Category: "运动", Description: "游泳爱好者", Icon: "🏊", Color: "#45B7D1", IsActive: true},
		{ID: "sport-004", Name: "瑜伽", Category: "运动", Description: "练习瑜伽", Icon: "🧘", Color: "#96CEB4", IsActive: true},
		{ID: "sport-005", Name: "篮球", Category: "运动", Description: "篮球运动", Icon: "🏀", Color: "#FFEAA7", IsActive: true},
		{ID: "sport-006", Name: "足球", Category: "运动", Description: "足球运动", Icon: "⚽", Color: "#DDA0DD", IsActive: true},

		// 音乐艺术
		{ID: "music-001", Name: "流行音乐", Category: "音乐", Description: "喜欢流行音乐", Icon: "🎵", Color: "#FF7675", IsActive: true},
		{ID: "music-002", Name: "古典音乐", Category: "音乐", Description: "古典音乐爱好者", Icon: "🎼", Color: "#6C5CE7", IsActive: true},
		{ID: "music-003", Name: "摇滚音乐", Category: "音乐", Description: "摇滚乐迷", Icon: "🎸", Color: "#A29BFE", IsActive: true},
		{ID: "music-004", Name: "唱歌", Category: "音乐", Description: "喜欢唱歌", Icon: "🎤", Color: "#FD79A8", IsActive: true},
		{ID: "music-005", Name: "乐器演奏", Category: "音乐", Description: "会演奏乐器", Icon: "🎹", Color: "#FDCB6E", IsActive: true},

		// 电影娱乐
		{ID: "movie-001", Name: "电影", Category: "娱乐", Description: "电影爱好者", Icon: "🎬", Color: "#E17055", IsActive: true},
		{ID: "movie-002", Name: "电视剧", Category: "娱乐", Description: "追剧达人", Icon: "📺", Color: "#00B894", IsActive: true},
		{ID: "movie-003", Name: "动漫", Category: "娱乐", Description: "动漫迷", Icon: "🎭", Color: "#00CEC9", IsActive: true},
		{ID: "movie-004", Name: "纪录片", Category: "娱乐", Description: "纪录片爱好者", Icon: "📹", Color: "#2D3436", IsActive: true},

		// 旅行探索
		{ID: "travel-001", Name: "旅行", Category: "旅行", Description: "热爱旅行", Icon: "✈️", Color: "#0984E3", IsActive: true},
		{ID: "travel-002", Name: "摄影", Category: "旅行", Description: "摄影爱好者", Icon: "📷", Color: "#6C5CE7", IsActive: true},
		{ID: "travel-003", Name: "徒步", Category: "旅行", Description: "徒步探险", Icon: "🥾", Color: "#00B894", IsActive: true},
		{ID: "travel-004", Name: "露营", Category: "旅行", Description: "户外露营", Icon: "⛺", Color: "#E17055", IsActive: true},

		// 美食烹饪
		{ID: "food-001", Name: "烹饪", Category: "美食", Description: "喜欢烹饪", Icon: "👨‍🍳", Color: "#E84393", IsActive: true},
		{ID: "food-002", Name: "烘焙", Category: "美食", Description: "烘焙爱好者", Icon: "🧁", Color: "#FDCB6E", IsActive: true},
		{ID: "food-003", Name: "品酒", Category: "美食", Description: "品酒师", Icon: "🍷", Color: "#6C5CE7", IsActive: true},
		{ID: "food-004", Name: "咖啡", Category: "美食", Description: "咖啡爱好者", Icon: "☕", Color: "#8D6E63", IsActive: true},

		// 读书学习
		{ID: "book-001", Name: "阅读", Category: "学习", Description: "爱好阅读", Icon: "📚", Color: "#2D3436", IsActive: true},
		{ID: "book-002", Name: "写作", Category: "学习", Description: "喜欢写作", Icon: "✍️", Color: "#636E72", IsActive: true},
		{ID: "book-003", Name: "学习", Category: "学习", Description: "终身学习者", Icon: "🎓", Color: "#00B894", IsActive: true},
		{ID: "book-004", Name: "语言学习", Category: "学习", Description: "学习外语", Icon: "🗣️", Color: "#0984E3", IsActive: true},

		// 科技数码
		{ID: "tech-001", Name: "编程", Category: "科技", Description: "程序员", Icon: "💻", Color: "#2D3436", IsActive: true},
		{ID: "tech-002", Name: "游戏", Category: "科技", Description: "游戏爱好者", Icon: "🎮", Color: "#6C5CE7", IsActive: true},
		{ID: "tech-003", Name: "数码产品", Category: "科技", Description: "数码达人", Icon: "📱", Color: "#00CEC9", IsActive: true},
		{ID: "tech-004", Name: "人工智能", Category: "科技", Description: "AI爱好者", Icon: "🤖", Color: "#A29BFE", IsActive: true},

		// 艺术创作
		{ID: "art-001", Name: "绘画", Category: "艺术", Description: "绘画爱好者", Icon: "🎨", Color: "#E84393", IsActive: true},
		{ID: "art-002", Name: "手工制作", Category: "艺术", Description: "手工达人", Icon: "✂️", Color: "#FDCB6E", IsActive: true},
		{ID: "art-003", Name: "设计", Category: "艺术", Description: "设计师", Icon: "🎯", Color: "#00B894", IsActive: true},
		{ID: "art-004", Name: "舞蹈", Category: "艺术", Description: "舞蹈爱好者", Icon: "💃", Color: "#FD79A8", IsActive: true},

		// 宠物动物
		{ID: "pet-001", Name: "猫咪", Category: "宠物", Description: "猫奴", Icon: "🐱", Color: "#FF7675", IsActive: true},
		{ID: "pet-002", Name: "狗狗", Category: "宠物", Description: "狗狗爱好者", Icon: "🐶", Color: "#FDCB6E", IsActive: true},
		{ID: "pet-003", Name: "动物保护", Category: "宠物", Description: "动物保护者", Icon: "🐾", Color: "#00B894", IsActive: true},
	}

	return db.Create(&tags).Error
}
