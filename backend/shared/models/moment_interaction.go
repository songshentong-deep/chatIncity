package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MomentLike 动态点赞表
type MomentLike struct {
	ID       string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MomentID string         `json:"moment_id" gorm:"type:varchar(36);not null;index"`
	UserID   string         `json:"user_id" gorm:"type:varchar(36);not null;index"`
	CreatedAt time.Time     `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// MomentComment 动态评论表
type MomentComment struct {
	ID       string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MomentID string         `json:"moment_id" gorm:"type:varchar(36);not null;index"`
	UserID   string         `json:"user_id" gorm:"type:varchar(36);not null;index"`
	Content  string         `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联
	User User `json:"user" gorm:"foreignKey:UserID"`
}

func (l *MomentLike) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}

func (c *MomentComment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}