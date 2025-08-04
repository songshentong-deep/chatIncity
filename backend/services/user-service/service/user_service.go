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
}

type userService struct {
	userRepo      repository.UserRepository
	interestRepo  repository.InterestRepository
	config        *config.Config
	attemptService *security.LoginAttemptService
}

func NewUserService(userRepo repository.UserRepository, interestRepo repository.InterestRepository, cfg *config.Config) UserService {
	return &userService{
		userRepo:      userRepo,
		interestRepo:  interestRepo,
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