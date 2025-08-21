package service

import (
	"fmt"
	"social-app/shared/config"
	"social-app/shared/models"
	"social-app/services/chat-service/repository"
	"time"
)

type ChatService interface {
	// 聊天室管理
	CreateChatRoom(req *CreateChatRoomRequest, creatorID string) (*ChatRoomResponse, error)
	GetChatRooms(userID string) ([]ChatRoomResponse, error)
	GetChatRoom(chatRoomID, userID string) (*ChatRoomResponse, error)
	GetOrCreatePrivateChat(user1ID, user2ID string) (*ChatRoomResponse, error)
	GetGroupMembers(userID, chatRoomID string) ([]MemberResponse, error)
	
	// 消息管理
	SendMessage(senderID string, req *SendMessageRequest) (*MessageResponse, error)
	GetMessages(userID, chatRoomID string, page, limit int) ([]MessageResponse, error)
	MarkAsRead(userID string, req *MarkAsReadRequest) error
	DeleteMessage(userID, messageID string) error
}

type chatService struct {
	chatRepo repository.ChatRepository
	config   *config.Config
}

func NewChatService(chatRepo repository.ChatRepository, cfg *config.Config) ChatService {
	return &chatService{
		chatRepo: chatRepo,
		config:   cfg,
	}
}

// 创建聊天室
func (s *chatService) CreateChatRoom(req *CreateChatRoomRequest, creatorID string) (*ChatRoomResponse, error) {
	// 创建聊天室
	chatRoom := &models.ChatRoom{
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   creatorID,
	}

	if err := s.chatRepo.CreateChatRoom(chatRoom); err != nil {
		return nil, fmt.Errorf("创建聊天室失败: %v", err)
	}

	// 添加创建者为管理员
	creatorMember := &models.ChatRoomMember{
		ChatRoomID: chatRoom.ID,
		UserID:     creatorID,
		Role:       models.MemberRoleAdmin,
		JoinedAt:   time.Now(),
	}

	if err := s.chatRepo.AddChatRoomMember(creatorMember); err != nil {
		return nil, fmt.Errorf("添加创建者失败: %v", err)
	}

	// 添加其他成员
	for _, memberID := range req.MemberIDs {
		if memberID != creatorID {
			member := &models.ChatRoomMember{
				ChatRoomID: chatRoom.ID,
				UserID:     memberID,
				Role:       models.MemberRoleMember,
				JoinedAt:   time.Now(),
			}
			if err := s.chatRepo.AddChatRoomMember(member); err != nil {
				return nil, fmt.Errorf("添加成员失败: %v", err)
			}
		}
	}

	return s.buildChatRoomResponse(chatRoom, creatorID)
}

// 获取用户的聊天室列表
func (s *chatService) GetChatRooms(userID string) ([]ChatRoomResponse, error) {
	chatRooms, err := s.chatRepo.GetChatRoomsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取聊天室列表失败: %v", err)
	}

	var responses []ChatRoomResponse
	for _, chatRoom := range chatRooms {
		response, err := s.buildChatRoomResponse(chatRoom, userID)
		if err != nil {
			continue // 跳过有错误的聊天室
		}
		responses = append(responses, *response)
	}

	return responses, nil
}

// 获取单个聊天室信息
func (s *chatService) GetChatRoom(chatRoomID, userID string) (*ChatRoomResponse, error) {
	// 检查用户是否是聊天室成员
	isMember, err := s.chatRepo.IsChatRoomMember(chatRoomID, userID)
	if err != nil {
		return nil, fmt.Errorf("检查成员身份失败: %v", err)
	}
	if !isMember {
		return nil, fmt.Errorf("您不是该聊天室的成员")
	}

	chatRoom, err := s.chatRepo.GetChatRoomByID(chatRoomID)
	if err != nil {
		return nil, fmt.Errorf("获取聊天室失败: %v", err)
	}

	return s.buildChatRoomResponse(chatRoom, userID)
}

