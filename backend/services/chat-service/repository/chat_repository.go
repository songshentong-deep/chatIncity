package repository

import (
	"context"
	"social-app/shared/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type ChatRepository interface {
	// 聊天房间相关
	CreateChatRoom(room *models.ChatRoom) error
	GetChatRoomByID(id string) (*models.ChatRoom, error)
	GetChatRoomsByUser(userID string) ([]models.ChatRoom, error)
	GetDirectChatRoom(user1ID, user2ID string) (*models.ChatRoom, error)
	UpdateChatRoom(room *models.ChatRoom) error
	
	// 消息相关
	CreateMessage(message *models.Message) error
	GetMessageByID(id string) (*models.Message, error)
	GetMessagesByChatID(chatID string, offset, limit int) ([]models.Message, int64, error)
	UpdateMessage(message *models.Message) error
	DeleteMessage(messageID string) error
	
	// 消息状态相关
	CreateMessageStatus(status *models.MessageStatus) error
	GetMessageStatus(messageID, userID string) (*models.MessageStatus, error)
	UpdateMessageStatus(status *models.MessageStatus) error
	GetUnreadCount(chatID, userID string) (int64, error)
}

type chatRepository struct {
	db      *gorm.DB
	mongoDB *mongo.Database
}

func NewChatRepository(db *gorm.DB, mongoDB *mongo.Database) ChatRepository {
	return &chatRepository{
		db:      db,
		mongoDB: mongoDB,
	}
}

// 聊天房间相关实现
func (r *chatRepository) CreateChatRoom(room *models.ChatRoom) error {
	return r.db.Create(room).Error
}

func (r *chatRepository) GetChatRoomByID(id string) (*models.ChatRoom, error) {
	var room models.ChatRoom
	err := r.db.Preload("LastMessage").Where("id = ?", id).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *chatRepository) GetChatRoomsByUser(userID string) ([]models.ChatRoom, error) {
	var rooms []models.ChatRoom
	err := r.db.Preload("LastMessage").
		Where("participants LIKE ?", "%"+userID+"%").
		Order("updated_at DESC").
		Find(&rooms).Error
	return rooms, err
}

func (r *chatRepository) GetDirectChatRoom(user1ID, user2ID string) (*models.ChatRoom, error) {
	var room models.ChatRoom
	// 查找包含两个用户的直接聊天房间
	err := r.db.Where("type = ? AND (participants = ? OR participants = ?)", 
		"direct", 
		user1ID+","+user2ID, 
		user2ID+","+user1ID).
		First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *chatRepository) UpdateChatRoom(room *models.ChatRoom) error {
	return r.db.Save(room).Error
}

// 消息相关实现（使用MongoDB）
func (r *chatRepository) CreateMessage(message *models.Message) error {
	collection := r.mongoDB.Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 如果没有ID，生成一个
	if message.ID == "" {
		message.ID = primitive.NewObjectID().Hex()
	}

	_, err := collection.InsertOne(ctx, message)
	return err
}

func (r *chatRepository) GetMessageByID(id string) (*models.Message, error) {
	collection := r.mongoDB.Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var message models.Message
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *chatRepository) GetMessagesByChatID(chatID string, offset, limit int) ([]models.Message, int64, error) {
	collection := r.mongoDB.Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取总数
	total, err := collection.CountDocuments(ctx, bson.M{"chat_id": chatID})
	if err != nil {
		return nil, 0, err
	}

	// 获取消息列表（按时间倒序）
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"timestamp", -1}})
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, bson.M{"chat_id": chatID}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (r *chatRepository) UpdateMessage(message *models.Message) error {
	collection := r.mongoDB.Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"id": message.ID}
	update := bson.M{"$set": message}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *chatRepository) DeleteMessage(messageID string) error {
	collection := r.mongoDB.Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"id": messageID})
	return err
}

// 消息状态相关实现
func (r *chatRepository) CreateMessageStatus(status *models.MessageStatus) error {
	return r.db.Create(status).Error
}

func (r *chatRepository) GetMessageStatus(messageID, userID string) (*models.MessageStatus, error) {
	var status models.MessageStatus
	err := r.db.Where("message_id = ? AND user_id = ?", messageID, userID).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *chatRepository) UpdateMessageStatus(status *models.MessageStatus) error {
	return r.db.Save(status).Error
}

func (r *chatRepository) GetUnreadCount(chatID, userID string) (int64, error) {
	// 从PostgreSQL获取已读消息ID列表
	var readMessageIDs []string
	err := r.db.Model(&models.MessageStatus{}).
		Select("message_id").
		Where("user_id = ? AND status = ?", userID, "read").
		Pluck("message_id", &readMessageIDs).Error
	if err != nil {
		return 0, err
	}

	// 从MongoDB查询未读消息数量
	collection := r.mongoDB.Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"chat_id":   chatID,
		"sender_id": bson.M{"$ne": userID},
	}

	// 如果有已读消息，排除它们
	if len(readMessageIDs) > 0 {
		filter["id"] = bson.M{"$nin": readMessageIDs}
	}

	count, err := collection.CountDocuments(ctx, filter)
	return count, err
}