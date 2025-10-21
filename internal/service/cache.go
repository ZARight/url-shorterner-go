package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"url-shortener/internal/repository"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	redisClient *redis.Client
}

func NewCacheService(redisClient *redis.Client) *CacheService {
	return &CacheService{
		redisClient: redisClient,
	}
}

// GetShortURL 从缓存获取短链接信息
func (c *CacheService) GetShortURL(ctx context.Context, shortCode string) (*repository.ShortURL, error) {
	key := fmt.Sprintf("short_url:%s", shortCode)

	val, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	var shortURL repository.ShortURL
	if err := json.Unmarshal([]byte(val), &shortURL); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	return &shortURL, nil
}

// SetShortURL 将短链接信息存入缓存
func (c *CacheService) SetShortURL(ctx context.Context, shortURL *repository.ShortURL) error {
	key := fmt.Sprintf("short_url:%s", shortURL.ShortCode)

	// 序列化数据
	data, err := json.Marshal(shortURL)
	if err != nil {
		return fmt.Errorf("failed to marshal short URL: %w", err)
	}

	// 设置缓存，过期时间为24小时
	err = c.redisClient.Set(ctx, key, data, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// DeleteShortURL 从缓存删除短链接信息
func (c *CacheService) DeleteShortURL(ctx context.Context, shortCode string) error {
	key := fmt.Sprintf("short_url:%s", shortCode)
	return c.redisClient.Del(ctx, key).Err()
}
