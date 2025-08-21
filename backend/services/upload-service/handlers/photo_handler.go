package handlers

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"social-app/services/upload-service/service"
	"social-app/shared/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type PhotoHandler struct {
	photoService service.PhotoService
}

func NewPhotoHandler(photoService service.PhotoService) *PhotoHandler {
	return &PhotoHandler{
		photoService: photoService,
	}
}

// UploadPhoto 上传照片
func (h *PhotoHandler) UploadPhoto(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取照片类型
	photoType := c.PostForm("type")
	if photoType == "" {
		photoType = string(models.PhotoTypeGallery)
	}

	// 验证照片类型
	if photoType != string(models.PhotoTypeAvatar) && photoType != string(models.PhotoTypeGallery) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的照片类型",
		})
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "获取上传文件失败: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// 验证文件类型
	if !isValidImageFile(header.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不支持的文件格式，请上传jpg、png或gif格式的图片",
		})
		return
	}

	// 验证文件大小 (10MB)
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件大小不能超过10MB",
		})
		return
	}

	// 读取文件内容
	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "读取文件失败: " + err.Error(),
		})
		return
	}

	// 创建上传请求
	uploadReq := &models.PhotoUploadRequest{
		UserID:      userID.(string),
		Type:        models.PhotoType(photoType),
		File:        fileData,
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
	}

	// 上传照片
	photo, err := h.photoService.UploadPhoto(uploadReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "上传照片失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data":    photo,
	})
}

// GetUserPhotos 获取用户照片
func (h *PhotoHandler) GetUserPhotos(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取照片类型
	photoType := c.Query("type")
	if photoType == "" {
		photoType = string(models.PhotoTypeGallery)
	}

	// 验证照片类型
	if photoType != string(models.PhotoTypeAvatar) && photoType != string(models.PhotoTypeGallery) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的照片类型",
		})
		return
	}

	photos, err := h.photoService.GetUserPhotos(userID.(string), models.PhotoType(photoType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取照片失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    photos,
	})
}

// GetPhoto 获取单张照片
func (h *PhotoHandler) GetPhoto(c *gin.Context) {
	photoID := c.Param("id")
	if photoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "照片ID不能为空",
		})
		return
	}

	photo, err := h.photoService.GetPhotoByID(photoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "照片不存在: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    photo,
	})
}

// GetPhotoData 获取图片数据
func (h *PhotoHandler) GetPhotoData(c *gin.Context) {
	photoID := c.Param("id")
	if photoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "照片ID不能为空",
		})
		return
	}

	// 获取照片信息
	_, err := h.photoService.GetPhotoByID(photoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "照片不存在: " + err.Error(),
		})
		return
	}

	// 从数据库获取图片数据
	photoData, err := h.photoService.GetPhotoData(photoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取图片数据失败: " + err.Error(),
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "image/jpeg")
	c.Header("Content-Length", fmt.Sprintf("%d", len(photoData)))
	c.Header("Cache-Control", "public, max-age=86400") // 缓存1天

	// 返回图片数据
	c.Data(http.StatusOK, "image/jpeg", photoData)
}

// UploadPhotoBytes Web平台字节上传
func (h *PhotoHandler) UploadPhotoBytes(c *gin.Context) {
	// 从认证中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    400,
			"message": "用户id不存在",
		})
		return
	}

	// 读取请求体
	fileData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "读取文件数据失败: " + err.Error(),
		})
		return
	}

	// 获取头信息
	photoType := c.GetHeader("X-Photo-Type")
	if photoType == "" {
		photoType = "avatar"
	}
	filename := c.GetHeader("X-Filename")
	if filename == "" {
		filename = "avatar.jpg"
	}
	contentType := c.GetHeader("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	// 创建上传请求
	uploadReq := &models.PhotoUploadRequest{
		UserID:      userID.(string),
		Type:        models.PhotoType(photoType),
		File:        fileData,
		Filename:    filename,
		ContentType: contentType,
	}

	// 上传照片到OBS
	photo, err := h.photoService.UploadPhoto(uploadReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "上传照片失败: " + err.Error(),
		})
		return
	}

	// 如果是头像，更新用户头像表
	if photoType == "avatar" {
		err = h.photoService.UpdateUserAvatar(userID.(string), photo.URL)
		if err != nil {
			fmt.Printf("更新用户头像失败: %v\n", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data":    photo,
	})
}

// DeletePhoto 删除照片
func (h *PhotoHandler) DeletePhoto(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	photoID := c.Param("id")
	if photoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "照片ID不能为空",
		})
		return
	}

	err := h.photoService.DeletePhoto(photoID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "删除照片失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// SetActiveAvatar 设置活跃头像
func (h *PhotoHandler) SetActiveAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	photoID := c.Param("id")
	if photoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "照片ID不能为空",
		})
		return
	}

	err := h.photoService.SetActiveAvatar(photoID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "设置头像失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "设置成功",
	})
}

// GetUserAvatar 获取用户头像
func (h *PhotoHandler) GetUserAvatar(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	avatarURL, err := h.photoService.GetUserAvatar(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户头像失败: " + err.Error(),
		})
		return
	}

	// 如果是OBS URL，转换为代理URL
	if avatarURL != "" && (strings.Contains(avatarURL, "obs.") || strings.Contains(avatarURL, "myhuaweicloud.com")) {
		avatarURL = fmt.Sprintf("/api/v1/proxy/image?url=%s", avatarURL)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"avatar_url": avatarURL,
		},
	})
}

// ProxyImage 图片代理接口
func (h *PhotoHandler) ProxyImage(c *gin.Context) {
	imageURL := c.Query("url")
	if imageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少图片URL参数",
		})
		return
	}

	// 只允许代理OBS图片
	if !strings.Contains(imageURL, "obs.") && !strings.Contains(imageURL, "myhuaweicloud.com") {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不允许的图片源",
		})
		return
	}

	// 请求图片
	resp, err := http.Get(imageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取图片失败",
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// 返回默认头像
		c.Header("Content-Type", "image/svg+xml")
		c.Header("Access-Control-Allow-Origin", "*")
		defaultAvatar := `<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg"><circle cx="50" cy="50" r="40" fill="#ccc"/></svg>`
		c.String(http.StatusOK, defaultAvatar)
		return
	}

	// 设置响应头
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Status(http.StatusOK)

	// 转发图片数据
	io.Copy(c.Writer, resp.Body)
}

// 验证图片文件格式
func isValidImageFile(filename string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	ext := strings.ToLower(filepath.Ext(filename))
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}
