package service

import "time"

// Token验证结果
type TokenValidationResult struct {
	Valid     bool      `json:"valid"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Token响应
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

// 身份认证请求
type IdentityVerificationRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	RealName string `json:"real_name" binding:"required"`
	IDNumber string `json:"id_number" binding:"required"`
	IDPhoto  string `json:"id_photo" binding:"required"` // base64编码的身份证照片
}

// 人脸认证请求
type FaceVerificationRequest struct {
	UserID         string `json:"user_id" binding:"required"`
	FaceImage      string `json:"face_image" binding:"required"`      // base64编码的人脸照片
	ReferenceImage string `json:"reference_image,omitempty"`          // 参考照片（可选）
}

// 活体检测请求
type LivenessDetectionRequest struct {
	UserID      string   `json:"user_id" binding:"required"`
	VideoFrames []string `json:"video_frames" binding:"required"` // base64编码的视频帧
}

// 验证结果
type VerificationResult struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// 人脸验证结果
type FaceVerificationResult struct {
	Success    bool                   `json:"success"`
	Confidence float64                `json:"confidence"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// 活体检测结果
type LivenessResult struct {
	IsLive     bool      `json:"is_live"`
	Score      float64   `json:"score"`
	Message    string    `json:"message"`
	Actions    []string  `json:"actions"`
	VerifiedAt time.Time `json:"verified_at"`
}

// OAuth结果
type OAuthResult struct {
	Provider    string        `json:"provider"`
	OpenID      string        `json:"open_id"`
	UnionID     string        `json:"union_id,omitempty"`
	UserInfo    OAuthUserInfo `json:"user_info"`
	AccessToken string        `json:"access_token"`
	ExpiresIn   int           `json:"expires_in"`
}

// OAuth用户信息
type OAuthUserInfo struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   string `json:"gender"`
}

// OCR请求
type OCRRequest struct {
	Image string `json:"image" binding:"required"` // base64编码的图片数据
}

// OCR结果
type OCRResult struct {
	Name     string `json:"name"`
	IDNumber string `json:"id_number"`
	Address  string `json:"address"`
	Birth    string `json:"birth"`
	Gender   string `json:"gender"`
	Nation   string `json:"nation"`
}

// 实名验证结果
type RealNameResult struct {
	Valid      bool    `json:"valid"`
	Message    string  `json:"message"`
	Confidence float64 `json:"confidence"`
}

// 人脸识别结果
type FaceRecognitionResult struct {
	Confidence    float64 `json:"confidence"`
	Message       string  `json:"message"`
	FaceQuality   float64 `json:"face_quality"`
	LivenessScore float64 `json:"liveness_score"`
}

// 活体检测结果
type LivenessDetectionResult struct {
	Score           float64  `json:"score"`
	Message         string   `json:"message"`
	DetectedActions []string `json:"detected_actions"`
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