package handlers

import (
	"net/http"
	"social-app/services/auth-service/models"
	"social-app/services/auth-service/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Nickname string `json:"nickname"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	result, err := h.authService.Register(req.Username, req.Password, req.Email, req.Nickname)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}

// 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	result, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}

// 验证Token
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	result, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}

// 刷新Token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	result, err := h.authService.RefreshToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}

// 身份认证
func (h *AuthHandler) VerifyIdentity(c *gin.Context) {
	var req models.IdentityVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	result, err := h.authService.VerifyIdentity(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}

// 人脸认证
func (h *AuthHandler) VerifyFace(c *gin.Context) {
	var req models.FaceVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	result, err := h.authService.VerifyFace(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}

// 活体检测
func (h *AuthHandler) VerifyLiveness(c *gin.Context) {
	var req models.LivenessDetectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	result, err := h.authService.VerifyLiveness(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}

// OCR身份证识别
func (h *AuthHandler) OCRIDCard(c *gin.Context) {
	var req struct {
		IDPhoto string `json:"id_photo" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	// 模拟OCR识别
	result := &models.OCRResult{
		Name:     "张三",
		IDNumber: "110101199001011234",
		Address:  "北京市朝阳区",
		Birth:    "1990-01-01",
		Gender:   "男",
		Nation:   "汉",
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}



// 获取认证状态
func (h *AuthHandler) GetVerificationStatus(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "用户ID不能为空"))
		return
	}

	// 模拟获取认证状态
	status := map[string]interface{}{
		"user_id": userID,
		"identity_verified": true,
		"face_verified": true,
		"liveness_verified": true,
		"verification_level": "high",
	}

	c.JSON(http.StatusOK, service.SuccessResponse(status))
}

// 重试认证
func (h *AuthHandler) RetryVerification(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "用户ID不能为空"))
		return
	}

	var req struct {
		VerificationType string `json:"verification_type" binding:"required"` // identity, face, liveness
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	// 模拟重试认证
	result := map[string]interface{}{
		"user_id": userID,
		"verification_type": req.VerificationType,
		"retry_allowed": true,
		"message": "认证重试已启动",
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}

// 上传头像
func (h *AuthHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "上传文件错误: "+err.Error()))
		return
	}

	// 检查文件大小（限制5MB）
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "文件大小超过限制（5MB）"))
		return
	}

	// 检查文件类型
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "不支持的文件类型，仅支持 JPEG、PNG、GIF"))
		return
	}

	// 模拟保存文件并返回URL
	result := map[string]interface{}{
		"filename": file.Filename,
		"size":     file.Size,
		"url":      "https://example.com/avatars/" + file.Filename,
		"message":  "头像上传成功",
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}

// 上传照片
func (h *AuthHandler) UploadPhoto(c *gin.Context) {
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "上传文件错误: "+err.Error()))
		return
	}

	// 检查文件大小（限制10MB）
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "文件大小超过限制（10MB）"))
		return
	}

	// 检查文件类型
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "不支持的文件类型，仅支持 JPEG、PNG、GIF"))
		return
	}

	// 模拟保存文件并返回URL
	result := map[string]interface{}{
		"filename": file.Filename,
		"size":     file.Size,
		"url":      "https://example.com/photos/" + file.Filename,
		"message":  "照片上传成功",
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}

// 微信登录
func (h *AuthHandler) WeChatLogin(c *gin.Context) {
	url := h.authService.GetWeChatLoginURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// 微信回调
func (h *AuthHandler) WeChatCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "缺少code参数"))
		return
	}

	result, err := h.authService.HandleWeChatCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}

// QQ登录
func (h *AuthHandler) QQLogin(c *gin.Context) {
	url := h.authService.GetQQLoginURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// QQ回调
func (h *AuthHandler) QQCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "缺少code参数"))
		return
	}

	result, err := h.authService.HandleQQCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}

// Apple登录
func (h *AuthHandler) AppleLogin(c *gin.Context) {
	url := h.authService.GetAppleLoginURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Apple回调
func (h *AuthHandler) AppleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "缺少code参数"))
		return
	}

	result, err := h.authService.HandleAppleCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(result))
}