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