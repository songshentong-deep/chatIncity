package handlers

import (
	"fmt"
	"net/http"
	"social-app/services/user-service/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	response, err := h.userService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(response))
}

// 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	response, err := h.userService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(response))
}

// 刷新token
func (h *UserHandler) RefreshToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "缺少Authorization头"))
		return
	}

	// 移除Bearer前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	response, err := h.userService.RefreshToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(response))
}

// 获取用户资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	response, err := h.userService.GetProfile(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, service.ErrorResponse(404, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(response))
}

// 更新用户资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	err := h.userService.UpdateProfile(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "资料更新成功"}))
}

// 上传头像
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "上传文件错误: "+err.Error()))
		return
	}

	// 这里应该实现文件上传到云存储的逻辑
	// 简化处理，返回一个模拟的URL
	avatarURL := "https://example.com/avatars/" + userID.(string) + "_" + file.Filename

	// 更新用户头像
	err = h.userService.UpdateProfile(userID.(string), &service.UpdateProfileRequest{
		Avatar: avatarURL,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{
		"avatar_url": avatarURL,
		"message":    "头像上传成功",
	}))
}

// 身份认证
func (h *UserHandler) VerifyIdentity(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	var req service.IdentityVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	err := h.userService.VerifyIdentity(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "身份认证成功"}))
}

// 人脸认证
func (h *UserHandler) VerifyFace(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	var req service.FaceVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	err := h.userService.VerifyFace(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "人脸认证成功"}))
}

// 获取兴趣标签
func (h *UserHandler) GetInterestTags(c *gin.Context) {
	tags, err := h.userService.GetInterestTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(tags))
}

// 上传多张照片
func (h *UserHandler) UploadPhotos(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	// 获取上传的文件
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "获取上传文件失败: "+err.Error()))
		return
	}

	files := form.File["photos"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "没有上传文件"))
		return
	}

	// 限制最多上传6张照片
	if len(files) > 6 {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "最多只能上传6张照片"))
		return
	}

	var photoURLs []string
	for i, file := range files {
		// 这里应该实现文件上传到云存储的逻辑
		// 简化处理，返回一个模拟的URL
		photoURL := fmt.Sprintf("https://example.com/photos/%s_%d_%s", userID.(string), i, file.Filename)
		photoURLs = append(photoURLs, photoURL)
	}

	// 更新用户照片
	err = h.userService.UpdatePhotos(userID.(string), photoURLs)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{
		"photo_urls": photoURLs,
		"message":    "照片上传成功",
	}))
}

// 删除照片
func (h *UserHandler) DeletePhoto(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	photoID := c.Param("photo_id")
	if photoID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "照片ID不能为空"))
		return
	}

	err := h.userService.DeletePhoto(userID.(string), photoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "照片删除成功"}))
}

// 更新兴趣标签
func (h *UserHandler) UpdateInterests(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	var req struct {
		Interests []string `json:"interests" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	// 限制兴趣标签数量
	if len(req.Interests) > 10 {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "最多只能选择10个兴趣标签"))
		return
	}

	err := h.userService.UpdateInterests(userID.(string), req.Interests)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "兴趣标签更新成功"}))
}

// 获取资料完整度
func (h *UserHandler) GetProfileCompleteness(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	completeness, err := h.userService.GetProfileCompleteness(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(completeness))
}
