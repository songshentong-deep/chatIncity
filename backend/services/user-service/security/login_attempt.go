package security

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// LoginAttemptService 登录尝试限制服务
type LoginAttemptService struct {
	redis *redis.Client
}

// NewLoginAttemptService 创建登录尝试服务
func NewLoginAttemptService(redisClient *redis.Client) *LoginAttemptService {
	return &LoginAttemptService{
		redis: redisClient,
	}
}

// CheckUserAttempts 检查用户登录尝试次数
func (s *LoginAttemptService) CheckUserAttempts(identifier string) error {
	ctx := context.Background()
	key := fmt.Sprintf("user_attempts:%s", identifier)
	
	attempts, err := s.redis.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return err
	}
	
	if attempts >= 5 {
		ttl := s.redis.TTL(ctx, key).Val()
		minutes := int(ttl.Minutes())
		if minutes <= 0 {
			minutes = 15
		}
		return fmt.Errorf("登录尝试次数过多，请在%d分钟后重试", minutes)
	}
	
	return nil
}

// CheckIPAttempts 检查IP登录尝试次数
func (s *LoginAttemptService) CheckIPAttempts(ip string) error {
	ctx := context.Background()
	key := fmt.Sprintf("ip_attempts:%s", ip)
	
	attempts, err := s.redis.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return err
	}
	
	if attempts >= 10 {
		ttl := s.redis.TTL(ctx, key).Val()
		minutes := int(ttl.Minutes())
		if minutes <= 0 {
			minutes = 60
		}
		return fmt.Errorf("该IP登录尝试次数过多，请在%d分钟后重试", minutes)
	}
	
	return nil
}

// RecordFailedAttempt 记录失败的登录尝试
func (s *LoginAttemptService) RecordFailedAttempt(identifier, ip string) error {
	ctx := context.Background()
	
	// 记录用户级别的失败尝试
	userKey := fmt.Sprintf("user_attempts:%s", identifier)
	pipe := s.redis.Pipeline()
	pipe.Incr(ctx, userKey)
	pipe.Expire(ctx, userKey, 15*time.Minute) // 15分钟后重置
	
	// 记录IP级别的失败尝试
	ipKey := fmt.Sprintf("ip_attempts:%s", ip)
	pipe.Incr(ctx, ipKey)
	pipe.Expire(ctx, ipKey, 60*time.Minute) // 1小时后重置
	
	_, err := pipe.Exec(ctx)
	return err
}

// ClearUserAttempts 清除用户的失败尝试记录
func (s *LoginAttemptService) ClearUserAttempts(identifier string) error {
	ctx := context.Background()
	key := fmt.Sprintf("user_attempts:%s", identifier)
	return s.redis.Del(ctx, key).Err()
}

// GetUserAttempts 获取用户当前尝试次数
func (s *LoginAttemptService) GetUserAttempts(identifier string) (int, error) {
	ctx := context.Background()
	key := fmt.Sprintf("user_attempts:%s", identifier)
	
	attempts, err := s.redis.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return attempts, err
}

// IsUserLocked 检查用户是否被锁定
func (s *LoginAttemptService) IsUserLocked(identifier string) (bool, int, error) {
	attempts, err := s.GetUserAttempts(identifier)
	if err != nil {
		return false, 0, err
	}
	
	if attempts >= 5 {
		ctx := context.Background()
		key := fmt.Sprintf("user_attempts:%s", identifier)
		ttl := s.redis.TTL(ctx, key).Val()
		return true, int(ttl.Minutes()), nil
	}
	
	return false, 0, nil
}