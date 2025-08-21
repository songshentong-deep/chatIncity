package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"social-app/services/upload-service/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	photoService service.PhotoService
}

func NewUploadHandler(photoService service.PhotoService) *UploadHandler {
	return &UploadHandler{
		photoService: photoService,
	}
}

// 上传图片
func (h *UploadHandler) UploadImage(c *gin.Context) {
	h.uploadFile(c, "image", []string{".jpg", ".jpeg", ".png", ".gif", ".webp"})
}

// 上传视频
func (h *UploadHandler) UploadVideo(c *gin.Context) {
	h.uploadFile(c, "video", []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm"})
}

// 通用文件上传
func (h *UploadHandler) UploadFile(c *gin.Context) {
	h.uploadFile(c, "file", []string{}) // 空数组表示允许所有文件类型
}

// 通用上传处理逻辑
func (h *UploadHandler) uploadFile(c *gin.Context, fileType string, allowedExts []string) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件上传失败: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// 验证文件类型
	if len(allowedExts) > 0 {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		allowed := false
		for _, allowedExt := range allowedExts {
			if ext == allowedExt {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": fmt.Sprintf("不支持的文件类型，仅支持: %s", strings.Join(allowedExts, ", ")),
			})
			return
		}
	}

	// 验证文件大小
	maxSize := h.getMaxFileSize(fileType)
	if header.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("文件大小超过限制，最大允许 %dMB", maxSize/(1024*1024)),
		})
		return
	}

	// 生成唯一文件名
	fileID := uuid.New().String()
	ext := filepath.Ext(header.Filename)
	filename := fileID + ext

	// 上传文件
	uploadResult, err := h.photoService.UploadFile(file, filename, userID.(string), fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "文件上传失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data": gin.H{
			"id":          fileID,
			"filename":    header.Filename,
			"url":         uploadResult.URL,
			"size":        header.Size,
			"type":        fileType,
			"uploaded_at": time.Now().Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}

// 获取文件
func (h *UploadHandler) GetFile(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件ID不能为空",
		})
		return
	}

	// 这里可以实现文件信息获取逻辑
	// 暂时返回基本信息
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"id":  fileID,
			"url": fmt.Sprintf("/uploads/%s", fileID),
		},
	})
}

// 图片代理接口 - 解决CORS问题
func (h *UploadHandler) ProxyImage(c *gin.Context) {
	imageURL := c.Query("url")
	if imageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "图片URL不能为空",
		})
		return
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 请求原始图片
	resp, err := client.Get(imageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取图片失败: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    502,
			"message": fmt.Sprintf("图片服务器返回错误: %d", resp.StatusCode),
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Header("Content-Length", resp.Header.Get("Content-Length"))
	c.Header("Cache-Control", "public, max-age=3600") // 缓存1小时
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET")
	c.Header("Access-Control-Allow-Headers", "Content-Type")

	// 直接将图片数据流式传输给客户端
	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

// 获取不同文件类型的最大大小限制
func (h *UploadHandler) getMaxFileSize(fileType string) int64 {
	switch fileType {
	case "image":
		return 10 * 1024 * 1024 // 10MB
	case "video":
		return 100 * 1024 * 1024 // 100MB
	default:
		return 50 * 1024 * 1024 // 50MB
	}
}

// 验证文件是否为图片
func (h *UploadHandler) isImageFile(header *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"}
	
	for _, imageExt := range imageExts {
		if ext == imageExt {
			return true
		}
	}
	return false
}

// 验证文件是否为视频
func (h *UploadHandler) isVideoFile(header *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	videoExts := []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm", ".mkv"}
	
	for _, videoExt := range videoExts {
		if ext == videoExt {
			return true
		}
	}
	return false
}