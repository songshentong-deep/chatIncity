package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"social-app/shared/config"
	"social-app/shared/models"
	"social-app/services/chat-service/repository"
	"social-app/services/chat-service/types"
	"strings"
	"time"

	"gorm.io/gorm"
)

type chatService struct {
	chatRepo repository.ChatRepository
	config   *config.Config
	hub      WebSocketHub
}

// WebSocketHub 接口，避免循环导入
type WebSocketHub interface {
	SendToUser(userID string, message []byte)
}

func NewChatService(chatRepo repository.ChatRepository, cfg *config.Config, hub WebSocketHub) ChatService {
	return &chatService{
		chatRepo: chatRepo,
		config:   cfg,
		hub:      hub,
	}
}

// 获取用户的聊天房间列表
func (s *chatService) GetChatRooms(userID string) ([]ChatRoomResponse, error) {
	rooms, err := s.chatRepo.GetChatRoomsByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("获取聊天房间失败: %v", err)
	}

	var response []ChatRoomResponse
	for _, room := range rooms {
		// 获取未读消息数量
		unreadCount, _ := s.chatRepo.GetUnreadCount(room.ID, userID)
		
		// 将字符串转换为数组
		participants := strings.Split(room.Participants, ",")
		
		roomResp := ChatRoomResponse{
			ID:           room.ID,
			Participants: participants,
			Type:         room.Type,
			UnreadCount:  int(unreadCount),
			CreatedAt:    room.CreatedAt,
			UpdatedAt:    room.UpdatedAt,
		}

		if room.Name != nil {
			roomResp.Name = *room.Name
		}
		if room.Avatar != nil {
			roomResp.Avatar = *room.Avatar
		}

		// 转换最后一条消息
		if room.LastMessage != nil {
			roomResp.LastMessage = s.convertToMessageResponse(room.LastMessage)
		}

		response = append(response, roomResp)
	}

	return response, nil
}

// 创建或获取聊天房间
func (s *chatService) CreateOrGetChatRoom(userID string, req *CreateChatRoomRequest) (*ChatRoomResponse, error) {
	if req.Type == "direct" {
		// 检查是否已存在直接聊天房间
		existingRoom, err := s.chatRepo.GetDirectChatRoom(userID, req.ParticipantID)
		if err == nil {
			// 房间已存在，返回现有房间
			return s.convertToChatRoomResponse(existingRoom, userID), nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("查询聊天房间失败: %v", err)
		}
	}

	// 创建新的聊天房间
	participants := userID + "," + req.ParticipantID
	room := &models.ChatRoom{
		Participants: participants,
		Type:         req.Type,
	}

	if req.Name != "" {
		room.Name = &req.Name
	}

	err := s.chatRepo.CreateChatRoom(room)
	if err != nil {
		return nil, fmt.Errorf("创建聊天房间失败: %v", err)
	}

	return s.convertToChatRoomResponse(room, userID), nil
}

