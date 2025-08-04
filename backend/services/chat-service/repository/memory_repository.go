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
	chatRooms      map[string]*models.ChatRoom
	messages       map[string]*models.Message
	messageStatus  map[string]*models.MessageStatus
	mutex          sync.RWMutex
}

func NewMemoryChatRepository() ChatRepository {
	return &MemoryChatRepository{
		chatRooms:     make(map[string]*models.ChatRoom),
		messages:      make(map[string]*models.Message),
		messageStatus: make(map[string]*models.MessageStatus),
	}
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

func (r *MemoryChatRepository) GetChatRoomByID(id string) (*models.ChatRoom, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	room, exists := r.chatRooms[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return room, nil
}

func (r *MemoryChatRepository) GetChatRoomsByUser(userID string) ([]models.ChatRoom, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var rooms []models.ChatRoom
	for _, room := range r.chatRooms {
		// 简单检查用户是否在参与者列表中
		if contains(room.Participants, userID) {
			rooms = append(rooms, *room)
		}
	}
	return rooms, nil
}

func (r *MemoryChatRepository) GetDirectChatRoom(user1ID, user2ID string) (*models.ChatRoom, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for _, room := range r.chatRooms {
		if room.Type == "direct" {
			participants := room.Participants
			if (participants == user1ID+","+user2ID) || (participants == user2ID+","+user1ID) {
				return room, nil
			}
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *MemoryChatRepository) UpdateChatRoom(room *models.ChatRoom) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.chatRooms[room.ID]; !exists {
		return gorm.ErrRecordNotFound
	}
	
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

func (r *MemoryChatRepository) GetMessageByID(id string) (*models.Message, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	message, exists := r.messages[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return message, nil
}

func (r *MemoryChatRepository) GetMessagesByChatID(chatID string, offset, limit int) ([]models.Message, int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var messages []models.Message
	for _, message := range r.messages {
		if message.ChatID == chatID {
			messages = append(messages, *message)
		}
	}
	
	// 按时间排序（简化实现）
	total := int64(len(messages))
	
	// 简单分页
	start := offset
	end := offset + limit
	if start > len(messages) {
		start = len(messages)
	}
	if end > len(messages) {
		end = len(messages)
	}
	
	if start < end {
		messages = messages[start:end]
	} else {
		messages = []models.Message{}
	}
	
	return messages, total, nil
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

// 消息状态相关实现
func (r *MemoryChatRepository) CreateMessageStatus(status *models.MessageStatus) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if status.ID == "" {
		status.ID = uuid.New().String()
	}
	
	key := status.MessageID + ":" + status.UserID
	r.messageStatus[key] = status
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

func (r *MemoryChatRepository) UpdateMessageStatus(status *models.MessageStatus) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	key := status.MessageID + ":" + status.UserID
	if _, exists := r.messageStatus[key]; !exists {
		return gorm.ErrRecordNotFound
	}
	
	r.messageStatus[key] = status
	return nil
}

func (r *MemoryChatRepository) GetUnreadCount(chatID, userID string) (int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	// 获取该聊天室的所有消息
	var unreadCount int64
	for _, message := range r.messages {
		if message.ChatID == chatID && message.SenderID != userID {
			// 检查是否已读
			key := message.ID + ":" + userID
			if status, exists := r.messageStatus[key]; !exists || status.Status != "read" {
				unreadCount++
			}
		}
	}
	
	return unreadCount, nil
}

// 辅助函数
func contains(participants, userID string) bool {
	// 简单的字符串包含检查
	return participants == userID || 
		   participants == userID+"," ||
		   participants == ","+userID ||
		   len(participants) > len(userID) && (participants[:len(userID)+1] == userID+"," || participants[len(participants)-len(userID)-1:] == ","+userID)
}