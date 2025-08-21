package service

import (
	"errors"
	"fmt"
	"social-app/shared/config"
	"social-app/shared/models"
	"social-app/services/user-service/repository"
	"social-app/services/user-service/security"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req *RegisterRequest) (*AuthResponse, error)
	Login(req *LoginRequest) (*AuthResponse, error)
	RefreshToken(token string) (*AuthResponse, error)
	GetProfile(userID string) (*UserProfileResponse, error)
	UpdateProfile(userID string, req *UpdateProfileRequest) error
	UpdatePhotos(userID string, photoURLs []string) error
	DeletePhoto(userID string, photoURL string) error
	UpdateInterests(userID string, interests []string) error
	GetProfileCompleteness(userID string) (*ProfileCompletenessResponse, error)
	VerifyIdentity(userID string, req *IdentityVerificationRequest) error
	VerifyFace(userID string, req *FaceVerificationRequest) error
	GetInterestTags() ([]InterestTagResponse, error)
	SearchUsers(query string, currentUserID string) ([]UserSearchResult, error)
	ChangePassword(phone, newPassword string) error
	SaveRecentAccount(req *SaveRecentAccountRequest) error
	GetRecentAccounts() ([]RecentAccountResponse, error)
	
	// 好友相关方法
	SendFriendRequest(fromUserID, toUserID string, message *string) error
	AcceptFriendRequest(requestID string) error
	RejectFriendRequest(requestID string) error
	GetFriends(userID string) ([]FriendResponse, error)
	GetFriendRequests(userID string) ([]FriendRequestResponse, error)
	GetSentFriendRequests(userID string) ([]FriendRequestResponse, error)
	DeleteFriend(userID, friendID string) error
	CheckFriendStatus(userID, targetUserID string) (*FriendStatusResponse, error)
	GetFriendRecommendations(userID string) ([]FriendRecommendationResponse, error)
	GetRandomUsers(userID string, limit int) ([]UserSearchResult, error)
}

type userService struct {
	userRepo      repository.UserRepository
	interestRepo  repository.InterestRepository
	friendRepo    repository.FriendRepository
	recentAccountRepo repository.RecentAccountRepository
	config        *config.Config
	attemptService *security.LoginAttemptService
}

func NewUserService(userRepo repository.UserRepository, interestRepo repository.InterestRepository, friendRepo repository.FriendRepository, recentAccountRepo repository.RecentAccountRepository, cfg *config.Config) UserService {
	return &userService{
		userRepo:      userRepo,
		interestRepo:  interestRepo,
		friendRepo:    friendRepo,
		recentAccountRepo: recentAccountRepo,
		config:        cfg,
		attemptService: nil, // 暂时设为nil，后续可以通过依赖注入添加Redis客户端
	}
}

