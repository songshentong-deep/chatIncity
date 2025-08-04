package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatRoom 聊天房间
type ChatRoom struct {
	ID           string     `json:"id" gorm:"primaryKey"`
	Participants string     `json:"participants"` // 存储为逗号分隔的字符串
	Type         string     `json:"type"` // 'direct', 'group'
	Name         *string    `json:"name,omitempty"`
	Avatar       *string    `json:"avatar,omitempty"`
	LastMessage  *Message   `json:"last_message,omitempty" gorm:"foreignKey:ChatID"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// Message 消息
type Message struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	ChatID    string         `json:"chat_id" gorm:"index"`
	SenderID  string         `json:"sender_id"`
	Content   MessageContent `json:"content" gorm:"embedded"`
	Type      string         `json:"type"` // 'text', 'image', 'video', 'audio', 'game', 'system'
	ReplyTo   *string        `json:"reply_to,omitempty"` // 回复的消息ID
	Timestamp time.Time      `json:"timestamp"`
	Status    string         `json:"status"` // 'sent', 'delivered', 'read'
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// MessageContent 消息内容
type MessageContent struct {
	Text     *string      `json:"text,omitempty"`
	Media    *MediaFile   `json:"media,omitempty" gorm:"embedded"`
	GameData *GameSession `json:"game_data,omitempty" gorm:"embedded"`
}

// MediaFile 媒体文件
type MediaFile struct {
	URL       string  `json:"url"`
	Type      string  `json:"type"` // 'image', 'video', 'audio'
	Size      int64   `json:"size"`
	Duration  *int    `json:"duration,omitempty"` // 音视频时长(秒)
	Thumbnail *string `json:"thumbnail,omitempty"` // 缩略图URL
}

// GameSession 游戏会话
type GameSession struct {
	GameType string `json:"game_type"` // 'truth_dare', 'draw_guess', 'quick_qa'
	Status   string `json:"status"`    // 'waiting', 'playing', 'finished'
	Data     string `json:"data" gorm:"type:text"` // 游戏数据(JSON字符串)
}

// MessageStatus 消息状态记录
type MessageStatus struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	MessageID string    `json:"message_id" gorm:"index"`
	UserID    string    `json:"user_id" gorm:"index"`
	Status    string    `json:"status"` // 'delivered', 'read'
	Timestamp time.Time `json:"timestamp"`
}

// CallSession 通话会话
type CallSession struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	ChatID       string    `json:"chat_id" gorm:"index"`
	InitiatorID  string    `json:"initiator_id"`
	ParticipantID string   `json:"participant_id"`
	Type         string    `json:"type"` // 'voice', 'video'
	Status       string    `json:"status"` // 'ringing', 'connected', 'ended', 'missed'
	StartTime    time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	Duration     int       `json:"duration"` // 通话时长(秒)
}

// BeforeCreate GORM钩子
func (cr *ChatRoom) BeforeCreate(tx *gorm.DB) error {
	if cr.ID == "" {
		cr.ID = uuid.New().String()
	}
	return nil
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

func (ms *MessageStatus) BeforeCreate(tx *gorm.DB) error {
	if ms.ID == "" {
		ms.ID = uuid.New().String()
	}
	return nil
}

func (cs *CallSession) BeforeCreate(tx *gorm.DB) error {
	if cs.ID == "" {
		cs.ID = uuid.New().String()
	}
	return nil
}