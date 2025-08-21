package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"social-app/services/upload-service/repository"
	"social-app/shared/models"
	"social-app/shared/storage"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PhotoService interface {
	UploadPhoto(req *models.PhotoUploadRequest) (*models.PhotoResponse, error)
	GetUserPhotos(userID string, photoType models.PhotoType) ([]*models.PhotoResponse, error)
	GetPhotoByID(photoID string) (*models.PhotoResponse, error)
	GetPhotoData(photoID string) ([]byte, error)
	DeletePhoto(photoID, userID string) error
	SetActiveAvatar(photoID, userID string) error
	GenerateThumbnail(photoID string) error
	UpdateUserAvatar(userID, avatarURL string) error
	GetUserAvatar(userID string) (string, error)
	UploadFile(file io.Reader, filename, userID, fileType string) (*FileUploadResult, error)
}

// 文件上传结果
type FileUploadResult struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

type photoService struct {
	photoRepo   repository.PhotoRepository
	storagePath string
	redisClient *redis.Client
	obsClient   *storage.OBSClient
	db          *gorm.DB
}

func NewPhotoService(photoRepo repository.PhotoRepository, storagePath string, redisClient *redis.Client, obsClient *storage.OBSClient, db *gorm.DB) PhotoService {
	return &photoService{
		photoRepo:   photoRepo,
		storagePath: storagePath,
		redisClient: redisClient,
		obsClient:   obsClient,
		db:          db,
	}
}

// UploadPhoto 上传照片
func (s *photoService) UploadPhoto(req *models.PhotoUploadRequest) (*models.PhotoResponse, error) {
	// 生成文件名
	filename := s.generateFilename(req.UserID, req.Type, req.Filename)

	// 创建存储目录
	dir := filepath.Join(s.storagePath, string(req.Type))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %v", err)
	}

	// 获取图片尺寸（如果失败则使用默认值）
	width, height := 100, 100
	if len(req.File) > 0 {
		if w, h, err := s.getImageDimensionsFromBytes(req.File); err == nil {
			width, height = w, h
		}
	}

	// 上传到OBS
	obsURL, err := s.uploadToOBS(filename, req.File, req.ContentType)
	if err != nil {
		return nil, fmt.Errorf("上传到OBS失败: %v", err)
	}

	// 创建照片记录
	photo := &models.Photo{
		UserID:       req.UserID,
		Type:         req.Type,
		Filename:     filename,
		OriginalName: req.Filename,
		ContentType:  req.ContentType,
		Size:         int64(len(req.File)),
		Width:        width,
		Height:       height,
		URL:          obsURL, // 使用OBS URL
		IsActive:     true,
		UploadedAt:   time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if req.Metadata != nil {
		photo.Metadata = *req.Metadata
	}

	// 如果是头像，先将其他头像设为非活跃状态
	if req.Type == models.PhotoTypeAvatar {
		if err := s.photoRepo.DeactivateUserPhotos(req.UserID, models.PhotoTypeAvatar); err != nil {
			return nil, fmt.Errorf("更新头像状态失败: %v", err)
		}
	}

	// 保存到数据库
	if err := s.photoRepo.CreatePhoto(photo); err != nil {
		return nil, fmt.Errorf("保存照片记录失败: %v", err)
	}

	// 异步生成缩略图
	go func() {
		if err := s.GenerateThumbnail(photo.ID.Hex()); err != nil {
			fmt.Printf("生成缩略图失败: %v\n", err)
		}
	}()

	return photo.ToResponse(), nil
}

// GetUserPhotos 获取用户照片
func (s *photoService) GetUserPhotos(userID string, photoType models.PhotoType) ([]*models.PhotoResponse, error) {
	photos, err := s.photoRepo.GetUserPhotos(userID, photoType)
	if err != nil {
		return nil, fmt.Errorf("获取用户照片失败: %v", err)
	}

	var responses []*models.PhotoResponse
	for _, photo := range photos {
		responses = append(responses, photo.ToResponse())
	}

	return responses, nil
}

// GetPhotoByID 根据ID获取照片
func (s *photoService) GetPhotoByID(photoID string) (*models.PhotoResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(photoID)
	if err != nil {
		return nil, fmt.Errorf("无效的照片ID: %v", err)
	}

	photo, err := s.photoRepo.GetPhotoByID(objectID)
	if err != nil {
		return nil, fmt.Errorf("获取照片失败: %v", err)
	}

	return photo.ToResponse(), nil
}

// GetPhotoData 获取照片数据（先从Redis查询，再从MongoDB查询）
func (s *photoService) GetPhotoData(photoID string) ([]byte, error) {
	// 先从Redis查询
	cacheKey := fmt.Sprintf("photo:data:%s", photoID)
	if cachedData, err := s.getFromRedis(cacheKey); err == nil && len(cachedData) > 0 {
		return cachedData, nil
	}

	// Redis中没有，从MongoDB查询
	objectID, err := primitive.ObjectIDFromHex(photoID)
	if err != nil {
		return nil, fmt.Errorf("无效的照片ID: %v", err)
	}

	photo, err := s.photoRepo.GetPhotoByID(objectID)
	if err != nil {
		return nil, fmt.Errorf("获取照片失败: %v", err)
	}

	if len(photo.Data) == 0 {
		return nil, fmt.Errorf("照片数据不存在")
	}

	// 同步到Redis（异步）
	go s.setToRedis(cacheKey, photo.Data, 24*time.Hour)

	return photo.Data, nil
}

