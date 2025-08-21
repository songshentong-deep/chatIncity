package repository

import (
	"context"
	"social-app/services/news-service/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NewsRepository interface {
	Create(news *models.News) error
	GetByID(id string) (*models.News, error)
	GetAll(page, limit int) ([]*models.News, error)
	GetByCategory(category models.NewsCategory, page, limit int) ([]*models.News, error)
	Update(news *models.News) error
	Delete(id string) error
	GetLatest(limit int) ([]*models.News, error)
	GetCategoryCount(category models.NewsCategory) (int64, error)
	IncrementViewCount(id string) error
	CheckExists(sourceURL string) (bool, error)
	DeleteOldNews(days int) error
}

type newsRepository struct {
	collection *mongo.Collection
}

func NewNewsRepository(db *mongo.Database) NewsRepository {
	return &newsRepository{
		collection: db.Collection("news"),
	}
}

func (r *newsRepository) Create(news *models.News) error {
	news.CreatedAt = time.Now()
	news.UpdatedAt = time.Now()
	news.IsActive = true
	
	_, err := r.collection.InsertOne(context.Background(), news)
	return err
}

func (r *newsRepository) GetByID(id string) (*models.News, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var news models.News
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&news)
	if err != nil {
		return nil, err
	}

	return &news, nil
}

func (r *newsRepository) GetAll(page, limit int) ([]*models.News, error) {
	skip := (page - 1) * limit
	
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "published_at", Value: -1}})

	cursor, err := r.collection.Find(context.Background(), bson.M{"is_active": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var news []*models.News
	for cursor.Next(context.Background()) {
		var n models.News
		if err := cursor.Decode(&n); err != nil {
			continue
		}
		news = append(news, &n)
	}

	return news, nil
}

func (r *newsRepository) GetByCategory(category models.NewsCategory, page, limit int) ([]*models.News, error) {
	skip := (page - 1) * limit
	
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "published_at", Value: -1}})

	filter := bson.M{
		"category":  category,
		"is_active": true,
	}

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var news []*models.News
	for cursor.Next(context.Background()) {
		var n models.News
		if err := cursor.Decode(&n); err != nil {
			continue
		}
		news = append(news, &n)
	}

	return news, nil
}

func (r *newsRepository) Update(news *models.News) error {
	news.UpdatedAt = time.Now()
	
	filter := bson.M{"_id": news.ID}
	update := bson.M{"$set": news}
	
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *newsRepository) Delete(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}

func (r *newsRepository) GetLatest(limit int) ([]*models.News, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "published_at", Value: -1}})

	cursor, err := r.collection.Find(context.Background(), bson.M{"is_active": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var news []*models.News
	for cursor.Next(context.Background()) {
		var n models.News
		if err := cursor.Decode(&n); err != nil {
			continue
		}
		news = append(news, &n)
	}

	return news, nil
}

func (r *newsRepository) GetCategoryCount(category models.NewsCategory) (int64, error) {
	filter := bson.M{
		"category":  category,
		"is_active": true,
	}
	
	return r.collection.CountDocuments(context.Background(), filter)
}

func (r *newsRepository) IncrementViewCount(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"view_count": 1}}
	
	_, err = r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *newsRepository) CheckExists(sourceURL string) (bool, error) {
	count, err := r.collection.CountDocuments(context.Background(), bson.M{"source_url": sourceURL})
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

func (r *newsRepository) DeleteOldNews(days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	
	filter := bson.M{
		"created_at": bson.M{"$lt": cutoffDate},
	}
	
	_, err := r.collection.DeleteMany(context.Background(), filter)
	return err
}