// 获取或创建私聊
func (s *chatService) GetOrCreatePrivateChat(user1ID, user2ID string) (*ChatRoomResponse, error) {
	// 尝试获取现有的私聊
	chatRoom, err := s.chatRepo.GetPrivateChatRoom(user1ID, user2ID)
	if err == nil {
		return s.buildChatRoomResponse(chatRoom, user1ID)
	}

	// 创建新的私聊
	req := &CreateChatRoomRequest{
		Type:      models.ChatRoomTypePrivate,
		MemberIDs: []string{user1ID, user2ID},
	}

	return s.CreateChatRoom(req, user1ID)
}

// 发送消息
func (s *chatService) SendMessage(senderID string, req *SendMessageRequest) (*MessageResponse, error) {
	// 检查用户是否是聊天室成员
	isMember, err := s.chatRepo.IsChatRoomMember(req.ChatRoomID, senderID)
	if err != nil {
		return nil, fmt.Errorf("检查成员身份失败: %v", err)
	}
	if !isMember {
		return nil, fmt.Errorf("您不是该聊天室的成员")
	}

	// 创建消息
	messageType := req.MessageType
	if messageType == "" {
		messageType = models.MessageTypeText
	}
	
	message := &models.Message{
		ChatRoomID:  req.ChatRoomID,
		SenderID:    senderID,
		Content:     req.Content,
		MessageType: messageType,
		Status:      models.MessageStatusSent,
	}

	if err := s.chatRepo.CreateMessage(message); err != nil {
		return nil, fmt.Errorf("发送消息失败: %v", err)
	}

	// 获取完整的消息信息
	fullMessage, err := s.chatRepo.GetMessageByID(message.ID)
	if err != nil {
		return nil, fmt.Errorf("获取消息失败: %v", err)
	}

	return s.buildMessageResponse(fullMessage), nil
}

// 获取消息列表
func (s *chatService) GetMessages(userID, chatRoomID string, page, limit int) ([]MessageResponse, error) {
	offset := (page - 1) * limit
	// 检查用户是否是聊天室成员
	isMember, err := s.chatRepo.IsChatRoomMember(chatRoomID, userID)
	if err != nil {
		return nil, fmt.Errorf("检查成员身份失败: %v", err)
	}
	if !isMember {
		return nil, fmt.Errorf("您不是该聊天室的成员")
	}

	messages, err := s.chatRepo.GetMessagesByRoomID(chatRoomID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("获取消息列表失败: %v", err)
	}

	var responses []MessageResponse
	for _, message := range messages {
		responses = append(responses, *s.buildMessageResponse(message))
	}

	return responses, nil
}

// 标记消息为已读
func (s *chatService) MarkAsRead(userID string, req *MarkAsReadRequest) error {
	// 检查用户是否是聊天室成员
	isMember, err := s.chatRepo.IsChatRoomMember(req.ChatRoomID, userID)
	if err != nil {
		return fmt.Errorf("检查成员身份失败: %v", err)
	}
	if !isMember {
		return fmt.Errorf("您不是该聊天室的成员")
	}

	return s.chatRepo.UpdateLastReadAt(req.ChatRoomID, userID)
}

// 删除消息
func (s *chatService) DeleteMessage(userID, messageID string) error {
	// 获取消息信息
	message, err := s.chatRepo.GetMessageByID(messageID)
	if err != nil {
		return fmt.Errorf("获取消息失败: %v", err)
	}

	// 检查用户是否有权限删除消息（只能删除自己的消息）
	if message.SenderID != userID {
		return fmt.Errorf("您只能删除自己的消息")
	}

	// 删除消息
	return s.chatRepo.DeleteMessage(messageID)
}

