package repository

import (
	"social-app/shared/models"
	"gorm.io/gorm"
)

type ChatRepository interface {
	// 聊天室管理
	CreateChatRoom(chatRoom *models.ChatRoom) error
	GetChatRoomByID(chatRoomID string) (*models.ChatRoom, error)
	GetChatRoomsByUserID(userID string) ([]*models.ChatRoom, error)
	GetPrivateChatRoom(user1ID, user2ID string) (*models.ChatRoom, error)
	
	// 聊天室成员管理
	AddChatRoomMember(member *models.ChatRoomMember) error
	GetChatRoomMembers(chatRoomID string) ([]*models.ChatRoomMember, error)
	IsChatRoomMember(chatRoomID, userID string) (bool, error)
	UpdateLastReadAt(chatRoomID, userID string) error
	
	// 消息管理
	CreateMessage(message *models.Message) error
	GetMessagesByRoomID(chatRoomID string, limit, offset int) ([]*models.Message, error)
	GetMessageByID(messageID string) (*models.Message, error)
	UpdateMessageStatus(messageID string, status models.MessageStatus) error
	DeleteMessage(messageID string) error
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

// 聊天室管理
func (r *chatRepository) CreateChatRoom(chatRoom *models.ChatRoom) error {
	return r.db.Create(chatRoom).Error
}

func (r *chatRepository) GetChatRoomByID(chatRoomID string) (*models.ChatRoom, error) {
	var chatRoom models.ChatRoom
	err := r.db.Preload("Creator").Where("id = ?", chatRoomID).First(&chatRoom).Error
	if err != nil {
		return nil, err
	}
	return &chatRoom, nil
}

func (r *chatRepository) GetChatRoomsByUserID(userID string) ([]*models.ChatRoom, error) {
	var chatRooms []*models.ChatRoom
	err := r.db.Joins("JOIN chat_room_members ON chat_rooms.id = chat_room_members.chat_room_id").
		Where("chat_room_members.user_id = ?", userID).
		Preload("Creator").
		Find(&chatRooms).Error
	return chatRooms, err
}

func (r *chatRepository) GetPrivateChatRoom(user1ID, user2ID string) (*models.ChatRoom, error) {
	var chatRoom models.ChatRoom
	
	// 查找同时包含两个用户的私聊聊天室
	err := r.db.Where("type = ? AND id IN (?)", models.ChatRoomTypePrivate,
		r.db.Table("chat_room_members").
			Select("chat_room_id").
			Where("user_id IN (?, ?)", user1ID, user2ID).
			Group("chat_room_id").
			Having("COUNT(DISTINCT user_id) = 2 AND COUNT(*) = 2")).
		Preload("Creator").
		First(&chatRoom).Error
		
	if err != nil {
		return nil, err
	}
	return &chatRoom, nil
}

// 聊天室成员管理
func (r *chatRepository) AddChatRoomMember(member *models.ChatRoomMember) error {
	return r.db.Create(member).Error
}

func (r *chatRepository) GetChatRoomMembers(chatRoomID string) ([]*models.ChatRoomMember, error) {
	var members []*models.ChatRoomMember
	err := r.db.Preload("User").Where("chat_room_id = ?", chatRoomID).Find(&members).Error
	return members, err
}

func (r *chatRepository) IsChatRoomMember(chatRoomID, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.ChatRoomMember{}).
		Where("chat_room_id = ? AND user_id = ?", chatRoomID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *chatRepository) UpdateLastReadAt(chatRoomID, userID string) error {
	return r.db.Model(&models.ChatRoomMember{}).
		Where("chat_room_id = ? AND user_id = ?", chatRoomID, userID).
		Update("last_read_at", gorm.Expr("NOW()")).Error
}

// 消息管理
func (r *chatRepository) CreateMessage(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *chatRepository) GetMessagesByRoomID(chatRoomID string, limit, offset int) ([]*models.Message, error) {
	var messages []*models.Message
	err := r.db.Preload("Sender").
		Where("chat_room_id = ?", chatRoomID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (r *chatRepository) GetMessageByID(messageID string) (*models.Message, error) {
	var message models.Message
	err := r.db.Preload("Sender").Where("id = ?", messageID).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *chatRepository) UpdateMessageStatus(messageID string, status models.MessageStatus) error {
	return r.db.Model(&models.Message{}).Where("id = ?", messageID).Update("status", status).Error
}

func (r *chatRepository) DeleteMessage(messageID string) error {
	return r.db.Delete(&models.Message{}, "id = ?", messageID).Error
}