// 获取聊天消息
func (s *chatService) GetMessages(userID string, chatID string, page, limit int) (*MessagesResponse, error) {
	// 验证用户是否有权限访问该聊天室
	room, err := s.chatRepo.GetChatRoomByID(chatID)
	if err != nil {
		return nil, fmt.Errorf("聊天房间不存在: %v", err)
	}

	// 检查用户是否有权限访问该聊天室
	participants := strings.Split(room.Participants, ",")
	hasPermission := false
	for _, participant := range participants {
		if participant == userID {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		return nil, errors.New("无权限访问该聊天室")
	}

	// 计算偏移量
	offset := (page - 1) * limit

	// 获取消息
	messages, total, err := s.chatRepo.GetMessagesByChatID(chatID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取消息失败: %v", err)
	}

	// 转换响应格式
	var messageResponses []MessageResponse
	for _, msg := range messages {
		messageResponses = append(messageResponses, *s.convertToMessageResponse(&msg))
	}

	// 计算是否还有更多消息
	hasMore := int64(offset+limit) < total

	return &MessagesResponse{
		Messages: messageResponses,
		Total:    int(total),
		Page:     page,
		Limit:    limit,
		HasMore:  hasMore,
	}, nil
}

// 发送消息
func (s *chatService) SendMessage(userID string, req *SendMessageRequest) (*MessageResponse, error) {
	// 验证用户是否有权限发送消息到该聊天室
	room, err := s.chatRepo.GetChatRoomByID(req.ChatID)
	if err != nil {
		return nil, fmt.Errorf("聊天房间不存在: %v", err)
	}

	participants := strings.Split(room.Participants, ",")
	hasPermission := false
	for _, participant := range participants {
		if participant == userID {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		return nil, errors.New("无权限发送消息到该聊天室")
	}

	// 创建消息
	message := &models.Message{
		ChatID:    req.ChatID,
		SenderID:  userID,
		Type:      req.Type,
		Timestamp: time.Now(),
		Status:    "sent",
	}

	// 设置消息内容
	if req.Content.Text != "" {
		message.Content.Text = &req.Content.Text
	}
	if req.Content.Media != nil {
		message.Content.Media = &models.MediaFile{
			URL:  req.Content.Media.URL,
			Type: req.Content.Media.Type,
			Size: req.Content.Media.Size,
		}
		if req.Content.Media.Duration > 0 {
			message.Content.Media.Duration = &req.Content.Media.Duration
		}
		if req.Content.Media.Thumbnail != "" {
			message.Content.Media.Thumbnail = &req.Content.Media.Thumbnail
		}
	}
	if req.Content.GameData != nil {
		message.Content.GameData = &models.GameSession{
			GameType: req.Content.GameData.GameType,
			Status:   req.Content.GameData.Status,
			Data:     req.Content.GameData.Data,
		}
	}

	if req.ReplyTo != "" {
		message.ReplyTo = &req.ReplyTo
	}

	// 保存消息
	err = s.chatRepo.CreateMessage(message)
	if err != nil {
		return nil, fmt.Errorf("发送消息失败: %v", err)
	}

	// 更新聊天房间的最后消息时间
	room.UpdatedAt = time.Now()
	s.chatRepo.UpdateChatRoom(room)

	// 为其他参与者创建消息状态记录
	for _, participantID := range participants {
		if participantID != userID {
			status := &models.MessageStatus{
				MessageID: message.ID,
				UserID:    participantID,
				Status:    "delivered",
				Timestamp: time.Now(),
			}
			s.chatRepo.CreateMessageStatus(status)
		}
	}

	// 通过WebSocket实时推送新消息给其他参与者
	messageResponse := s.convertToMessageResponse(message)
	wsMessage := types.WSMessage{
		Type:   types.WSMessageTypeNewMessage,
		Data:   messageResponse,
		ChatID: req.ChatID,
	}
	
	if msgBytes, err := json.Marshal(wsMessage); err == nil {
		for _, participantID := range participants {
			if participantID != userID && s.hub != nil {
				s.hub.SendToUser(participantID, msgBytes)
			}
		}
	}

	return messageResponse, nil
}

// 标记消息为已读
func (s *chatService) MarkAsRead(userID string, req *MarkAsReadRequest) error {
	// 验证用户权限
	room, err := s.chatRepo.GetChatRoomByID(req.ChatID)
	if err != nil {
		return fmt.Errorf("聊天房间不存在: %v", err)
	}

	participants := strings.Split(room.Participants, ",")
	hasPermission := false
	for _, participant := range participants {
		if participant == userID {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		return errors.New("无权限操作该聊天室")
	}

	// 查找或创建消息状态
	status, err := s.chatRepo.GetMessageStatus(req.MessageID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新的状态记录
			status = &models.MessageStatus{
				MessageID: req.MessageID,
				UserID:    userID,
				Status:    "read",
				Timestamp: time.Now(),
			}
			return s.chatRepo.CreateMessageStatus(status)
		}
		return fmt.Errorf("获取消息状态失败: %v", err)
	}

	// 更新状态
	status.Status = "read"
	status.Timestamp = time.Now()
	return s.chatRepo.UpdateMessageStatus(status)
}

// 删除消息
func (s *chatService) DeleteMessage(userID string, messageID string) error {
	// 验证消息是否存在且用户有权限删除
	message, err := s.chatRepo.GetMessageByID(messageID)
	if err != nil {
		return fmt.Errorf("消息不存在: %v", err)
	}

	if message.SenderID != userID {
		return errors.New("只能删除自己发送的消息")
	}

	return s.chatRepo.DeleteMessage(messageID)
}

// 上传媒体文件
func (s *chatService) UploadMedia(userID string, file *multipart.FileHeader, mediaType string) (*MediaUploadResponse, error) {
	// 验证文件类型和大小
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".mp4": true, ".mov": true, ".avi": true,
		".mp3": true, ".wav": true, ".aac": true,
	}

	if !allowedExts[ext] {
		return nil, errors.New("不支持的文件类型")
	}

	// 检查文件大小限制
	maxSize := int64(50 * 1024 * 1024) // 50MB
	if file.Size > maxSize {
		return nil, errors.New("文件大小超过限制")
	}

	// 模拟文件上传（实际应该上传到云存储）
	filename := fmt.Sprintf("%s_%d%s", userID, time.Now().Unix(), ext)
	url := fmt.Sprintf("https://example.com/chat-media/%s", filename)

	response := &MediaUploadResponse{
		URL:  url,
		Type: mediaType,
		Size: file.Size,
	}

	// 如果是视频或音频，可以设置时长（这里模拟）
	if mediaType == "video" || mediaType == "audio" {
		response.Duration = 30 // 模拟30秒
	}

	// 如果是视频，可以生成缩略图（这里模拟）
	if mediaType == "video" {
		response.Thumbnail = fmt.Sprintf("https://example.com/chat-thumbnails/%s_thumb.jpg", filename)
	}

	return response, nil
}

// 辅助方法：转换聊天房间为响应格式
func (s *chatService) convertToChatRoomResponse(room *models.ChatRoom, userID string) *ChatRoomResponse {
	unreadCount, _ := s.chatRepo.GetUnreadCount(room.ID, userID)
	
	// 将字符串转换为数组
	participants := strings.Split(room.Participants, ",")
	
	response := &ChatRoomResponse{
		ID:           room.ID,
		Participants: participants,
		Type:         room.Type,
		UnreadCount:  int(unreadCount),
		CreatedAt:    room.CreatedAt,
		UpdatedAt:    room.UpdatedAt,
	}

	if room.Name != nil {
		response.Name = *room.Name
	}
	if room.Avatar != nil {
		response.Avatar = *room.Avatar
	}
	if room.LastMessage != nil {
		response.LastMessage = s.convertToMessageResponse(room.LastMessage)
	}

	return response
}

// 辅助方法：转换消息为响应格式
func (s *chatService) convertToMessageResponse(message *models.Message) *MessageResponse {
	response := &MessageResponse{
		ID:        message.ID,
		ChatID:    message.ChatID,
		SenderID:  message.SenderID,
		Type:      message.Type,
		Timestamp: message.Timestamp,
		Status:    message.Status,
		Content:   MessageContentResponse{},
	}

	if message.Content.Text != nil {
		response.Content.Text = *message.Content.Text
	}
	if message.Content.Media != nil {
		response.Content.Media = &MediaFileResponse{
			URL:  message.Content.Media.URL,
			Type: message.Content.Media.Type,
			Size: message.Content.Media.Size,
		}
		if message.Content.Media.Duration != nil {
			response.Content.Media.Duration = *message.Content.Media.Duration
		}
		if message.Content.Media.Thumbnail != nil {
			response.Content.Media.Thumbnail = *message.Content.Media.Thumbnail
		}
	}
	if message.Content.GameData != nil {
		response.Content.GameData = &GameDataResponse{
			GameType: message.Content.GameData.GameType,
			Status:   message.Content.GameData.Status,
			Data:     message.Content.GameData.Data,
		}
	}

	if message.ReplyTo != nil {
		response.ReplyTo = *message.ReplyTo
	}

	return response
}