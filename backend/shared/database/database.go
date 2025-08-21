package database

import (
	"fmt"
	"log"
	"social-app/shared/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect 连接到数据库
func Connect(databaseURL string) (*gorm.DB, error) {
	// 如果没有提供完整的URL，使用环境变量构建
	if databaseURL == "" {
		databaseURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			"localhost", "5432", "chat", "password", "social_app", "disable")
	}

	// 配置GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(databaseURL), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// 自动迁移数据库表
	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("数据库连接成功")
	return db, nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Friend{},
		&models.FriendRequest{},
		&models.Message{},
		&models.ChatRoom{},
		&models.ChatRoomMember{},
		&models.CallSession{},
	)
}

// GetDB 获取数据库连接（用于测试）
func GetDB() *gorm.DB {
	db, err := Connect("")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return db
}