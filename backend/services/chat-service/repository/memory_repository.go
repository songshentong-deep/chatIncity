package repository

import (
	"social-app/shared/models"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MemoryChatRepository 内存实现的聊天仓库，用于测试
type MemoryChatRepository struct {
	chatRooms     map[string]*models.ChatRoom
	messages      map[string]*models.Message
	messageStatus map[string]*models.MessageStatus
	mutex         sync.RWMutex
}

// 聊天房间相关实现
func (r *MemoryChatRepository) CreateChatRoom(room *models.ChatRoom) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if room.ID == "" {
		room.ID = uuid.New().String()
	}
	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()

	r.chatRooms[room.ID] = room
	return nil
}

// 消息相关实现
func (r *MemoryChatRepository) CreateMessage(message *models.Message) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if message.ID == "" {
		message.ID = uuid.New().String()
	}

	r.messages[message.ID] = message
	return nil
}

func (r *MemoryChatRepository) UpdateMessage(message *models.Message) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.messages[message.ID]; !exists {
		return gorm.ErrRecordNotFound
	}

	r.messages[message.ID] = message
	return nil
}

func (r *MemoryChatRepository) DeleteMessage(messageID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.messages[messageID]; !exists {
		return gorm.ErrRecordNotFound
	}

	delete(r.messages, messageID)
	return nil
}

func (r *MemoryChatRepository) GetMessageStatus(messageID, userID string) (*models.MessageStatus, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	key := messageID + ":" + userID
	status, exists := r.messageStatus[key]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return status, nil
}

// 辅助函数
func contains(participants, userID string) bool {
	// 简单的字符串包含检查
	return participants == userID ||
		participants == userID+"," ||
		participants == ","+userID ||
		len(participants) > len(userID) && (participants[:len(userID)+1] == userID+"," || participants[len(participants)-len(userID)-1:] == ","+userID)
}
