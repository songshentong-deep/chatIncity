package config

import (
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	JWT      JWTConfig      `json:"jwt"`
	Upload   UploadConfig   `json:"upload"`
	External ExternalConfig `json:"external"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `json:"port"`
	Mode string `json:"mode"` // debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	PostgreSQL PostgreSQLConfig `json:"postgresql"`
	MongoDB    MongoDBConfig    `json:"mongodb"`
}

// PostgreSQLConfig PostgreSQL配置
type PostgreSQLConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	SSLMode  string `json:"ssl_mode"`
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	URI      string `json:"uri"`
	Database string `json:"database"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `json:"secret"`
	ExpireTime int    `json:"expire_time"` // 小时
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	MaxSize    int64  `json:"max_size"`    // 最大文件大小(字节)
	AllowTypes []string `json:"allow_types"` // 允许的文件类型
	StoragePath string `json:"storage_path"` // 存储路径
}

// ExternalConfig 外部服务配置
type ExternalConfig struct {
	FaceAPI    FaceAPIConfig    `json:"face_api"`
	MapAPI     MapAPIConfig     `json:"map_api"`
	PushAPI    PushAPIConfig    `json:"push_api"`
	ContentAPI ContentAPIConfig `json:"content_api"`
}

// FaceAPIConfig 人脸识别API配置
type FaceAPIConfig struct {
	Provider string `json:"provider"` // 'baidu', 'tencent', 'aliyun'
	APIKey   string `json:"api_key"`
	Secret   string `json:"secret"`
	Endpoint string `json:"endpoint"`
}

// MapAPIConfig 地图API配置
type MapAPIConfig struct {
	Provider string `json:"provider"` // 'amap', 'baidu', 'google'
	APIKey   string `json:"api_key"`
	Secret   string `json:"secret"`
}

// PushAPIConfig 推送API配置
type PushAPIConfig struct {
	Provider string `json:"provider"` // 'jpush', 'firebase'
	AppKey   string `json:"app_key"`
	Secret   string `json:"secret"`
}

// ContentAPIConfig 内容审核API配置
type ContentAPIConfig struct {
	Provider string `json:"provider"` // 'baidu', 'tencent', 'aliyun'
	APIKey   string `json:"api_key"`
	Secret   string `json:"secret"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 加载.env文件
	godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			PostgreSQL: PostgreSQLConfig{
				Host:     getEnv("POSTGRES_HOST", "localhost"),
				Port:     getEnv("POSTGRES_PORT", "5432"),
				User:     getEnv("POSTGRES_USER", "postgres"),
				Password: getEnv("POSTGRES_PASSWORD", "password"),
				DBName:   getEnv("POSTGRES_DB", "social_app"),
				SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
			},
			MongoDB: MongoDBConfig{
				URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
				Database: getEnv("MONGODB_DATABASE", "social_app"),
			},
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key"),
			ExpireTime: getEnvAsInt("JWT_EXPIRE_TIME", 24),
		},
		Upload: UploadConfig{
			MaxSize:     getEnvAsInt64("UPLOAD_MAX_SIZE", 10*1024*1024), // 10MB
			AllowTypes:  []string{"image/jpeg", "image/png", "image/gif", "video/mp4"},
			StoragePath: getEnv("UPLOAD_STORAGE_PATH", "./uploads"),
		},
		External: ExternalConfig{
			FaceAPI: FaceAPIConfig{
				Provider: getEnv("FACE_API_PROVIDER", "baidu"),
				APIKey:   getEnv("FACE_API_KEY", ""),
				Secret:   getEnv("FACE_API_SECRET", ""),
				Endpoint: getEnv("FACE_API_ENDPOINT", ""),
			},
			MapAPI: MapAPIConfig{
				Provider: getEnv("MAP_API_PROVIDER", "amap"),
				APIKey:   getEnv("MAP_API_KEY", ""),
				Secret:   getEnv("MAP_API_SECRET", ""),
			},
			PushAPI: PushAPIConfig{
				Provider: getEnv("PUSH_API_PROVIDER", "jpush"),
				AppKey:   getEnv("PUSH_API_KEY", ""),
				Secret:   getEnv("PUSH_API_SECRET", ""),
			},
			ContentAPI: ContentAPIConfig{
				Provider: getEnv("CONTENT_API_PROVIDER", "baidu"),
				APIKey:   getEnv("CONTENT_API_KEY", ""),
				Secret:   getEnv("CONTENT_API_SECRET", ""),
			},
		},
	}

	return config, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为int
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsInt64 获取环境变量并转换为int64
func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}