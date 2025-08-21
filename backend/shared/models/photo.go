package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Photo 照片模型 - 存储在MongoDB中
type Photo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string            `bson:"user_id" json:"user_id"`
	Type        PhotoType         `bson:"type" json:"type"` // avatar, gallery
	Filename    string            `bson:"filename" json:"filename"`
	OriginalName string           `bson:"original_name" json:"original_name"`
	ContentType string            `bson:"content_type" json:"content_type"`
	Size        int64             `bson:"size" json:"size"`
	Width       int               `bson:"width" json:"width"`
	Height      int               `bson:"height" json:"height"`
	URL         string            `bson:"url" json:"url"`
	ThumbnailURL string           `bson:"thumbnail_url,omitempty" json:"thumbnail_url,omitempty"`
	Data        []byte            `bson:"data" json:"-"` // 存储图片数据，JSON中不返回
	IsActive    bool              `bson:"is_active" json:"is_active"`
	UploadedAt  time.Time         `bson:"uploaded_at" json:"uploaded_at"`
	CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
	
	// 元数据
	Metadata PhotoMetadata `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

// PhotoType 照片类型
type PhotoType string

const (
	PhotoTypeAvatar  PhotoType = "avatar"
	PhotoTypeGallery PhotoType = "gallery"
)

// PhotoMetadata 照片元数据
type PhotoMetadata struct {
	Camera       string            `bson:"camera,omitempty" json:"camera,omitempty"`
	Location     *PhotoLocation    `bson:"location,omitempty" json:"location,omitempty"`
	Tags         []string          `bson:"tags,omitempty" json:"tags,omitempty"`
	Description  string            `bson:"description,omitempty" json:"description,omitempty"`
	ExifData     map[string]string `bson:"exif_data,omitempty" json:"exif_data,omitempty"`
}

// PhotoLocation 照片位置信息
type PhotoLocation struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
	Address   string  `bson:"address,omitempty" json:"address,omitempty"`
	City      string  `bson:"city,omitempty" json:"city,omitempty"`
	Country   string  `bson:"country,omitempty" json:"country,omitempty"`
}

// PhotoUploadRequest 照片上传请求
type PhotoUploadRequest struct {
	UserID      string            `json:"user_id"`
	Type        PhotoType         `json:"type"`
	File        []byte            `json:"file"`
	Filename    string            `json:"filename"`
	ContentType string            `json:"content_type"`
	Metadata    *PhotoMetadata    `json:"metadata,omitempty"`
}

// PhotoResponse 照片响应
type PhotoResponse struct {
	ID           string        `json:"id"`
	UserID       string        `json:"user_id"`
	Type         string        `json:"type"`
	URL          string        `json:"url"`
	ThumbnailURL string        `json:"thumbnail_url,omitempty"`
	Width        int           `json:"width"`
	Height       int           `json:"height"`
	Size         int64         `json:"size"`
	IsActive     bool          `json:"is_active"`
	UploadedAt   string        `json:"uploaded_at"`
	Metadata     PhotoMetadata `json:"metadata,omitempty"`
}

// ToResponse 转换为响应格式
func (p *Photo) ToResponse() *PhotoResponse {
	return &PhotoResponse{
		ID:           p.ID.Hex(),
		UserID:       p.UserID,
		Type:         string(p.Type),
		URL:          p.URL,
		ThumbnailURL: p.ThumbnailURL,
		Width:        p.Width,
		Height:       p.Height,
		Size:         p.Size,
		IsActive:     p.IsActive,
		UploadedAt:   p.UploadedAt.Format("2006-01-02 15:04:05"),
		Metadata:     p.Metadata,
	}
}