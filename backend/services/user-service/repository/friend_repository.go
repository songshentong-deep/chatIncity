package repository

import (
	"social-app/shared/models"
	"gorm.io/gorm"
)

type FriendRepository interface {
	// 好友关系管理
	CreateFriend(friend *models.Friend) error
	GetFriendsByUserID(userID string) ([]*models.Friend, error)
	GetFriendByUserIDs(userID, friendID string) (*models.Friend, error)
	UpdateFriendStatus(friendID string, status models.FriendStatus) error
	DeleteFriend(friendID string) error
	
	// 好友请求管理
	CreateFriendRequest(request *models.FriendRequest) error
	GetFriendRequestsByToUserID(toUserID string) ([]*models.FriendRequest, error)
	GetFriendRequestsByFromUserID(fromUserID string) ([]*models.FriendRequest, error)
	GetFriendRequestByID(requestID string) (*models.FriendRequest, error)
	GetFriendRequestByUserIDs(fromUserID, toUserID string) (*models.FriendRequest, error)
	UpdateFriendRequestStatus(requestID string, status models.FriendRequestStatus) error
	DeleteFriendRequest(requestID string) error
}

type friendRepository struct {
	db *gorm.DB
}

func NewFriendRepository(db *gorm.DB) FriendRepository {
	return &friendRepository{db: db}
}

// 好友关系管理
func (r *friendRepository) CreateFriend(friend *models.Friend) error {
	return r.db.Create(friend).Error
}

func (r *friendRepository) GetFriendsByUserID(userID string) ([]*models.Friend, error) {
	var friends []*models.Friend
	err := r.db.Preload("Friend").Where("user_id = ? AND status = ?", userID, models.FriendStatusAccepted).Find(&friends).Error
	return friends, err
}

func (r *friendRepository) GetFriendByUserIDs(userID, friendID string) (*models.Friend, error) {
	var friend models.Friend
	err := r.db.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", 
		userID, friendID, friendID, userID).First(&friend).Error
	if err != nil {
		return nil, err
	}
	return &friend, nil
}

func (r *friendRepository) UpdateFriendStatus(friendID string, status models.FriendStatus) error {
	return r.db.Model(&models.Friend{}).Where("id = ?", friendID).Update("status", status).Error
}

func (r *friendRepository) DeleteFriend(friendID string) error {
	return r.db.Delete(&models.Friend{}, "id = ?", friendID).Error
}

// 好友请求管理
func (r *friendRepository) CreateFriendRequest(request *models.FriendRequest) error {
	return r.db.Create(request).Error
}

func (r *friendRepository) GetFriendRequestsByToUserID(toUserID string) ([]*models.FriendRequest, error) {
	var requests []*models.FriendRequest
	err := r.db.Preload("FromUser").Where("to_user_id = ? AND status = ?", toUserID, models.FriendRequestStatusPending).Find(&requests).Error
	return requests, err
}

func (r *friendRepository) GetFriendRequestsByFromUserID(fromUserID string) ([]*models.FriendRequest, error) {
	var requests []*models.FriendRequest
	err := r.db.Preload("ToUser").Where("from_user_id = ?", fromUserID).Find(&requests).Error
	return requests, err
}

func (r *friendRepository) GetFriendRequestByID(requestID string) (*models.FriendRequest, error) {
	var request models.FriendRequest
	err := r.db.Preload("FromUser").Preload("ToUser").Where("id = ?", requestID).First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *friendRepository) GetFriendRequestByUserIDs(fromUserID, toUserID string) (*models.FriendRequest, error) {
	var request models.FriendRequest
	err := r.db.Where("from_user_id = ? AND to_user_id = ? AND status = ?", 
		fromUserID, toUserID, models.FriendRequestStatusPending).First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *friendRepository) UpdateFriendRequestStatus(requestID string, status models.FriendRequestStatus) error {
	return r.db.Model(&models.FriendRequest{}).Where("id = ?", requestID).Update("status", status).Error
}

func (r *friendRepository) DeleteFriendRequest(requestID string) error {
	return r.db.Delete(&models.FriendRequest{}, "id = ?", requestID).Error
}