// 注册用户
func (s *userService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// 检查用户是否已存在
	if req.Phone != "" {
		if _, err := s.userRepo.GetByPhone(req.Phone); err == nil {
			return nil, errors.New("手机号已被注册")
		}
	}
	
	if req.Email != "" {
		if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
			return nil, errors.New("邮箱已被注册")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 创建用户
	user := &models.User{
		Phone: req.Phone,
		Email: &req.Email,
		Profile: models.UserProfile{
			Nickname: req.Nickname,
			Age:      req.Age,
			Gender:   req.Gender,
		},
		Preferences: models.UserPreferences{
			AgeRange: models.AgeRange{
				Min: 18,
				Max: 35,
			},
			GenderFilter:  "all",
			DistanceRange: 10,
			ShowOnline:    true,
			ShowDistance:  true,
		},
		Privacy: models.PrivacySettings{
			ProfileVisibility: "public",
			LocationVisible:   true,
			OnlineStatus:      true,
			ReadReceipts:      true,
		},
	}

	// 存储加密后的密码（这里简化处理，实际应该有专门的密码字段）
	hashedPasswordStr := string(hashedPassword)
	user.Profile.Bio = &hashedPasswordStr // 临时存储，实际应该有password字段

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	// 生成JWT token
	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	return &AuthResponse{
		Token:     token,
		ExpiresIn: s.config.JWT.ExpireTime * 3600, // 转换为秒
		User: UserResponse{
			ID:       user.ID,
			Phone:    user.Phone,
			Email:    user.Email,
			Nickname: user.Profile.Nickname,
			Avatar:   user.Profile.Avatar,
			Age:      user.Profile.Age,
			Gender:   user.Profile.Gender,
		},
	}, nil
}

// 用户登录
func (s *userService) Login(req *LoginRequest) (*AuthResponse, error) {
	// 确定用户标识符
	identifier := req.Phone
	if identifier == "" {
		identifier = req.Email
	}
	if identifier == "" {
		return nil, errors.New("请提供手机号或邮箱")
	}

	// 检查登录尝试次数限制
	if s.attemptService != nil {
		if err := s.attemptService.CheckUserAttempts(identifier); err != nil {
			return nil, err
		}
		
		// 这里应该从请求中获取真实IP，简化处理使用固定值
		clientIP := "127.0.0.1" // 实际应该从gin.Context中获取
		if err := s.attemptService.CheckIPAttempts(clientIP); err != nil {
			return nil, err
		}
	}

	var user *models.User
	var err error

	// 根据登录方式获取用户
	if req.Phone != "" {
		user, err = s.userRepo.GetByPhone(req.Phone)
	} else if req.Email != "" {
		user, err = s.userRepo.GetByEmail(req.Email)
	}

	// 统一错误信息，防止用户枚举
	if err != nil {
		if s.attemptService != nil {
			s.attemptService.RecordFailedAttempt(identifier, "127.0.0.1")
		}
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if user.Profile.Bio == nil {
		if s.attemptService != nil {
			s.attemptService.RecordFailedAttempt(identifier, "127.0.0.1")
		}
		return nil, errors.New("用户名或密码错误")
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(*user.Profile.Bio), []byte(req.Password))
	if err != nil {
		if s.attemptService != nil {
			s.attemptService.RecordFailedAttempt(identifier, "127.0.0.1")
		}
		return nil, errors.New("用户名或密码错误")
	}

	// 登录成功，清除失败尝试记录
	if s.attemptService != nil {
		s.attemptService.ClearUserAttempts(identifier)
	}

	// 生成JWT token
	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	return &AuthResponse{
		Token:     token,
		ExpiresIn: s.config.JWT.ExpireTime * 3600,
		User: UserResponse{
			ID:       user.ID,
			Phone:    user.Phone,
			Email:    user.Email,
			Nickname: user.Profile.Nickname,
			Avatar:   user.Profile.Avatar,
			Age:      user.Profile.Age,
			Gender:   user.Profile.Gender,
		},
	}, nil
}

// 刷新token
func (s *userService) RefreshToken(tokenString string) (*AuthResponse, error) {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("无效的token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("无效的token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("无效的用户ID")
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %v", err)
	}

	// 生成新token
	newToken, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	return &AuthResponse{
		Token:     newToken,
		ExpiresIn: s.config.JWT.ExpireTime * 3600,
		User: UserResponse{
			ID:       user.ID,
			Phone:    user.Phone,
			Email:    user.Email,
			Nickname: user.Profile.Nickname,
			Avatar:   user.Profile.Avatar,
			Age:      user.Profile.Age,
			Gender:   user.Profile.Gender,
		},
	}, nil
}

// 获取用户资料
func (s *userService) GetProfile(userID string) (*UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %v", err)
	}

	return &UserProfileResponse{
		ID:           user.ID,
		Phone:        user.Phone,
		Email:        user.Email,
		Nickname:     user.Profile.Nickname,
		Avatar:       user.Profile.Avatar,
		Age:          user.Profile.Age,
		Gender:       user.Profile.Gender,
		Bio:          user.Profile.Bio,
		Interests:    user.Profile.Interests,
		Photos:       user.Profile.Photos,
		Location:     user.Profile.Location,
		Verification: user.Verification,
		Preferences:  user.Preferences,
		Privacy:      user.Privacy,
	}, nil
}

// 更新用户资料
func (s *userService) UpdateProfile(userID string, req *UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	// 更新用户资料
	if req.Nickname != "" {
		user.Profile.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Profile.Avatar = req.Avatar
	}
	if req.Age > 0 {
		user.Profile.Age = req.Age
	}
	if req.Gender != "" {
		user.Profile.Gender = req.Gender
	}
	if req.Bio != "" {
		user.Profile.Bio = &req.Bio
	}
	if len(req.Interests) > 0 {
		user.Profile.Interests = req.Interests
		// 需要明确更新数组字段
		if err := s.userRepo.Update(user); err != nil {
			return fmt.Errorf("更新兴趣标签失败: %v", err)
		}
	}
	if len(req.Photos) > 0 {
		user.Profile.Photos = req.Photos
	}

	return s.userRepo.Update(user)
}

// 身份认证
func (s *userService) VerifyIdentity(userID string, req *IdentityVerificationRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	// 这里应该调用第三方身份认证API
	// 简化处理，直接标记为已认证
	now := time.Now()
	verification := &models.VerificationStatus{
		Identity: models.IdentityVerification{
			Verified:   true,
			IDNumber:   &req.IDNumber,
			RealName:   &req.RealName,
			VerifiedAt: &now,
		},
		Face: user.Verification.Face, // 保持原有人脸认证状态
	}

	return s.userRepo.UpdateVerification(userID, verification)
}

// 人脸认证
func (s *userService) VerifyFace(userID string, req *FaceVerificationRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	// 这里应该调用第三方人脸识别API
	// 简化处理，直接标记为已认证
	now := time.Now()
	verification := &models.VerificationStatus{
		Identity: user.Verification.Identity, // 保持原有身份认证状态
		Face: models.FaceVerification{
			Verified:   true,
			Confidence: 0.95, // 模拟95%置信度
			VerifiedAt: &now,
		},
	}

	return s.userRepo.UpdateVerification(userID, verification)
}


// 生成JWT token
func (s *userService) generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(s.config.JWT.ExpireTime) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

// 更新用户照片
func (s *userService) UpdatePhotos(userID string, photoURLs []string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	// 合并现有照片和新照片，限制最多6张
	allPhotos := append(user.Profile.Photos, photoURLs...)
	if len(allPhotos) > 6 {
		allPhotos = allPhotos[:6]
	}

	user.Profile.Photos = allPhotos
	return s.userRepo.Update(user)
}

// 删除用户照片
func (s *userService) DeletePhoto(userID string, photoURL string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	// 从照片列表中移除指定照片
	var newPhotos []string
	for _, photo := range user.Profile.Photos {
		if photo != photoURL {
			newPhotos = append(newPhotos, photo)
		}
	}

	user.Profile.Photos = newPhotos
	return s.userRepo.Update(user)
}

// 更新用户兴趣标签
func (s *userService) UpdateInterests(userID string, interests []string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	// 验证兴趣标签是否存在
	if len(interests) > 0 {
		validTags, err := s.interestRepo.GetByIDs(interests)
		if err != nil {
			return fmt.Errorf("验证兴趣标签失败: %v", err)
		}

		if len(validTags) != len(interests) {
			return errors.New("包含无效的兴趣标签")
		}

		// 更新标签使用次数
		for _, tagID := range interests {
			s.interestRepo.UpdateUsageCount(tagID)
		}
	}

	user.Profile.Interests = interests
	return s.userRepo.Update(user)
}

// 获取资料完整度
func (s *userService) GetProfileCompleteness(userID string) (*ProfileCompletenessResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %v", err)
	}

	completeness := &ProfileCompletenessResponse{
		UserID:           userID,
		OverallScore:     0,
		CompletedFields:  []string{},
		MissingFields:    []string{},
		Recommendations:  []string{},
	}

	totalFields := 10
	completedCount := 0

	// 检查基本信息
	if user.Profile.Nickname != "" {
		completeness.CompletedFields = append(completeness.CompletedFields, "nickname")
		completedCount++
	} else {
		completeness.MissingFields = append(completeness.MissingFields, "nickname")
		completeness.Recommendations = append(completeness.Recommendations, "添加昵称让其他人更容易记住你")
	}

	if user.Profile.Avatar != "" {
		completeness.CompletedFields = append(completeness.CompletedFields, "avatar")
		completedCount++
	} else {
		completeness.MissingFields = append(completeness.MissingFields, "avatar")
		completeness.Recommendations = append(completeness.Recommendations, "上传头像提高匹配成功率")
	}

	if user.Profile.Age > 0 {
		completeness.CompletedFields = append(completeness.CompletedFields, "age")
		completedCount++
	} else {
		completeness.MissingFields = append(completeness.MissingFields, "age")
	}

	if user.Profile.Gender != "" {
		completeness.CompletedFields = append(completeness.CompletedFields, "gender")
		completedCount++
	} else {
		completeness.MissingFields = append(completeness.MissingFields, "gender")
	}

	if user.Profile.Bio != nil && *user.Profile.Bio != "" {
		completeness.CompletedFields = append(completeness.CompletedFields, "bio")
		completedCount++
	} else {
		completeness.MissingFields = append(completeness.MissingFields, "bio")
		completeness.Recommendations = append(completeness.Recommendations, "添加个人简介展示你的个性")
	}

	if len(user.Profile.Interests) > 0 {
		completeness.CompletedFields = append(completeness.CompletedFields, "interests")
		completedCount++
	} else {
		completeness.MissingFields = append(completeness.MissingFields, "interests")
		completeness.Recommendations = append(completeness.Recommendations, "选择兴趣标签找到志同道合的人")
	}

	if len(user.Profile.Photos) > 0 {
		completeness.CompletedFields = append(completeness.CompletedFields, "photos")
		completedCount++
	} else {
		completeness.MissingFields = append(completeness.MissingFields, "photos")
		completeness.Recommendations = append(completeness.Recommendations, "上传更多照片展示真实的自己")
	}

	if user.Profile.Location != nil {
		completeness.CompletedFields = append(completeness.CompletedFields, "location")
		completedCount++
	} else {
		completeness.MissingFields = append(completeness.MissingFields, "location")
		completeness.Recommendations = append(completeness.Recommendations, "设置位置信息找到附近的人")
	}

	if user.Verification.Identity.Verified {
		completeness.CompletedFields = append(completeness.CompletedFields, "identity_verification")
		completedCount++
	} else {
		completeness.MissingFields = append(completeness.MissingFields, "identity_verification")
		completeness.Recommendations = append(completeness.Recommendations, "完成身份认证提高账号可信度")
	}

	if user.Verification.Face.Verified {
		completeness.CompletedFields = append(completeness.CompletedFields, "face_verification")
		completedCount++
	} else {
		completeness.MissingFields = append(completeness.MissingFields, "face_verification")
		completeness.Recommendations = append(completeness.Recommendations, "完成人脸认证确保账号安全")
	}

	// 计算完整度百分比
	completeness.OverallScore = int((float64(completedCount) / float64(totalFields)) * 100)

	return completeness, nil
}

