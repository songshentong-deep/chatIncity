package service

import (
	"mime/multipart"
	"time"
)

// 创建聊天房间请求
type CreateChatRoomRequest struct {
	ParticipantID string `json:"participant_id" binding:"required"`
	Type          string `json:"type" binding:"required,oneof=direct group"`
	Name          string `json:"name,omitempty"`
}

// 发送消息请求
type SendMessageRequest struct {
	ChatID   string                `json:"chat_id" binding:"required"`
	Type     string                `json:"type" binding:"required,oneof=text image video audio game system"`
	Content  MessageContentRequest `json:"content" binding:"required"`
	ReplyTo  string                `json:"reply_to,omitempty"`
}

// 消息内容请求
type MessageContentRequest struct {
	Text     string            `json:"text,omitempty"`
	Media    *MediaFileRequest `json:"media,omitempty"`
	GameData *GameDataRequest  `json:"game_data,omitempty"`
}

// 媒体文件请求
type MediaFileRequest struct {
	URL       string `json:"url" binding:"required"`
	Type      string `json:"type" binding:"required,oneof=image video audio"`
	Size      int64  `json:"size"`
	Duration  int    `json:"duration,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

// 游戏数据请求
type GameDataRequest struct {
	GameType string `json:"game_type" binding:"required,oneof=truth_dare draw_guess quick_qa"`
	Status   string `json:"status" binding:"required,oneof=waiting playing finished"`
	Data     string `json:"data"`
}

// 标记已读请求
type MarkAsReadRequest struct {
	ChatID    string `json:"chat_id" binding:"required"`
	MessageID string `json:"message_id" binding:"required"`
}

// 聊天房间响应
type ChatRoomResponse struct {
	ID            string                 `json:"id"`
	Participants  []string               `json:"participants"`
	Type          string                 `json:"type"`
	Name          string                 `json:"name,omitempty"`
	Avatar        string                 `json:"avatar,omitempty"`
	LastMessage   *MessageResponse       `json:"last_message,omitempty"`
	UnreadCount   int                    `json:"unread_count"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// 消息响应
type MessageResponse struct {
	ID        string                    `json:"id"`
	ChatID    string                    `json:"chat_id"`
	SenderID  string                    `json:"sender_id"`
	Content   MessageContentResponse    `json:"content"`
	Type      string                    `json:"type"`
	ReplyTo   string                    `json:"reply_to,omitempty"`
	Timestamp time.Time                 `json:"timestamp"`
	Status    string                    `json:"status"`
}

// 消息内容响应
type MessageContentResponse struct {
	Text     string                `json:"text,omitempty"`
	Media    *MediaFileResponse    `json:"media,omitempty"`
	GameData *GameDataResponse     `json:"game_data,omitempty"`
}

// 媒体文件响应
type MediaFileResponse struct {
	URL       string `json:"url"`
	Type      string `json:"type"`
	Size      int64  `json:"size"`
	Duration  int    `json:"duration,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

// 游戏数据响应
type GameDataResponse struct {
	GameType string `json:"game_type"`
	Status   string `json:"status"`
	Data     string `json:"data"`
}

// 分页消息响应
type MessagesResponse struct {
	Messages    []MessageResponse `json:"messages"`
	Total       int               `json:"total"`
	Page        int               `json:"page"`
	Limit       int               `json:"limit"`
	HasMore     bool              `json:"has_more"`
}

// 媒体上传响应
type MediaUploadResponse struct {
	URL       string `json:"url"`
	Type      string `json:"type"`
	Size      int64  `json:"size"`
	Duration  int    `json:"duration,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}



// API响应包装
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 成功响应
func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// 错误响应
func ErrorResponse(code int, message string) APIResponse {
	return APIResponse{
		Code:    code,
		Message: message,
	}
}

// 聊天服务接口
type ChatService interface {
	GetChatRooms(userID string) ([]ChatRoomResponse, error)
	CreateOrGetChatRoom(userID string, req *CreateChatRoomRequest) (*ChatRoomResponse, error)
	GetMessages(userID string, chatID string, page, limit int) (*MessagesResponse, error)
	SendMessage(userID string, req *SendMessageRequest) (*MessageResponse, error)
	MarkAsRead(userID string, req *MarkAsReadRequest) error
	DeleteMessage(userID string, messageID string) error
	UploadMedia(userID string, file *multipart.FileHeader, mediaType string) (*MediaUploadResponse, error)
}