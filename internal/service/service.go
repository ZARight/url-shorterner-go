package service

import (
	"url-shortener/internal/repository"

	"github.com/redis/go-redis/v9"
)

// Service 聚合所有业务服务
type Service struct {
	Shorten *ShortenService
	Cache   *CacheService
}

// NewService 创建所有服务的实例
func NewService(repo repository.ShortURLRepository, redisClient *redis.Client) *Service {
	cacheService := NewCacheService(redisClient)

	return &Service{
		Shorten: NewShortenService(repo, cacheService),
		Cache:   cacheService,
	}
}