// 获取兴趣标签
func (s *userService) GetInterestTags() ([]InterestTagResponse, error) {
	tags, err := s.interestRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("获取兴趣标签失败: %v", err)
	}

	var response []InterestTagResponse
	for _, tag := range tags {
		response = append(response, InterestTagResponse{
			ID:          tag.ID,
			Name:        tag.Name,
			Category:    tag.Category,
			Description: tag.Description,
			Icon:        tag.Icon,
			Color:       tag.Color,
			UsageCount:  tag.UsageCount,
		})
	}

	return response, nil
}

// 搜索用户
func (s *userService) SearchUsers(query string, currentUserID string) ([]UserSearchResult, error) {
	if query == "" {
		return []UserSearchResult{}, nil
	}

	// 从仓库层搜索用户
	users, err := s.userRepo.SearchUsers(query, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("搜索用户失败: %v", err)
	}

	// 转换为响应格式
	var results []UserSearchResult
	for _, user := range users {
		result := UserSearchResult{
			ID:       user.ID,
			Phone:    user.Phone,
			Nickname: user.Profile.Nickname,
			Avatar:   user.Profile.Avatar,
			Age:      user.Profile.Age,
			Gender:   user.Profile.Gender,
		}
		
		// 添加简介（如果有的话）
		if user.Profile.Bio != nil {
			result.Bio = *user.Profile.Bio
		}
		
		results = append(results, result)
	}

	return results, nil
}

