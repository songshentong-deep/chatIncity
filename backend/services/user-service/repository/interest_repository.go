package repository

import (
	"social-app/shared/models"
	"gorm.io/gorm"
)

type InterestRepository interface {
	GetAll() ([]models.InterestTag, error)
	GetByCategory(category string) ([]models.InterestTag, error)
	GetByIDs(ids []string) ([]models.InterestTag, error)
	UpdateUsageCount(id string) error
}

type interestRepository struct {
	db *gorm.DB
}

func NewInterestRepository(db *gorm.DB) InterestRepository {
	return &interestRepository{db: db}
}

func (r *interestRepository) GetAll() ([]models.InterestTag, error) {
	var tags []models.InterestTag
	err := r.db.Where("is_active = ?", true).Order("category, usage_count DESC").Find(&tags).Error
	return tags, err
}

func (r *interestRepository) GetByCategory(category string) ([]models.InterestTag, error) {
	var tags []models.InterestTag
	err := r.db.Where("category = ? AND is_active = ?", category, true).Order("usage_count DESC").Find(&tags).Error
	return tags, err
}

func (r *interestRepository) GetByIDs(ids []string) ([]models.InterestTag, error) {
	var tags []models.InterestTag
	err := r.db.Where("id IN ? AND is_active = ?", ids, true).Find(&tags).Error
	return tags, err
}

func (r *interestRepository) UpdateUsageCount(id string) error {
	return r.db.Model(&models.InterestTag{}).Where("id = ?", id).UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}