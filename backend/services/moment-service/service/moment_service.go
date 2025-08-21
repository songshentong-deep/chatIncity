package service

import (
	"fmt"
	"social-app/shared/models"
	"social-app/services/moment-service/repository"
)

type MomentService interface {
	CreateMoment(req *CreateMomentRequest) error
	GetMyMoments(userID string) ([]MomentResponse, error)
	GetFriendsMoments(userID string) ([]MomentResponse, error)
	DeleteMoment(momentID, userID string) error
}

type momentService struct {
	momentRepo      repository.MomentRepository
	interactionRepo repository.InteractionRepository
}

func NewMomentService(momentRepo repository.MomentRepository, interactionRepo repository.InteractionRepository) MomentService {
	return &momentService{
		momentRepo:      momentRepo,
		interactionRepo: interactionRepo,
	}
}

func (s *momentService) CreateMoment(req *CreateMomentRequest) error {
	// 处理图片URL列表
	images := models.StringArray{}
	if req.Images != nil && len(*req.Images) > 0 {
		images = models.StringArray(*req.Images)
	}
	
	moment := &models.Moment{
		UserID:   req.UserID,
		Content:  req.Content,
		Images:   images,
		VideoURL: req.VideoURL,
		Type:     req.Type,
	}

	return s.momentRepo.Create(moment)
}

func (s *momentService) GetMyMoments(userID string) ([]MomentResponse, error) {
	moments, err := s.momentRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取我的动态失败: %v", err)
	}

	var response []MomentResponse
	for _, moment := range moments {
		response = append(response, MomentResponse{
			ID:           moment.ID,
			UserID:       moment.UserID,
			UserNickname: moment.User.Profile.Nickname,
			UserAvatar:   moment.User.Profile.Avatar,
			Content:      moment.Content,
			Images:       []string(moment.Images),
			VideoURL:     moment.VideoURL,
			Type:         moment.Type,
			LikesCount:   s.getLikesCount(moment.ID),
			CommentsCount: s.getCommentsCount(moment.ID),
			IsLiked:      s.isLiked(moment.ID, userID),
			CreatedAt:    moment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

func (s *momentService) GetFriendsMoments(userID string) ([]MomentResponse, error) {
	moments, err := s.momentRepo.GetFriendsMoments(userID)
	if err != nil {
		return nil, fmt.Errorf("获取好友动态失败: %v", err)
	}

	var response []MomentResponse
	for _, moment := range moments {
		response = append(response, MomentResponse{
			ID:           moment.ID,
			UserID:       moment.UserID,
			UserNickname: moment.User.Profile.Nickname,
			UserAvatar:   moment.User.Profile.Avatar,
			Content:      moment.Content,
			Images:       []string(moment.Images),
			VideoURL:     moment.VideoURL,
			Type:         moment.Type,
			LikesCount:   s.getLikesCount(moment.ID),
			CommentsCount: s.getCommentsCount(moment.ID),
			IsLiked:      s.isLiked(moment.ID, userID),
			CreatedAt:    moment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

func (s *momentService) DeleteMoment(momentID, userID string) error {
	// 验证动态是否属于当前用户
	moment, err := s.momentRepo.GetByID(momentID)
	if err != nil {
		return fmt.Errorf("动态不存在")
	}

	if moment.UserID != userID {
		return fmt.Errorf("无权限删除此动态")
	}

	return s.momentRepo.Delete(momentID)
}

func (s *momentService) getLikesCount(momentID string) int {
	count, _ := s.interactionRepo.GetLikesCount(momentID)
	return int(count)
}

func (s *momentService) getCommentsCount(momentID string) int {
	count, _ := s.interactionRepo.GetCommentsCount(momentID)
	return int(count)
}

func (s *momentService) isLiked(momentID, userID string) bool {
	liked, _ := s.interactionRepo.IsLiked(momentID, userID)
	return liked
}