// 发送好友请求
func (s *userService) SendFriendRequest(fromUserID, toUserID string, message *string) error {
	// 检查是否已经是好友
	existingFriend, _ := s.friendRepo.GetFriendByUserIDs(fromUserID, toUserID)
	if existingFriend != nil {
		return fmt.Errorf("已经是好友关系")
	}

	// 检查是否已经发送过请求
	existingRequest, _ := s.friendRepo.GetFriendRequestByUserIDs(fromUserID, toUserID)
	if existingRequest != nil {
		return fmt.Errorf("已经发送过好友请求")
	}

	// 创建好友请求
	request := &models.FriendRequest{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Message:    message,
		Status:     models.FriendRequestStatusPending,
	}

	return s.friendRepo.CreateFriendRequest(request)
}

// 接受好友请求
func (s *userService) AcceptFriendRequest(requestID string) error {
	// 获取好友请求
	request, err := s.friendRepo.GetFriendRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("好友请求不存在")
	}

	if request.Status != models.FriendRequestStatusPending {
		return fmt.Errorf("好友请求已处理")
	}

	// 创建双向好友关系
	friend1 := &models.Friend{
		UserID:   request.FromUserID,
		FriendID: request.ToUserID,
		Status:   models.FriendStatusAccepted,
	}

	friend2 := &models.Friend{
		UserID:   request.ToUserID,
		FriendID: request.FromUserID,
		Status:   models.FriendStatusAccepted,
	}

	// 创建好友关系
	if err := s.friendRepo.CreateFriend(friend1); err != nil {
		return fmt.Errorf("创建好友关系失败: %v", err)
	}

	if err := s.friendRepo.CreateFriend(friend2); err != nil {
		return fmt.Errorf("创建好友关系失败: %v", err)
	}

	// 更新请求状态
	if err := s.friendRepo.UpdateFriendRequestStatus(requestID, models.FriendRequestStatusAccepted); err != nil {
		return fmt.Errorf("更新请求状态失败: %v", err)
	}

	return nil
}

