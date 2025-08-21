package repository

import (
	"social-app/shared/models"
	"gorm.io/gorm"
)

type MomentRepository interface {
	Create(moment *models.Moment) error
	GetByID(id string) (*models.Moment, error)
	GetByUserID(userID string) ([]*models.Moment, error)
	GetFriendsMoments(userID string) ([]*models.Moment, error)
	Delete(id string) error
}

type momentRepository struct {
	db *gorm.DB
}

func NewMomentRepository(db *gorm.DB) MomentRepository {
	return &momentRepository{db: db}
}

func (r *momentRepository) Create(moment *models.Moment) error {
	return r.db.Create(moment).Error
}

func (r *momentRepository) GetByID(id string) (*models.Moment, error) {
	var moment models.Moment
	err := r.db.Preload("User").Where("id = ?", id).First(&moment).Error
	if err != nil {
		return nil, err
	}
	return &moment, nil
}

func (r *momentRepository) GetByUserID(userID string) ([]*models.Moment, error) {
	var moments []*models.Moment
	err := r.db.Preload("User").Where("user_id = ?", userID).Order("created_at DESC").Find(&moments).Error
	return moments, err
}

func (r *momentRepository) GetFriendsMoments(userID string) ([]*models.Moment, error) {
	var moments []*models.Moment
	
	// 获取好友的动态，包括自己的动态
	err := r.db.Preload("User").
		Where("user_id IN (SELECT friend_id FROM friends WHERE user_id = ? AND status = 'accepted') OR user_id = ?", userID, userID).
		Order("created_at DESC").
		Limit(50).
		Find(&moments).Error
	
	return moments, err
}

func (r *momentRepository) Delete(id string) error {
	return r.db.Delete(&models.Moment{}, "id = ?", id).Error
}