package models

import (
	"time"
	"gorm.io/gorm"
)

// RecentAccount 最近登录账号模型
type RecentAccount struct {
	ID               string         `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Phone            string         `json:"phone" gorm:"type:varchar(20);not null;index"`
	Password         string         `json:"password,omitempty" gorm:"type:text"`
	RememberPassword bool           `json:"remember_password" gorm:"default:false"`
	LastLoginAt      time.Time      `json:"last_login_at" gorm:"not null"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}