// 拒绝好友请求
func (s *userService) RejectFriendRequest(requestID string) error {
	// 获取好友请求
	request, err := s.friendRepo.GetFriendRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("好友请求不存在")
	}

	if request.Status != models.FriendRequestStatusPending {
		return fmt.Errorf("好友请求已处理")
	}

	// 更新请求状态
	return s.friendRepo.UpdateFriendRequestStatus(requestID, models.FriendRequestStatusRejected)
}

// 获取好友列表
func (s *userService) GetFriends(userID string) ([]FriendResponse, error) {
	friends, err := s.friendRepo.GetFriendsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取好友列表失败: %v", err)
	}

	var response []FriendResponse
	for _, friend := range friends {
		friendResponse := FriendResponse{
			ID:       friend.ID,
			UserID:   friend.UserID,
			FriendID: friend.FriendID,
			Nickname: friend.Friend.Profile.Nickname,
			Avatar:   friend.Friend.Profile.Avatar,
			Phone:    friend.Friend.Phone,
			Gender:   friend.Friend.Profile.Gender,
			Status:   string(friend.Status),
			CreatedAt: friend.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if friend.Friend.Email != nil {
			friendResponse.Email = *friend.Friend.Email
		}

		response = append(response, friendResponse)
	}

	return response, nil
}

// 获取好友请求列表
func (s *userService) GetFriendRequests(userID string) ([]FriendRequestResponse, error) {
	requests, err := s.friendRepo.GetFriendRequestsByToUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取好友请求失败: %v", err)
	}

	var response []FriendRequestResponse
	for _, request := range requests {
		requestResponse := FriendRequestResponse{
			ID:               request.ID,
			FromUserID:       request.FromUserID,
			ToUserID:         request.ToUserID,
			FromUserNickname: request.FromUser.Profile.Nickname,
			FromUserAvatar:   request.FromUser.Profile.Avatar,
			FromUserPhone:    request.FromUser.Phone,
			Status:           string(request.Status),
			CreatedAt:        request.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if request.Message != nil {
			requestResponse.Message = *request.Message
		}

		response = append(response, requestResponse)
	}

	return response, nil
}

