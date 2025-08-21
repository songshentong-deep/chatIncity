package service

import (
	"social-app/shared/models"
)

// 注册请求
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"required"`
	Age      int    `json:"age" binding:"required,min=18,max=100"`
	Gender   string `json:"gender" binding:"required,oneof=male female other"`
}

// 登录请求
type LoginRequest struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

// 更新资料请求
type UpdateProfileRequest struct {
	Nickname  string   `json:"nickname"`
	Avatar    string   `json:"avatar"`
	Age       int      `json:"age"`
	Gender    string   `json:"gender"`
	Bio       string   `json:"bio"`
	Interests []string `json:"interests"`
	Photos    []string `json:"photos"`
}

// 身份认证请求
type IdentityVerificationRequest struct {
	IDNumber string `json:"id_number" binding:"required"`
	RealName string `json:"real_name" binding:"required"`
	IDPhoto  string `json:"id_photo" binding:"required"` // base64编码的身份证照片
}

// 人脸认证请求
type FaceVerificationRequest struct {
	FaceImage string `json:"face_image" binding:"required"` // base64编码的人脸照片
}

// 更改密码请求
type ChangePasswordRequest struct {
	Phone       string `json:"phone" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// 保存最近登录账号请求
type SaveRecentAccountRequest struct {
	Phone          string `json:"phone" binding:"required"`
	Password       string `json:"password,omitempty"`
	RememberPassword bool `json:"remember_password"`
}

// 最近登录账号响应
type RecentAccountResponse struct {
	Phone          string `json:"phone"`
	Password       string `json:"password,omitempty"`
	RememberPassword bool `json:"remember_password"`
	LastLoginAt    string `json:"last_login_at"`
}

// 用户响应
type UserResponse struct {
	ID       string  `json:"id"`
	Phone    string  `json:"phone"`
	Email    *string `json:"email,omitempty"`
	Nickname string  `json:"nickname"`
	Avatar   string  `json:"avatar"`
	Age      int     `json:"age"`
	Gender   string  `json:"gender"`
}

// 认证响应
type AuthResponse struct {
	Token     string       `json:"token"`
	ExpiresIn int          `json:"expires_in"` // 秒
	User      UserResponse `json:"user"`
}

// 用户资料响应
type UserProfileResponse struct {
	ID           string                        `json:"id"`
	Phone        string                        `json:"phone"`
	Email        *string                       `json:"email,omitempty"`
	Nickname     string                        `json:"nickname"`
	Avatar       string                        `json:"avatar"`
	Age          int                           `json:"age"`
	Gender       string                        `json:"gender"`
	Bio          *string                       `json:"bio,omitempty"`
	Interests    []string                      `json:"interests"`
	Photos       []string                      `json:"photos"`
	Location     *models.GeoLocation           `json:"location,omitempty"`
	Verification models.VerificationStatus     `json:"verification"`
	Preferences  models.UserPreferences        `json:"preferences"`
	Privacy      models.PrivacySettings        `json:"privacy"`
}

// 兴趣标签响应
type InterestTagResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	UsageCount  int    `json:"usage_count"`
}

// API响应包装
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 成功响应
func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// 错误响应
func ErrorResponse(code int, message string) APIResponse {
	return APIResponse{
		Code:    code,
		Message: message,
	}
}

// 资料完整度响应
type ProfileCompletenessResponse struct {
	UserID          string   `json:"user_id"`
	OverallScore    int      `json:"overall_score"`    // 总体完整度百分比
	CompletedFields []string `json:"completed_fields"` // 已完成的字段
	MissingFields   []string `json:"missing_fields"`   // 缺失的字段
	Recommendations []string `json:"recommendations"`  // 完善建议
}

// 用户搜索结果
type UserSearchResult struct {
	ID       string `json:"id"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	Bio      string `json:"bio,omitempty"`
}

// 好友响应
type FriendResponse struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	FriendID string `json:"friend_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Email    string `json:"email,omitempty"`
	Gender   string `json:"gender"`
	Status   string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// 好友请求响应
type FriendRequestResponse struct {
	ID               string `json:"id"`
	FromUserID       string `json:"from_user_id"`
	ToUserID         string `json:"to_user_id"`
	FromUserNickname string `json:"from_user_nickname"`
	FromUserAvatar   string `json:"from_user_avatar"`
	FromUserPhone    string `json:"from_user_phone"`
	Message          string `json:"message,omitempty"`
	Status           string `json:"status"`
	CreatedAt        string `json:"created_at"`
}

// 发送好友请求
type SendFriendRequestRequest struct {
	ToUserID string  `json:"to_user_id" binding:"required"`
	Message  *string `json:"message"`
}

// 好友状态响应
type FriendStatusResponse struct {
	Status         string `json:"status"`          // "friend", "request_sent", "request_received", "none"
	IsFriend       bool   `json:"is_friend"`       // 是否已经是好友
	CanSendRequest bool   `json:"can_send_request"` // 是否可以发送好友请求
	Message        string `json:"message"`         // 状态描述信息
}

// 好友推荐响应
type FriendRecommendationResponse struct {
	UserID          string   `json:"user_id"`
	Nickname        string   `json:"nickname"`
	Avatar          string   `json:"avatar"`
	Gender          string   `json:"gender"`
	Age             int      `json:"age"`
	CommonInterests int      `json:"common_interests"`
	Distance        float64  `json:"distance"`
	Interests       []string `json:"interests"`
	Bio             string   `json:"bio,omitempty"`
}