// DeletePhoto 删除照片
func (s *photoService) DeletePhoto(photoID, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(photoID)
	if err != nil {
		return fmt.Errorf("无效的照片ID: %v", err)
	}

	// 获取照片信息
	photo, err := s.photoRepo.GetPhotoByID(objectID)
	if err != nil {
		return fmt.Errorf("获取照片失败: %v", err)
	}

	// 检查权限
	if photo.UserID != userID {
		return fmt.Errorf("无权限删除此照片")
	}

	// 删除文件
	filePath := filepath.Join(s.storagePath, string(photo.Type), photo.Filename)
	if err := os.Remove(filePath); err != nil {
		fmt.Printf("删除文件失败: %v\n", err) // 不阻断流程
	}

	// 删除缩略图
	if photo.ThumbnailURL != "" {
		thumbnailPath := filepath.Join(s.storagePath, "thumbnails", photo.Filename)
		os.Remove(thumbnailPath)
	}

	// 从数据库删除
	return s.photoRepo.DeletePhoto(objectID)
}

// SetActiveAvatar 设置活跃头像
func (s *photoService) SetActiveAvatar(photoID, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(photoID)
	if err != nil {
		return fmt.Errorf("无效的照片ID: %v", err)
	}

	// 获取照片信息
	photo, err := s.photoRepo.GetPhotoByID(objectID)
	if err != nil {
		return fmt.Errorf("获取照片失败: %v", err)
	}

	// 检查权限和类型
	if photo.UserID != userID {
		return fmt.Errorf("无权限操作此照片")
	}

	if photo.Type != models.PhotoTypeAvatar {
		return fmt.Errorf("只能设置头像类型的照片为活跃状态")
	}

	// 先将其他头像设为非活跃状态
	if err := s.photoRepo.DeactivateUserPhotos(userID, models.PhotoTypeAvatar); err != nil {
		return fmt.Errorf("更新头像状态失败: %v", err)
	}

	// 设置当前头像为活跃状态
	return s.photoRepo.SetPhotoActive(objectID, true)
}

// GenerateThumbnail 生成缩略图
func (s *photoService) GenerateThumbnail(photoID string) error {
	objectID, err := primitive.ObjectIDFromHex(photoID)
	if err != nil {
		return fmt.Errorf("无效的照片ID: %v", err)
	}

	photo, err := s.photoRepo.GetPhotoByID(objectID)
	if err != nil {
		return fmt.Errorf("获取照片失败: %v", err)
	}

	// 对于OBS文件，直接使用原图URL作为缩略图
	if strings.Contains(photo.URL, "obs.") || strings.Contains(photo.URL, "myhuaweicloud.com") {
		return s.photoRepo.UpdateThumbnailURL(objectID, photo.URL)
	}

	// 本地文件的缩略图生成逻辑保持不变
	originalPath := filepath.Join(s.storagePath, string(photo.Type), photo.Filename)
	thumbnailDir := filepath.Join(s.storagePath, "thumbnails")
	if err := os.MkdirAll(thumbnailDir, 0755); err != nil {
		return fmt.Errorf("创建缩略图目录失败: %v", err)
	}
	thumbnailPath := filepath.Join(thumbnailDir, photo.Filename)
	if err := s.createThumbnail(originalPath, thumbnailPath, 200, 200); err != nil {
		return fmt.Errorf("生成缩略图失败: %v", err)
	}
	thumbnailURL := s.generateURL("thumbnails", photo.Filename)
	return s.photoRepo.UpdateThumbnailURL(objectID, thumbnailURL)
}

// 辅助方法

func (s *photoService) generateFilename(userID string, photoType models.PhotoType, originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%s_%d%s", userID, photoType, timestamp, ext)
}

func (s *photoService) generateURL(photoType models.PhotoType, filename string) string {
	return fmt.Sprintf("/uploads/%s/%s", photoType, filename)
}

func (s *photoService) saveFile(data []byte, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func (s *photoService) getImageDimensions(filePath string) (int, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}

func (s *photoService) getImageDimensionsFromBytes(data []byte) (int, int, error) {
	reader := bytes.NewReader(data)
	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		return 0, 0, err
	}
	return img.Width, img.Height, nil
}

func (s *photoService) createThumbnail(srcPath, dstPath string, maxWidth, maxHeight int) error {
	// 打开原图
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 解码图片
	img, _, err := image.Decode(srcFile)
	if err != nil {
		return err
	}

	// 计算缩略图尺寸
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	if width > maxWidth || height > maxHeight {
		ratio := float64(width) / float64(height)
		if width > height {
			width = maxWidth
			height = int(float64(maxWidth) / ratio)
		} else {
			height = maxHeight
			width = int(float64(maxHeight) * ratio)
		}
	}

	// 创建缩略图文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 这里简化处理，实际应该使用图片缩放库
	// 目前直接复制原图作为缩略图
	srcFile.Seek(0, io.SeekStart)
	_, err = io.Copy(dstFile, srcFile)

	return err
}