// 获取发送的好友请求列表
func (s *userService) GetSentFriendRequests(userID string) ([]FriendRequestResponse, error) {
	requests, err := s.friendRepo.GetFriendRequestsByFromUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取发送的好友请求失败: %v", err)
	}

	var response []FriendRequestResponse
	for _, request := range requests {
		requestResponse := FriendRequestResponse{
			ID:               request.ID,
			FromUserID:       request.FromUserID,
			ToUserID:         request.ToUserID,
			FromUserNickname: request.FromUser.Profile.Nickname,
			FromUserAvatar:   request.FromUser.Profile.Avatar,
			FromUserPhone:    request.FromUser.Phone,
			Status:           string(request.Status),
			CreatedAt:        request.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if request.Message != nil {
			requestResponse.Message = *request.Message
		}

		response = append(response, requestResponse)
	}

	return response, nil
}

// 删除好友
func (s *userService) DeleteFriend(userID, friendID string) error {
	// 获取好友关系
	friend, err := s.friendRepo.GetFriendByUserIDs(userID, friendID)
	if err != nil {
		return fmt.Errorf("好友关系不存在")
	}

	// 删除双向好友关系
	// 这里需要删除两条记录
	return s.friendRepo.DeleteFriend(friend.ID)
}

// 检查好友状态
func (s *userService) CheckFriendStatus(userID, targetUserID string) (*FriendStatusResponse, error) {
	// 检查是否已经是好友
	friend, err := s.friendRepo.GetFriendByUserIDs(userID, targetUserID)
	if err == nil && friend != nil {
		return &FriendStatusResponse{
			Status:    "friend",
			IsFriend:  true,
			CanSendRequest: false,
			Message:   "已经是好友",
		}, nil
	}

	// 检查是否有待处理的好友请求（我发送的）
	sentRequest, err := s.friendRepo.GetFriendRequestByUserIDs(userID, targetUserID)
	if err == nil && sentRequest != nil {
		if sentRequest.Status == models.FriendRequestStatusPending {
			return &FriendStatusResponse{
				Status:    "request_sent",
				IsFriend:  false,
				CanSendRequest: false,
				Message:   "已发送好友请求，等待对方确认",
			}, nil
		}
	}

	// 检查是否有待处理的好友请求（对方发送的）
	receivedRequest, err := s.friendRepo.GetFriendRequestByUserIDs(targetUserID, userID)
	if err == nil && receivedRequest != nil {
		if receivedRequest.Status == models.FriendRequestStatusPending {
			return &FriendStatusResponse{
				Status:    "request_received",
				IsFriend:  false,
				CanSendRequest: false,
				Message:   "对方已向你发送好友请求",
			}, nil
		}
	}

	// 可以发送好友请求
	return &FriendStatusResponse{
		Status:    "none",
		IsFriend:  false,
		CanSendRequest: true,
		Message:   "可以发送好友请求",
	}, nil
}

// 更改密码
func (s *userService) ChangePassword(phone, newPassword string) error {
	// 获取用户
	user, err := s.userRepo.GetByPhone(phone)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	// 更新密码（使用Bio字段作为临时存储）
	hashedPasswordStr := string(hashedPassword)
	user.Profile.Bio = &hashedPasswordStr

	return s.userRepo.Update(user)
}

// 保存最近登录账号
func (s *userService) SaveRecentAccount(req *SaveRecentAccountRequest) error {
	account := &models.RecentAccount{
		Phone:            req.Phone,
		Password:         req.Password,
		RememberPassword: req.RememberPassword,
	}
	return s.recentAccountRepo.Save(account)
}

// 获取最近登录账号
func (s *userService) GetRecentAccounts() ([]RecentAccountResponse, error) {
	accounts, err := s.recentAccountRepo.GetAll()
	if err != nil {
		return nil, err
	}
	
	var response []RecentAccountResponse
	for _, account := range accounts {
		response = append(response, RecentAccountResponse{
			Phone:            account.Phone,
			Password:         account.Password,
			RememberPassword: account.RememberPassword,
			LastLoginAt:      account.LastLoginAt.Format("2006-01-02 15:04:05"),
		})
	}
	
	return response, nil
}

