package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"social-app/services/user-service/service"
	"strconv"
	"strings"
	"time"

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
	print("获取用户资料被执行了")
	if !exists {
		print("获取用户资料被执行了")
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
	print(userID)
	if !exists {
		print("更新用户资料被执行了")
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

	// 获取更新后的用户资料
	updatedProfile, err := h.userService.GetProfile(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, "获取更新后的资料失败"))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(updatedProfile))
}

// 上传头像
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		print("上传头像被执行了")
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "上传文件错误: "+err.Error()))
		return
	}

	// 验证文件类型
	if !isValidImageFile(file.Filename) {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "不支持的文件格式，请上传jpg、png或gif格式的图片"))
		return
	}

	// 验证文件大小 (5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "文件大小不能超过5MB"))
		return
	}

	// 生成文件名
	filename := fmt.Sprintf("avatar_%s_%d.jpg", userID.(string), time.Now().Unix())

	// 创建上传目录
	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, "创建上传目录失败"))
		return
	}

	// 保存文件
	filepath := fmt.Sprintf("%s/%s", uploadDir, filename)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, "保存文件失败: "+err.Error()))
		return
	}

	// 生成访问URL (在生产环境中应该使用CDN或云存储)
	avatarURL := fmt.Sprintf("http://localhost:8001/uploads/avatars/%s", filename)

	// 更新用户头像
	err = h.userService.UpdateProfile(userID.(string), &service.UpdateProfileRequest{
		Avatar: avatarURL,
	})
	if err != nil {
		// 如果更新失败，删除已上传的文件
		os.Remove(filepath)
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{
		"url": avatarURL,
	}))
}

// 验证图片文件格式
func isValidImageFile(filename string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	ext := strings.ToLower(filepath.Ext(filename))
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

// 身份认证
func (h *UserHandler) VerifyIdentity(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		print("身份认证被执行了")
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

// 搜索用户
func (h *UserHandler) SearchUsers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "搜索关键词不能为空"))
		return
	}

	results, err := h.userService.SearchUsers(query, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(results))
}

// 发送好友请求
func (h *UserHandler) SendFriendRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	var req service.SendFriendRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	// 检查不能添加自己为好友
	if userID.(string) == req.ToUserID {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "不能添加自己为好友"))
		return
	}

	err := h.userService.SendFriendRequest(userID.(string), req.ToUserID, req.Message)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "好友请求发送成功"}))
}

// 接受好友请求
func (h *UserHandler) AcceptFriendRequest(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	requestID := c.Param("request_id")
	if requestID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求ID不能为空"))
		return
	}

	err := h.userService.AcceptFriendRequest(requestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "好友请求已接受"}))
}

// 拒绝好友请求
func (h *UserHandler) RejectFriendRequest(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	requestID := c.Param("request_id")
	if requestID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求ID不能为空"))
		return
	}

	err := h.userService.RejectFriendRequest(requestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "好友请求已拒绝"}))
}

// 获取好友列表
func (h *UserHandler) GetFriends(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	friends, err := h.userService.GetFriends(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(friends))
}

// 获取好友请求列表
func (h *UserHandler) GetFriendRequests(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	requests, err := h.userService.GetFriendRequests(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(requests))
}

// 获取发送的好友请求列表
func (h *UserHandler) GetSentFriendRequests(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	requests, err := h.userService.GetSentFriendRequests(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(requests))
}

// 删除好友
func (h *UserHandler) DeleteFriend(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	friendID := c.Param("friend_id")
	if friendID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "好友ID不能为空"))
		return
	}

	err := h.userService.DeleteFriend(userID.(string), friendID)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "好友删除成功"}))
}

// 检查好友状态
func (h *UserHandler) CheckFriendStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	targetUserID := c.Param("user_id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "用户ID不能为空"))
		return
	}

	status, err := h.userService.CheckFriendStatus(userID.(string), targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(status))
}

// 获取指定用户资料
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "用户ID不能为空"))
		return
	}

	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, service.ErrorResponse(404, "用户不存在"))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(profile))
}

// 获取用户生活动态
func (h *UserHandler) GetUserMoments(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "用户ID不能为空"))
		return
	}

	// 返回空数组，后续实现动态功能
	c.JSON(http.StatusOK, service.SuccessResponse([]interface{}{}))
}

// 更改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	err := h.userService.ChangePassword(req.Phone, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "密码修改成功"}))
}

// 保存最近登录账号
func (h *UserHandler) SaveRecentAccount(c *gin.Context) {
	var req service.SaveRecentAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, "请求参数错误: "+err.Error()))
		return
	}

	err := h.userService.SaveRecentAccount(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(gin.H{"message": "保存成功"}))
}

// 获取最近登录账号
func (h *UserHandler) GetRecentAccounts(c *gin.Context) {
	accounts, err := h.userService.GetRecentAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(accounts))
}

// 获取推荐好友
func (h *UserHandler) GetFriendRecommendations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	recommendations, err := h.userService.GetFriendRecommendations(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(recommendations))
}

// 获取随机用户
func (h *UserHandler) GetRandomUsers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse(401, "未授权"))
		return
	}

	// 获取限制数量参数，默认为4
	limit := 4
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	users, err := h.userService.GetRandomUsers(userID.(string), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse(users))
}
