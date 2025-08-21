package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 新闻分类
type NewsCategory string

const (
	CategoryTech     NewsCategory = "tech"     // 科技
	CategorySports   NewsCategory = "sports"   // 体育
	CategoryPolitics NewsCategory = "politics" // 政治
	CategoryMilitary NewsCategory = "military" // 军事
	CategoryEntertainment NewsCategory = "entertainment" // 娱乐
	CategoryBusiness NewsCategory = "business" // 商业
	CategoryHealth   NewsCategory = "health"   // 健康
	CategoryGeneral  NewsCategory = "general"  // 综合
)

// 新闻模型
type News struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string            `json:"title" bson:"title"`
	Description string            `json:"description" bson:"description"`
	Content     string            `json:"content" bson:"content"`
	ImageURL    string            `json:"image_url" bson:"image_url"`
	SourceURL   string            `json:"source_url" bson:"source_url"`
	Source      string            `json:"source" bson:"source"`
	Author      string            `json:"author" bson:"author"`
	Category    NewsCategory      `json:"category" bson:"category"`
	Tags        []string          `json:"tags" bson:"tags"`
	PublishedAt time.Time         `json:"published_at" bson:"published_at"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
	ViewCount   int64             `json:"view_count" bson:"view_count"`
	IsActive    bool              `json:"is_active" bson:"is_active"`
}

// 新闻源配置
type NewsSource struct {
	ID       string       `json:"id" bson:"_id"`
	Name     string       `json:"name" bson:"name"`
	URL      string       `json:"url" bson:"url"`
	Category NewsCategory `json:"category" bson:"category"`
	APIKey   string       `json:"api_key" bson:"api_key"`
	IsActive bool         `json:"is_active" bson:"is_active"`
}

// 新闻API响应结构
type NewsAPIResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"source"`
		Author      string    `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		URL         string    `json:"url"`
		URLToImage  string    `json:"urlToImage"`
		PublishedAt time.Time `json:"publishedAt"`
		Content     string    `json:"content"`
	} `json:"articles"`
}

// 分类信息
type CategoryInfo struct {
	Category    NewsCategory `json:"category"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Count       int64        `json:"count"`
}

// 获取分类中文名称
func (c NewsCategory) GetChineseName() string {
	switch c {
	case CategoryTech:
		return "科技"
	case CategorySports:
		return "体育"
	case CategoryPolitics:
		return "政治"
	case CategoryMilitary:
		return "军事"
	case CategoryEntertainment:
		return "娱乐"
	case CategoryBusiness:
		return "商业"
	case CategoryHealth:
		return "健康"
	case CategoryGeneral:
		return "综合"
	default:
		return "未知"
	}
}

// 获取所有分类
func GetAllCategories() []CategoryInfo {
	return []CategoryInfo{
		{Category: CategoryTech, Name: "科技", Description: "科技新闻和创新资讯"},
		{Category: CategorySports, Name: "体育", Description: "体育赛事和运动资讯"},
		{Category: CategoryPolitics, Name: "政治", Description: "政治新闻和时事评论"},
		{Category: CategoryMilitary, Name: "军事", Description: "军事新闻和国防资讯"},
		{Category: CategoryEntertainment, Name: "娱乐", Description: "娱乐新闻和明星资讯"},
		{Category: CategoryBusiness, Name: "商业", Description: "商业新闻和财经资讯"},
		{Category: CategoryHealth, Name: "健康", Description: "健康资讯和医疗新闻"},
		{Category: CategoryGeneral, Name: "综合", Description: "综合新闻和社会资讯"},
	}
}