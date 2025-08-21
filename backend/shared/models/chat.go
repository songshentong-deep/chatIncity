package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Message 消息模型
type Message struct {
	ID         string         `json:"id" gorm:"primaryKey"`
	ChatRoomID string         `json:"chat_room_id" gorm:"index"`
	SenderID   string         `json:"sender_id" gorm:"index"`
	Content    string         `json:"content"`
	MessageType MessageType   `json:"message_type" gorm:"type:varchar(20);default:'text'"`
	Status     MessageStatus  `json:"status" gorm:"type:varchar(20);default:'sent'"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联
	ChatRoom ChatRoom `json:"chat_room,omitempty" gorm:"foreignKey:ChatRoomID"`
	Sender   User     `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
}

// ChatRoom 聊天室模型
type ChatRoom struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	Type        ChatRoomType   `json:"type" gorm:"type:varchar(20);default:'private'"`
	Name        *string        `json:"name,omitempty"`
	Description *string        `json:"description,omitempty"`
	CreatedBy   string         `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联
	Creator User `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
}

// ChatRoomMember 聊天室成员模型
type ChatRoomMember struct {
	ID         string         `json:"id" gorm:"primaryKey"`
	ChatRoomID string         `json:"chat_room_id" gorm:"index"`
	UserID     string         `json:"user_id" gorm:"index"`
	Role       MemberRole     `json:"role" gorm:"type:varchar(20);default:'member'"`
	JoinedAt   time.Time      `json:"joined_at"`
	LastReadAt *time.Time     `json:"last_read_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联
	ChatRoom ChatRoom `json:"chat_room,omitempty" gorm:"foreignKey:ChatRoomID"`
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// MessageContent 消息内容（扩展功能）
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

// CallSession 通话会话
type CallSession struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	ChatRoomID   string    `json:"chat_room_id" gorm:"index"`
	InitiatorID  string    `json:"initiator_id"`
	ParticipantID string   `json:"participant_id"`
	Type         string    `json:"type"` // 'voice', 'video'
	Status       string    `json:"status"` // 'ringing', 'connected', 'ended', 'missed'
	StartTime    time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	Duration     int       `json:"duration"` // 通话时长(秒)
}

// MessageType 消息类型枚举
type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeAudio MessageType = "audio"
	MessageTypeVideo MessageType = "video"
	MessageTypeFile  MessageType = "file"
)

// MessageStatus 消息状态枚举
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
)

// ChatRoomType 聊天室类型枚举
type ChatRoomType string

const (
	ChatRoomTypePrivate ChatRoomType = "private"
	ChatRoomTypeGroup   ChatRoomType = "group"
)

// MemberRole 成员角色枚举
type MemberRole string

const (
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleMember MemberRole = "member"
)

// BeforeCreate GORM钩子，创建前生成UUID
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

func (cr *ChatRoom) BeforeCreate(tx *gorm.DB) error {
	if cr.ID == "" {
		cr.ID = uuid.New().String()
	}
	return nil
}

func (crm *ChatRoomMember) BeforeCreate(tx *gorm.DB) error {
	if crm.ID == "" {
		crm.ID = uuid.New().String()
	}
	return nil
}

func (cs *CallSession) BeforeCreate(tx *gorm.DB) error {
	if cs.ID == "" {
		cs.ID = uuid.New().String()
	}
	return nil
}