package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Post 动态内容
type Post struct {
	ID           string           `json:"id" gorm:"primaryKey"`
	UserID       string           `json:"user_id" gorm:"index"`
	Content      PostContent      `json:"content" gorm:"embedded"`
	Visibility   string           `json:"visibility"` // 'public', 'friends', 'nearby'
	Interactions PostInteractions `json:"interactions" gorm:"embedded"`
	Moderation   ModerationStatus `json:"moderation" gorm:"embedded"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `json:"-" gorm:"index"`
}

// PostContent 动态内容
type PostContent struct {
	Text     *string      `json:"text,omitempty"`
	Images   []string     `json:"images,omitempty" gorm:"type:text[]"`
	Video    *string      `json:"video,omitempty"`
	Location *GeoLocation `json:"location,omitempty" gorm:"embedded"`
	Tags     []string     `json:"tags,omitempty" gorm:"type:text[]"`
}

// PostInteractions 动态互动数据
type PostInteractions struct {
	Likes    int `json:"likes"`
	Comments int `json:"comments"`
	Shares   int `json:"shares"`
	Views    int `json:"views"`
}

// ModerationStatus 审核状态
type ModerationStatus struct {
	Status     string     `json:"status"` // 'pending', 'approved', 'rejected', 'flagged'
	ReviewedBy *string    `json:"reviewed_by,omitempty"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
	Reason     *string    `json:"reason,omitempty"`
	AIScore    float64    `json:"ai_score"` // AI审核评分
}

// PostLike 动态点赞记录
type PostLike struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	PostID    string    `json:"post_id" gorm:"index"`
	UserID    string    `json:"user_id" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
}

// PostComment 动态评论
type PostComment struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	PostID    string    `json:"post_id" gorm:"index"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	ParentID  *string   `json:"parent_id,omitempty"` // 回复评论的ID
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// PostShare 动态分享记录
type PostShare struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	PostID    string    `json:"post_id" gorm:"index"`
	UserID    string    `json:"user_id" gorm:"index"`
	Platform  string    `json:"platform"` // 'internal', 'wechat', 'weibo'
	CreatedAt time.Time `json:"created_at"`
}

// Report 举报记录
type Report struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	ReporterID   string    `json:"reporter_id" gorm:"index"`
	TargetType   string    `json:"target_type"` // 'user', 'post', 'message'
	TargetID     string    `json:"target_id" gorm:"index"`
	Reason       string    `json:"reason"`
	Description  *string   `json:"description,omitempty"`
	Status       string    `json:"status"` // 'pending', 'processing', 'resolved', 'dismissed'
	ProcessedBy  *string   `json:"processed_by,omitempty"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// InterestTag 兴趣标签
type InterestTag struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex"`
	Category    string         `json:"category"`
	Description string         `json:"description"`
	Icon        string         `json:"icon"`
	Color       string         `json:"color"`
	UsageCount  int            `json:"usage_count"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate GORM钩子
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

func (pl *PostLike) BeforeCreate(tx *gorm.DB) error {
	if pl.ID == "" {
		pl.ID = uuid.New().String()
	}
	return nil
}

func (pc *PostComment) BeforeCreate(tx *gorm.DB) error {
	if pc.ID == "" {
		pc.ID = uuid.New().String()
	}
	return nil
}

func (ps *PostShare) BeforeCreate(tx *gorm.DB) error {
	if ps.ID == "" {
		ps.ID = uuid.New().String()
	}
	return nil
}

func (r *Report) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

func (it *InterestTag) BeforeCreate(tx *gorm.DB) error {
	if it.ID == "" {
		it.ID = uuid.New().String()
	}
	return nil
}