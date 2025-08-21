package repository

import (
	"social-app/shared/models"
	"time"
	"gorm.io/gorm"
)

type RecentAccountRepository interface {
	Save(account *models.RecentAccount) error
	GetAll() ([]models.RecentAccount, error)
}

type recentAccountRepository struct {
	db *gorm.DB
}

func NewRecentAccountRepository(db *gorm.DB) RecentAccountRepository {
	return &recentAccountRepository{db: db}
}

func (r *recentAccountRepository) Save(account *models.RecentAccount) error {
	// 先查找是否已存在该手机号的记录
	var existing models.RecentAccount
	err := r.db.Where("phone = ?", account.Phone).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		// 不存在则创建新记录
		account.LastLoginAt = time.Now()
		return r.db.Create(account).Error
	} else if err != nil {
		return err
	}
	
	// 存在则更新记录
	existing.Password = account.Password
	existing.RememberPassword = account.RememberPassword
	existing.LastLoginAt = time.Now()
	return r.db.Save(&existing).Error
}

func (r *recentAccountRepository) GetAll() ([]models.RecentAccount, error) {
	var accounts []models.RecentAccount
	err := r.db.Order("last_login_at DESC").Limit(5).Find(&accounts).Error
	return accounts, err
}