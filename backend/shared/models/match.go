package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserMatch 用户匹配记录
type UserMatch struct {
	ID           string          `json:"id" gorm:"primaryKey"`
	UserID       string          `json:"user_id" gorm:"index"`
	TargetUserID string          `json:"target_user_id" gorm:"index"`
	Score        float64         `json:"score"`
	Factors      MatchingFactors `json:"factors" gorm:"embedded"`
	Status       string          `json:"status"` // 'pending', 'liked', 'passed', 'mutual'
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// MatchingFactors 匹配因子
type MatchingFactors struct {
	InterestSimilarity    float64 `json:"interest_similarity"`
	LocationDistance      float64 `json:"location_distance"`
	BehaviorCompatibility float64 `json:"behavior_compatibility"`
	ActivityLevel         float64 `json:"activity_level"`
}

// UserBehavior 用户行为数据
type UserBehavior struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	UserID          string    `json:"user_id" gorm:"index"`
	ActionType      string    `json:"action_type"` // 'like', 'pass', 'chat', 'view_profile'
	TargetUserID    *string   `json:"target_user_id,omitempty"`
	TargetContentID *string   `json:"target_content_id,omitempty"`
	Duration        int       `json:"duration"` // 持续时间(秒)
	Timestamp       time.Time `json:"timestamp"`
}

// MatchPreference 匹配偏好记录
type MatchPreference struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	UserID       string    `json:"user_id" gorm:"index"`
	TargetUserID string    `json:"target_user_id" gorm:"index"`
	Action       string    `json:"action"` // 'like', 'pass', 'super_like'
	CreatedAt    time.Time `json:"created_at"`
}

// BeforeCreate GORM钩子
func (um *UserMatch) BeforeCreate(tx *gorm.DB) error {
	if um.ID == "" {
		um.ID = uuid.New().String()
	}
	return nil
}

func (ub *UserBehavior) BeforeCreate(tx *gorm.DB) error {
	if ub.ID == "" {
		ub.ID = uuid.New().String()
	}
	return nil
}

func (mp *MatchPreference) BeforeCreate(tx *gorm.DB) error {
	if mp.ID == "" {
		mp.ID = uuid.New().String()
	}
	return nil
}