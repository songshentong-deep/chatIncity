package repository

import (
	"context"
	"fmt"
	"social-app/shared/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PhotoRepository interface {
	CreatePhoto(photo *models.Photo) error
	GetPhotoByID(id primitive.ObjectID) (*models.Photo, error)
	GetUserPhotos(userID string, photoType models.PhotoType) ([]*models.Photo, error)
	GetActiveUserAvatar(userID string) (*models.Photo, error)
	UpdateThumbnailURL(id primitive.ObjectID, thumbnailURL string) error
	SetPhotoActive(id primitive.ObjectID, active bool) error
	DeactivateUserPhotos(userID string, photoType models.PhotoType) error
	DeletePhoto(id primitive.ObjectID) error
	GetPhotosByIDs(ids []primitive.ObjectID) ([]*models.Photo, error)
}

type photoRepository struct {
	collection *mongo.Collection
}

func NewPhotoRepository(db *mongo.Database) PhotoRepository {
	collection := db.Collection("photos")
	
	// 创建索引
	ctx := context.Background()
	
	// 用户ID索引
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	})
	
	// 用户ID和类型复合索引
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "type", Value: 1},
		},
	})
	
	// 活跃状态索引
	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "type", Value: 1},
			{Key: "is_active", Value: 1},
		},
	})

	return &photoRepository{
		collection: collection,
	}
}

// CreatePhoto 创建照片记录
func (r *photoRepository) CreatePhoto(photo *models.Photo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, photo)
	if err != nil {
		return fmt.Errorf("插入照片记录失败: %v", err)
	}

	photo.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetPhotoByID 根据ID获取照片
func (r *photoRepository) GetPhotoByID(id primitive.ObjectID) (*models.Photo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var photo models.Photo
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&photo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("照片不存在")
		}
		return nil, fmt.Errorf("查询照片失败: %v", err)
	}

	return &photo, nil
}

// GetUserPhotos 获取用户照片
func (r *photoRepository) GetUserPhotos(userID string, photoType models.PhotoType) ([]*models.Photo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": userID,
		"type":    photoType,
	}

	// 按上传时间倒序排列
	opts := options.Find().SetSort(bson.D{{Key: "uploaded_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询用户照片失败: %v", err)
	}
	defer cursor.Close(ctx)

	var photos []*models.Photo
	for cursor.Next(ctx) {
		var photo models.Photo
		if err := cursor.Decode(&photo); err != nil {
			continue
		}
		photos = append(photos, &photo)
	}

	return photos, nil
}

// GetActiveUserAvatar 获取用户活跃头像
func (r *photoRepository) GetActiveUserAvatar(userID string) (*models.Photo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id":   userID,
		"type":      models.PhotoTypeAvatar,
		"is_active": true,
	}

	var photo models.Photo
	err := r.collection.FindOne(ctx, filter).Decode(&photo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("用户没有活跃头像")
		}
		return nil, fmt.Errorf("查询活跃头像失败: %v", err)
	}

	return &photo, nil
}

// UpdateThumbnailURL 更新缩略图URL
func (r *photoRepository) UpdateThumbnailURL(id primitive.ObjectID, thumbnailURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"thumbnail_url": thumbnailURL,
			"updated_at":    time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新缩略图URL失败: %v", err)
	}

	return nil
}

// SetPhotoActive 设置照片活跃状态
func (r *photoRepository) SetPhotoActive(id primitive.ObjectID, active bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"is_active":  active,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("设置照片活跃状态失败: %v", err)
	}

	return nil
}

// DeactivateUserPhotos 将用户指定类型的所有照片设为非活跃状态
func (r *photoRepository) DeactivateUserPhotos(userID string, photoType models.PhotoType) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": userID,
		"type":    photoType,
	}
	update := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("批量更新照片状态失败: %v", err)
	}

	return nil
}

// DeletePhoto 删除照片
func (r *photoRepository) DeletePhoto(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除照片记录失败: %v", err)
	}

	return nil
}

// GetPhotosByIDs 根据ID列表获取照片
func (r *photoRepository) GetPhotosByIDs(ids []primitive.ObjectID) ([]*models.Photo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": bson.M{"$in": ids},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("批量查询照片失败: %v", err)
	}
	defer cursor.Close(ctx)

	var photos []*models.Photo
	for cursor.Next(ctx) {
		var photo models.Photo
		if err := cursor.Decode(&photo); err != nil {
			continue
		}
		photos = append(photos, &photo)
	}

	return photos, nil
}