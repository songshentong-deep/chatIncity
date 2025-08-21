package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// 上传图片
func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}
	defer file.Close()

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的图片格式"})
		return
	}

	// 生成文件名
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	filename := timestamp + ext

	// 模拟上传到OBS
	obsURL := fmt.Sprintf("https://chat-892c.obs.cn-south-4.myhuaweicloud.com/dongtai/tupian/%s", filename)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data": gin.H{
			"url": obsURL,
		},
	})
}

// 上传视频
func (h *UploadHandler) UploadVideo(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}
	defer file.Close()

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".mp4" && ext != ".avi" && ext != ".mov" && ext != ".wmv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的视频格式"})
		return
	}

	// 生成文件名
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	filename := timestamp + ext

	// 模拟上传到OBS
	obsURL := fmt.Sprintf("https://chat-892c.obs.cn-south-4.myhuaweicloud.com/dongtai/shipin/%s", filename)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传成功",
		"data": gin.H{
			"url": obsURL,
		},
	})
}
