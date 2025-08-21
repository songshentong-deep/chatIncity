package models

import (
	"time"
	"gorm.io/gorm"
)

// UserAvatar 用户头像模型
type UserAvatar struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    string         `json:"user_id" gorm:"uniqueIndex;not null"`
	AvatarURL string         `json:"avatar_url" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (UserAvatar) TableName() string {
	return "user_avatars"
}