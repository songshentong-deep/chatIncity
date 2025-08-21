package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Friend 好友关系模型
type Friend struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	UserID    string         `json:"user_id" gorm:"index"`
	FriendID  string         `json:"friend_id" gorm:"index"`
	Status    FriendStatus   `json:"status" gorm:"type:varchar(20);default:'pending'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联用户信息
	User   User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Friend User `json:"friend,omitempty" gorm:"foreignKey:FriendID"`
}

// FriendRequest 好友请求模型
type FriendRequest struct {
	ID           string               `json:"id" gorm:"primaryKey"`
	FromUserID   string               `json:"from_user_id" gorm:"index"`
	ToUserID     string               `json:"to_user_id" gorm:"index"`
	Message      *string              `json:"message,omitempty"`
	Status       FriendRequestStatus  `json:"status" gorm:"type:varchar(20);default:'pending'"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	DeletedAt    gorm.DeletedAt       `json:"-" gorm:"index"`
	
	// 关联用户信息
	FromUser User `json:"from_user,omitempty" gorm:"foreignKey:FromUserID"`
	ToUser   User `json:"to_user,omitempty" gorm:"foreignKey:ToUserID"`
}

// FriendStatus 好友状态枚举
type FriendStatus string

const (
	FriendStatusPending  FriendStatus = "pending"
	FriendStatusAccepted FriendStatus = "accepted"
	FriendStatusBlocked  FriendStatus = "blocked"
	FriendStatusRejected FriendStatus = "rejected"
)

// FriendRequestStatus 好友请求状态枚举
type FriendRequestStatus string

const (
	FriendRequestStatusPending  FriendRequestStatus = "pending"
	FriendRequestStatusAccepted FriendRequestStatus = "accepted"
	FriendRequestStatusRejected FriendRequestStatus = "rejected"
)

// BeforeCreate GORM钩子，创建前生成UUID
func (f *Friend) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return nil
}

func (fr *FriendRequest) BeforeCreate(tx *gorm.DB) error {
	if fr.ID == "" {
		fr.ID = uuid.New().String()
	}
	return nil
}