// 获取群成员列表
func (s *chatService) GetGroupMembers(userID, chatRoomID string) ([]MemberResponse, error) {
	// 检查用户是否是聊天室成员
	isMember, err := s.chatRepo.IsChatRoomMember(chatRoomID, userID)
	if err != nil {
		return nil, fmt.Errorf("检查成员身份失败: %v", err)
	}
	if !isMember {
		return nil, fmt.Errorf("您不是该聊天室的成员")
	}

	// 获取聊天室成员
	members, err := s.chatRepo.GetChatRoomMembers(chatRoomID)
	if err != nil {
		return nil, fmt.Errorf("获取群成员失败: %v", err)
	}

	var memberResponses []MemberResponse
	for _, member := range members {
		nickname := "未知用户"
		avatar := ""
		
		if member.User.Profile.Nickname != "" {
			nickname = member.User.Profile.Nickname
		}
		if member.User.Profile.Avatar != "" {
			avatar = member.User.Profile.Avatar
		}
		
		memberResponse := MemberResponse{
			ID:       member.ID,
			UserID:   member.UserID,
			Nickname: nickname,
			Avatar:   avatar,
			Role:     string(member.Role),
			JoinedAt: member.JoinedAt.Format("2006-01-02 15:04:05"),
		}
		memberResponses = append(memberResponses, memberResponse)
	}

	return memberResponses, nil
}

// 构建聊天室响应
func (s *chatService) buildChatRoomResponse(chatRoom *models.ChatRoom, userID string) (*ChatRoomResponse, error) {
	response := &ChatRoomResponse{
		ID:        chatRoom.ID,
		Type:      string(chatRoom.Type),
		CreatedBy: chatRoom.CreatedBy,
		CreatedAt: chatRoom.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 获取聊天室成员
	members, err := s.chatRepo.GetChatRoomMembers(chatRoom.ID)
	if err == nil {
		var memberResponses []MemberResponse
		for _, member := range members {
			nickname := "未知用户"
			avatar := ""
			
			if member.User.Profile.Nickname != "" {
				nickname = member.User.Profile.Nickname
			}
			if member.User.Profile.Avatar != "" {
				avatar = member.User.Profile.Avatar
			}
			
			memberResponse := MemberResponse{
				ID:       member.ID,
				UserID:   member.UserID,
				Nickname: nickname,
				Avatar:   avatar,
				Role:     string(member.Role),
				JoinedAt: member.JoinedAt.Format("2006-01-02 15:04:05"),
			}
			memberResponses = append(memberResponses, memberResponse)
		}
		response.Members = memberResponses
	}

	// 设置聊天室名称
	if chatRoom.Name != nil {
		response.Name = *chatRoom.Name
	} else if chatRoom.Type == models.ChatRoomTypePrivate {
		// 对于私聊，使用对方的昵称作为聊天室名称
		if err == nil {
			for _, member := range members {
				if member.UserID != userID {
					if member.User.Profile.Nickname != "" {
						response.Name = member.User.Profile.Nickname
					} else {
						response.Name = "好友"
					}
					break
				}
			}
		}
		if response.Name == "" {
			response.Name = "私聊"
		}
	}

	// 设置描述
	if chatRoom.Description != nil {
		response.Description = *chatRoom.Description
	}

	// 获取最后一条消息
	messages, err := s.chatRepo.GetMessagesByRoomID(chatRoom.ID, 1, 0)
	if err == nil && len(messages) > 0 {
		response.LastMessage = s.buildMessageResponse(messages[0])
	}

	// TODO: 计算未读消息数量
	response.UnreadCount = 0

	return response, nil
}

// 构建消息响应
func (s *chatService) buildMessageResponse(message *models.Message) *MessageResponse {
	senderName := "未知用户"
	senderAvatar := ""
	
	// 安全地访问发送者信息
	if message.Sender.ID != "" {
		if message.Sender.Profile.Nickname != "" {
			senderName = message.Sender.Profile.Nickname
		}
		if message.Sender.Profile.Avatar != "" {
			senderAvatar = message.Sender.Profile.Avatar
		}
	}
	
	return &MessageResponse{
		ID:           message.ID,
		ChatRoomID:   message.ChatRoomID,
		SenderID:     message.SenderID,
		SenderName:   senderName,
		SenderAvatar: senderAvatar,
		Content:      message.Content,
		MessageType:  string(message.MessageType),
		Status:       string(message.Status),
		CreatedAt:    message.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}