package service

// 创建动态请求
type CreateMomentRequest struct {
	UserID   string    `json:"user_id"`
	Content  string    `json:"content"`
	Images   *[]string `json:"images,omitempty"`    // 图片URL列表（已上传的）
	VideoURL *string   `json:"video_url,omitempty"` // 视频URL（已上传的）
	Type     string    `json:"type" binding:"required,oneof=text image video"`
}

// 动态响应
type MomentResponse struct {
	ID            string   `json:"id"`
	UserID        string   `json:"user_id"`
	UserNickname  string   `json:"user_nickname"`
	UserAvatar    string   `json:"user_avatar"`
	Content       string   `json:"content"`
	Images        []string `json:"images"`
	VideoURL      *string  `json:"video_url"`
	Type          string   `json:"type"`
	LikesCount    int      `json:"likes_count"`
	CommentsCount int      `json:"comments_count"`
	IsLiked       bool     `json:"is_liked"`
	CreatedAt     string   `json:"created_at"`
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
