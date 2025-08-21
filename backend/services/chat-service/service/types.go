package service

import (
	"social-app/shared/models"
)

// 发送消息请求
type SendMessageRequest struct {
	ChatRoomID  string                `json:"chat_room_id" binding:"required"`
	Content     string                `json:"content" binding:"required"`
	MessageType models.MessageType    `json:"message_type"`
}

// 创建聊天室请求
type CreateChatRoomRequest struct {
	Type        models.ChatRoomType `json:"type" binding:"required"`
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	MemberIDs   []string            `json:"member_ids" binding:"required"`
}

// 获取或创建私聊请求
type GetOrCreatePrivateChatRequest struct {
	FriendID string `json:"friend_id" binding:"required"`
}

// 标记已读请求
type MarkAsReadRequest struct {
	ChatRoomID string `json:"chat_room_id" binding:"required"`
}

// 消息响应
type MessageResponse struct {
	ID           string `json:"id"`
	ChatRoomID   string `json:"chat_room_id"`
	SenderID     string `json:"sender_id"`
	SenderName   string `json:"sender_name"`
	SenderAvatar string `json:"sender_avatar"`
	Content      string `json:"content"`
	MessageType  string `json:"message_type"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

// 聊天室响应
type ChatRoomResponse struct {
	ID           string            `json:"id"`
	Type         string            `json:"type"`
	Name         string            `json:"name"`
	Description  string            `json:"description,omitempty"`
	CreatedBy    string            `json:"created_by"`
	CreatedAt    string            `json:"created_at"`
	LastMessage  *MessageResponse  `json:"last_message,omitempty"`
	UnreadCount  int               `json:"unread_count"`
	Members      []MemberResponse  `json:"members,omitempty"`
}

// 成员响应
type MemberResponse struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
	JoinedAt string `json:"joined_at"`
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