package repository

import (
	"social-app/shared/models"
	"gorm.io/gorm"
)

type InteractionRepository interface {
	// 点赞相关
	LikeMoment(momentID, userID string) error
	UnlikeMoment(momentID, userID string) error
	IsLiked(momentID, userID string) (bool, error)
	GetMomentLikes(momentID string) ([]*models.MomentLike, error)
	GetLikesCount(momentID string) (int64, error)
	
	// 评论相关
	CreateComment(comment *models.MomentComment) error
	GetMomentComments(momentID string) ([]*models.MomentComment, error)
	GetCommentsCount(momentID string) (int64, error)
	DeleteComment(commentID, userID string) error
}

type interactionRepository struct {
	db *gorm.DB
}

func NewInteractionRepository(db *gorm.DB) InteractionRepository {
	return &interactionRepository{db: db}
}

// 点赞动态
func (r *interactionRepository) LikeMoment(momentID, userID string) error {
	like := &models.MomentLike{
		MomentID: momentID,
		UserID:   userID,
	}
	return r.db.Create(like).Error
}

// 取消点赞
func (r *interactionRepository) UnlikeMoment(momentID, userID string) error {
	return r.db.Where("moment_id = ? AND user_id = ?", momentID, userID).Delete(&models.MomentLike{}).Error
}

// 检查是否已点赞
func (r *interactionRepository) IsLiked(momentID, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.MomentLike{}).Where("moment_id = ? AND user_id = ?", momentID, userID).Count(&count).Error
	return count > 0, err
}

// 获取动态点赞列表
func (r *interactionRepository) GetMomentLikes(momentID string) ([]*models.MomentLike, error) {
	var likes []*models.MomentLike
	err := r.db.Preload("User").Where("moment_id = ?", momentID).Order("created_at DESC").Find(&likes).Error
	return likes, err
}

// 获取点赞数量
func (r *interactionRepository) GetLikesCount(momentID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.MomentLike{}).Where("moment_id = ?", momentID).Count(&count).Error
	return count, err
}

// 创建评论
func (r *interactionRepository) CreateComment(comment *models.MomentComment) error {
	return r.db.Create(comment).Error
}

// 获取动态评论列表
func (r *interactionRepository) GetMomentComments(momentID string) ([]*models.MomentComment, error) {
	var comments []*models.MomentComment
	err := r.db.Preload("User").Where("moment_id = ?", momentID).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

// 获取评论数量
func (r *interactionRepository) GetCommentsCount(momentID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.MomentComment{}).Where("moment_id = ?", momentID).Count(&count).Error
	return count, err
}

// 删除评论
func (r *interactionRepository) DeleteComment(commentID, userID string) error {
	return r.db.Where("id = ? AND user_id = ?", commentID, userID).Delete(&models.MomentComment{}).Error
}