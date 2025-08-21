package models

import "time"

// AuthResult 认证结果
type AuthResult struct {
	Token     string   `json:"token"`
	User      UserInfo `json:"user"`
	ExpiresIn int      `json:"expires_in"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

// TokenValidationResult token验证结果
type TokenValidationResult struct {
	Valid     bool      `json:"valid"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// TokenResponse token响应
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

// IdentityVerificationRequest 身份验证请求
type IdentityVerificationRequest struct {
	RealName string `json:"real_name"`
	IDNumber string `json:"id_number"`
	IDPhoto  string `json:"id_photo"`
}

// VerificationResult 验证结果
type VerificationResult struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// FaceVerificationRequest 人脸验证请求
type FaceVerificationRequest struct {
	FaceImage      string `json:"face_image"`
	ReferenceImage string `json:"reference_image"`
}

// FaceVerificationResult 人脸验证结果
type FaceVerificationResult struct {
	Success    bool                   `json:"success"`
	Confidence float64                `json:"confidence"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details"`
}

// LivenessDetectionRequest 活体检测请求
type LivenessDetectionRequest struct {
	VideoFrames []string `json:"video_frames"`
}

// LivenessResult 活体检测结果
type LivenessResult struct {
	IsLive     bool      `json:"is_live"`
	Score      float64   `json:"score"`
	Message    string    `json:"message"`
	Actions    []string  `json:"actions"`
	VerifiedAt time.Time `json:"verified_at"`
}

// OCRRequest OCR请求
type OCRRequest struct {
	Image string `json:"image"`
}

// OCRResult OCR结果
type OCRResult struct {
	Name     string `json:"name"`
	IDNumber string `json:"id_number"`
	Address  string `json:"address"`
	Birth    string `json:"birth"`
	Gender   string `json:"gender"`
	Nation   string `json:"nation"`
}

// OAuthResult OAuth结果
type OAuthResult struct {
	Provider    string        `json:"provider"`
	OpenID      string        `json:"open_id"`
	UnionID     string        `json:"union_id,omitempty"`
	UserInfo    OAuthUserInfo `json:"user_info"`
	AccessToken string        `json:"access_token"`
	ExpiresIn   int           `json:"expires_in"`
}

// OAuthUserInfo OAuth用户信息
type OAuthUserInfo struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   string `json:"gender"`
}

// RealNameResult 实名验证结果
type RealNameResult struct {
	Valid      bool    `json:"valid"`
	Message    string  `json:"message"`
	Confidence float64 `json:"confidence"`
}

// FaceRecognitionResult 人脸识别结果
type FaceRecognitionResult struct {
	Confidence    float64 `json:"confidence"`
	Message       string  `json:"message"`
	FaceQuality   float64 `json:"face_quality"`
	LivenessScore float64 `json:"liveness_score"`
}

// LivenessDetectionResult 活体检测结果
type LivenessDetectionResult struct {
	Score           float64  `json:"score"`
	Message         string   `json:"message"`
	DetectedActions []string `json:"detected_actions"`
}