package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"url-shortener/internal/repository"
)

// ShortenService 处理短链接的业务逻辑
type ShortenService struct {
	repo  repository.ShortURLRepository
	cache *CacheService
}

// NewShortenService 创建ShortenService实例
func NewShortenService(repo repository.ShortURLRepository, cache *CacheService) *ShortenService {
	return &ShortenService{
		repo:  repo,
		cache: cache,
	}
}

// CreateShortURLRequest 创建短链接的请求结构
type CreateShortURLRequest struct {
	LongURL string `json:"long_url" binding:"required,url"`
}

// CreateShortURLResponse 创建短链接的响应结构
type CreateShortURLResponse struct {
	ShortURL  string `json:"short_url"`
	ShortCode string `json:"short_code"`
	LongURL   string `json:"long_url"`
}

// CreateShortURL 创建短链接的核心业务逻辑
func (s *ShortenService) CreateShortURL(ctx context.Context, req *CreateShortURLRequest) (*CreateShortURLResponse, error) {
	// 1. 验证长链接
	if strings.TrimSpace(req.LongURL) == "" {
		return nil, fmt.Errorf("long URL cannot be empty")
	}

	// 2. 生成短码（使用MD5哈希 + 截取前6位作为简单实现）
	shortCode := s.generateShortCode(req.LongURL)

	// 3. 检查短码是否已存在，如果存在则重新生成（处理哈希冲突）
	for i := 0; i < 3; i++ { // 最多尝试3次
		exists, err := s.repo.ExistsByShortCode(ctx, shortCode)
		if err != nil {
			return nil, fmt.Errorf("failed to check short code existence: %w", err)
		}

		if !exists {
			break // 短码可用，跳出循环
		}

		// 短码已存在，重新生成（添加随机后缀）
		shortCode = s.generateShortCode(req.LongURL + string(rune(i)))
	}

	// 4. 创建短链接记录
	shortURL := &repository.ShortURL{
		ShortCode: shortCode,
		LongURL:   req.LongURL,
	}

	if err := s.repo.Create(ctx, shortURL); err != nil {
		return nil, fmt.Errorf("failed to create short URL: %w", err)
	}

	// 在创建数据库记录后，也设置缓存
	if err := s.cache.SetShortURL(ctx, shortURL); err != nil {
		fmt.Printf("⚠️ Failed to set cache after creation: %v\n", err)
		// 不返回错误，因为主要业务逻辑已经完成
	}

	// 5. 返回响应
	response := &CreateShortURLResponse{
		ShortURL:  fmt.Sprintf("http://localhost:8080/%s", shortCode), // 暂时硬编码，后续可以从配置读取
		ShortCode: shortCode,
		LongURL:   req.LongURL,
	}

	return response, nil
}

// generateShortCode 生成短码的辅助函数
// 目前使用简单的MD5哈希，后续可以改进为更分布式的方案
func (s *ShortenService) generateShortCode(longURL string) string {
	hash := md5.Sum([]byte(longURL))
	hashStr := hex.EncodeToString(hash[:])

	// 取前6位作为短码
	return hashStr[:6]
}

// GetLongURL 根据短码获取长链接
func (s *ShortenService) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	if strings.TrimSpace(shortCode) == "" {
		return "", fmt.Errorf("short code cannot be empty")
	}

	// 1. 先尝试从缓存获取
	cachedShortURL, err := s.cache.GetShortURL(ctx, shortCode)
	if err != nil {
		return "", fmt.Errorf("failed to get from cache: %w", err)
	}

	// 2. 如果缓存命中，直接返回
	if cachedShortURL != nil {
		fmt.Printf("✅ Cache HIT for short code: %s\n", shortCode)
		return cachedShortURL.LongURL, nil
	}

	fmt.Printf("❌ Cache MISS for short code: %s\n", shortCode)

	// 3. 缓存未命中，查询数据库
	shortURL, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", fmt.Errorf("failed to find short URL: %w", err)
	}

	if shortURL == nil {
		return "", fmt.Errorf("short URL not found for code: %s", shortCode)
	}

	// 4. 将查询结果存入缓存
	if err := s.cache.SetShortURL(ctx, shortURL); err != nil {
		fmt.Printf("⚠️ Failed to set cache: %v\n", err)
		// 不返回错误，因为主要业务逻辑已经完成
	}

	return shortURL.LongURL, nil
}
