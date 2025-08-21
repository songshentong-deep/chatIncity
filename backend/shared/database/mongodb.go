package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient MongoDB客户端
type MongoClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// ConnectMongoDB 连接到MongoDB
func ConnectMongoDB(uri, database string) (*MongoClient, error) {
	// 设置连接选项
	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetMaxPoolSize(10)
	clientOptions.SetConnectTimeout(10 * time.Second)

	// 连接到MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	db := client.Database(database)
	
	log.Printf("MongoDB连接成功: %s/%s", uri, database)
	
	return &MongoClient{
		Client:   client,
		Database: db,
	}, nil
}

// Close 关闭MongoDB连接
func (mc *MongoClient) Close() error {
	return mc.Client.Disconnect(context.Background())
}

// GetCollection 获取集合
func (mc *MongoClient) GetCollection(name string) *mongo.Collection {
	return mc.Database.Collection(name)
}