// 获取推荐好友
func (s *userService) GetFriendRecommendations(userID string) ([]FriendRecommendationResponse, error) {
	// 获取当前用户信息
	currentUser, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %v", err)
	}

	// 获取已有好友列表
	existingFriends, _ := s.friendRepo.GetFriendsByUserID(userID)
	existingFriendIDs := make(map[string]bool)
	for _, friend := range existingFriends {
		existingFriendIDs[friend.FriendID] = true
	}

	// 获取已发送的好友请求
	sentRequests, _ := s.friendRepo.GetFriendRequestsByFromUserID(userID)
	sentRequestIDs := make(map[string]bool)
	for _, request := range sentRequests {
		if request.Status == models.FriendRequestStatusPending {
			sentRequestIDs[request.ToUserID] = true
		}
	}

	// 获取所有用户（除了自己和已有好友）
	allUsers, err := s.userRepo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %v", err)
	}

	var recommendations []FriendRecommendationResponse
	for _, user := range allUsers {
		// 跳过自己、已有好友和已发送请求的用户
		if user.ID == userID || existingFriendIDs[user.ID] || sentRequestIDs[user.ID] {
			continue
		}

		// 计算共同兴趣数量
		commonInterests := s.calculateCommonInterests(currentUser.Profile.Interests, user.Profile.Interests)

		// 只推荐有共同兴趣的用户
		if commonInterests > 0 {
			recommendation := FriendRecommendationResponse{
				UserID:          user.ID,
				Nickname:        user.Profile.Nickname,
				Avatar:          user.Profile.Avatar,
				Gender:          user.Profile.Gender,
				Age:             user.Profile.Age,
				CommonInterests: commonInterests,
				Distance:        0.0, // 简化处理，后续可以根据位置计算
				Interests:       user.Profile.Interests,
			}

			if user.Profile.Bio != nil {
				recommendation.Bio = *user.Profile.Bio
			}

			recommendations = append(recommendations, recommendation)
		}
	}

	// 按共同兴趣数量排序，返回前10个
	for i := 0; i < len(recommendations)-1; i++ {
		for j := i + 1; j < len(recommendations); j++ {
			if recommendations[i].CommonInterests < recommendations[j].CommonInterests {
				recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
			}
		}
	}

	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}

	return recommendations, nil
}

// 计算共同兴趣数量
func (s *userService) calculateCommonInterests(interests1, interests2 []string) int {
	interestMap := make(map[string]bool)
	for _, interest := range interests1 {
		interestMap[interest] = true
	}

	commonCount := 0
	for _, interest := range interests2 {
		if interestMap[interest] {
			commonCount++
		}
	}

	return commonCount
}

// 获取随机用户
func (s *userService) GetRandomUsers(userID string, limit int) ([]UserSearchResult, error) {
	// 获取当前用户信息
	currentUser, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %v", err)
	}

	// 获取已有好友列表
	existingFriends, _ := s.friendRepo.GetFriendsByUserID(userID)
	existingFriendIDs := make(map[string]bool)
	for _, friend := range existingFriends {
		existingFriendIDs[friend.FriendID] = true
	}

	// 获取随机异性用户
	users, err := s.userRepo.GetRandomUsersByGender(userID, s.getOppositeGender(currentUser.Profile.Gender), limit*2) // 获取更多用户以便筛选
	if err != nil {
		return nil, fmt.Errorf("获取随机用户失败: %v", err)
	}

	var results []UserSearchResult
	for _, user := range users {
		// 跳过已有好友
		if existingFriendIDs[user.ID] {
			continue
		}

		result := UserSearchResult{
			ID:       user.ID,
			Phone:    user.Phone,
			Nickname: user.Profile.Nickname,
			Avatar:   user.Profile.Avatar,
			Age:      user.Profile.Age,
			Gender:   user.Profile.Gender,
		}

		if user.Profile.Bio != nil {
			result.Bio = *user.Profile.Bio
		}

		results = append(results, result)

		// 达到限制数量就停止
		if len(results) >= limit {
			break
		}
	}

	return results, nil
}

// 获取相反性别
func (s *userService) getOppositeGender(gender string) string {
	switch gender {
	case "male":
		return "female"
	case "female":
		return "male"
	default:
		return "" // 对于other性别，返回空字符串表示不限制
	}
}