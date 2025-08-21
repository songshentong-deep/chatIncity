package repository

import (
	"social-app/shared/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByPhone(phone string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
	UpdateVerification(id string, verification *models.VerificationStatus) error
	SearchUsers(query string, excludeUserID string) ([]*models.User, error)
	GetAllUsers() ([]*models.User, error)
	GetRandomUsersByGender(excludeUserID string, gender string, limit int) ([]*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *userRepository) UpdateVerification(id string, verification *models.VerificationStatus) error {
	updates := map[string]interface{}{}
	
	// 更新身份认证字段
	if verification.Identity.Verified {
		updates["verified"] = verification.Identity.Verified
		updates["id_number"] = verification.Identity.IDNumber
		updates["real_name"] = verification.Identity.RealName
		updates["verified_at"] = verification.Identity.VerifiedAt
	}
	
	// 更新人脸认证字段 - 需要添加到数据库表中
	// 注意：这里需要确保数据库表有对应的字段
	updates["confidence"] = verification.Face.Confidence
	
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *userRepository) SearchUsers(query string, excludeUserID string) ([]*models.User, error) {
	var users []*models.User
	
	// 构建搜索查询
	db := r.db.Where("id != ?", excludeUserID)
	
	// 支持按手机号、昵称搜索
	db = db.Where("phone LIKE ? OR nickname LIKE ?", "%"+query+"%", "%"+query+"%")
	
	// 限制返回结果数量
	db = db.Limit(20)
	
	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	
	return users, nil
}

func (r *userRepository) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetRandomUsersByGender(excludeUserID string, gender string, limit int) ([]*models.User, error) {
	var users []*models.User
	
	// 构建查询
	db := r.db.Where("id != ?", excludeUserID)
	
	// 如果指定了性别，添加性别过滤
	if gender != "" {
		db = db.Where("gender = ?", gender)
	}
	
	// 只返回有头像的用户
	db = db.Where("avatar != '' AND avatar IS NOT NULL")
	
	// 随机排序并限制数量
	db = db.Order("RANDOM()").Limit(limit)
	
	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	
	return users, nil
}