// Redis操作辅助方法
func (s *photoService) getFromRedis(key string) ([]byte, error) {
	if s.redisClient == nil {
		return nil, fmt.Errorf("Redis客户端未初始化")
	}
	return s.redisClient.Get(context.Background(), key).Bytes()
}

func (s *photoService) setToRedis(key string, data []byte, expiration time.Duration) {
	if s.redisClient == nil {
		return
	}
	ctx := context.Background()
	if err := s.redisClient.Set(ctx, key, data, expiration).Err(); err != nil {
		fmt.Printf("Redis设置失败: %v\n", err)
	}
}

// 上传到OBS
func (s *photoService) uploadToOBS(key string, data []byte, contentType string) (string, error) {
	if s.obsClient == nil {
		return "", fmt.Errorf("OBS客户端未初始化")
	}
	return s.obsClient.UploadFile(fmt.Sprintf("touxiang/%s", key), data, contentType)
}

// UpdateUserAvatar 更新用户头像
func (s *photoService) UpdateUserAvatar(userID, avatarURL string) error {
	if s.db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	// 查找是否已存在用户头像记录
	var avatar models.UserAvatar
	err := s.db.Where("user_id = ?", userID).First(&avatar).Error

	if err == gorm.ErrRecordNotFound {
		// 不存在，创建新记录
		avatar = models.UserAvatar{
			UserID:    userID,
			AvatarURL: avatarURL,
		}
		return s.db.Create(&avatar).Error
	} else if err != nil {
		return fmt.Errorf("查询用户头像失败: %v", err)
	} else {
		// 存在，更新URL
		return s.db.Model(&avatar).Update("avatar_url", avatarURL).Error
	}
}

// GetUserAvatar 获取用户头像
func (s *photoService) GetUserAvatar(userID string) (string, error) {
	if s.db == nil {
		return "", fmt.Errorf("数据库连接未初始化")
	}

	var avatar models.UserAvatar
	err := s.db.Where("user_id = ?", userID).First(&avatar).Error

	if err == gorm.ErrRecordNotFound {
		return "", nil // 没有头像记录
	} else if err != nil {
		return "", fmt.Errorf("查询用户头像失败: %v", err)
	}

	return avatar.AvatarURL, nil
}

// UploadFile 通用文件上传
func (s *photoService) UploadFile(file io.Reader, filename, userID, fileType string) (*FileUploadResult, error) {
	// 读取文件数据
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 生成唯一文件名
	uniqueFilename := s.generateUniqueFilename(userID, fileType, filename)

	// 确定内容类型
	contentType := s.getContentType(filename, fileType)

	// 上传到OBS
	var obsURL string
	if s.obsClient != nil {
		// 根据文件类型选择不同的存储路径
		var obsKey string
		switch fileType {
		case "image":
			obsKey = fmt.Sprintf("images/%s", uniqueFilename)
		case "video":
			obsKey = fmt.Sprintf("videos/%s", uniqueFilename)
		default:
			obsKey = fmt.Sprintf("files/%s", uniqueFilename)
		}

		obsURL, err = s.obsClient.UploadFile(obsKey, data, contentType)
		if err != nil {
			return nil, fmt.Errorf("上传到OBS失败: %v", err)
		}
	} else {
		// 本地存储
		dir := filepath.Join(s.storagePath, fileType)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("创建存储目录失败: %v", err)
		}

		filePath := filepath.Join(dir, uniqueFilename)
		if err := s.saveFile(data, filePath); err != nil {
			return nil, fmt.Errorf("保存文件失败: %v", err)
		}

		obsURL = fmt.Sprintf("/uploads/%s/%s", fileType, uniqueFilename)
	}

	return &FileUploadResult{
		URL:      obsURL,
		Filename: uniqueFilename,
		Size:     int64(len(data)),
	}, nil
}

// 生成唯一文件名
func (s *photoService) generateUniqueFilename(userID, fileType, originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%s_%d%s", userID, fileType, timestamp, ext)
}

// 根据文件名和类型确定内容类型
func (s *photoService) getContentType(filename, fileType string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch fileType {
	case "image":
		switch ext {
		case ".jpg", ".jpeg":
			return "image/jpeg"
		case ".png":
			return "image/png"
		case ".gif":
			return "image/gif"
		case ".webp":
			return "image/webp"
		default:
			return "image/jpeg"
		}
	case "video":
		switch ext {
		case ".mp4":
			return "video/mp4"
		case ".avi":
			return "video/avi"
		case ".mov":
			return "video/quicktime"
		case ".wmv":
			return "video/x-ms-wmv"
		case ".flv":
			return "video/x-flv"
		case ".webm":
			return "video/webm"
		default:
			return "video/mp4"
		}
	default:
		return "application/octet-stream"
	}
}
