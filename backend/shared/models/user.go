package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           string             `json:"id" gorm:"primaryKey"`
	Phone        string             `json:"phone" gorm:"uniqueIndex"`
	Email        *string            `json:"email,omitempty"`
	Profile      UserProfile        `json:"profile" gorm:"embedded"`
	Verification VerificationStatus `json:"verification" gorm:"embedded"`
	Preferences  UserPreferences    `json:"preferences" gorm:"embedded"`
	Privacy      PrivacySettings    `json:"privacy" gorm:"embedded"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	DeletedAt    gorm.DeletedAt     `json:"-" gorm:"index"`
}

// UserProfile 用户资料
type UserProfile struct {
	Nickname  string       `json:"nickname"`
	Avatar    string       `json:"avatar"`
	Age       int          `json:"age"`
	Gender    string       `json:"gender"` // 'male', 'female', 'other'
	Bio       *string      `json:"bio,omitempty"`
	Interests []string     `json:"interests" gorm:"type:text[]"`
	Photos    []string     `json:"photos" gorm:"type:text[]"`
	Location  *GeoLocation `json:"location,omitempty" gorm:"embedded"`
}

// GeoLocation 地理位置
type GeoLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
}

// VerificationStatus 认证状态
type VerificationStatus struct {
	Identity IdentityVerification `json:"identity" gorm:"embedded"`
	Face     FaceVerification     `json:"face" gorm:"embedded"`
}

// IdentityVerification 身份认证
type IdentityVerification struct {
	Verified   bool       `json:"verified"`
	IDNumber   *string    `json:"id_number,omitempty"`
	RealName   *string    `json:"real_name,omitempty"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}

// FaceVerification 人脸认证
type FaceVerification struct {
	Verified   bool       `json:"verified"`
	Confidence float64    `json:"confidence"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}

// UserPreferences 用户偏好设置
type UserPreferences struct {
	AgeRange      AgeRange `json:"age_range" gorm:"embedded"`
	GenderFilter  string   `json:"gender_filter"` // 'all', 'male', 'female'
	DistanceRange int      `json:"distance_range"` // 搜索半径(公里)
	ShowOnline    bool     `json:"show_online"`
	ShowDistance  bool     `json:"show_distance"`
}

// AgeRange 年龄范围
type AgeRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// PrivacySettings 隐私设置
type PrivacySettings struct {
	ProfileVisibility string `json:"profile_visibility"` // 'public', 'friends', 'private'
	LocationVisible   bool   `json:"location_visible"`
	OnlineStatus      bool   `json:"online_status"`
	ReadReceipts      bool   `json:"read_receipts"`
}

// BeforeCreate GORM钩子，创建前生成UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}