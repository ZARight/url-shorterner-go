package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// mysqlShortURLRepository 是 ShortURLRepository 的 MySQL 实现
type mysqlShortURLRepository struct {
	db *gorm.DB
}

// NewMysqlShortURLRepository 创建MySQL仓库实例
func NewMysqlShortURLRepository(db *gorm.DB) ShortURLRepository {
	return &mysqlShortURLRepository{db: db}
}

// Create 实现创建短链接的方法
func (r *mysqlShortURLRepository) Create(ctx context.Context, shortURL *ShortURL) error {
	result := r.db.WithContext(ctx).Create(shortURL)
	if result.Error != nil {
		return fmt.Errorf("failed to create short URL: %w", result.Error)
	}
	return nil
}

// FindByShortCode 实现根据短码查找的方法
func (r *mysqlShortURLRepository) FindByShortCode(ctx context.Context, shortCode string) (*ShortURL, error) {
	var shortURL ShortURL
	result := r.db.WithContext(ctx).Where("short_code = ?", shortCode).First(&shortURL)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 没有找到记录，返回nil
		}
		return nil, fmt.Errorf("failed to find short URL: %w", result.Error)
	}

	return &shortURL, nil
}

// ExistsByShortCode 检查短码是否存在
func (r *mysqlShortURLRepository) ExistsByShortCode(ctx context.Context, shortCode string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&ShortURL{}).Where("short_code = ?", shortCode).Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check short code existence: %w", result.Error)
	}

	return count > 0, nil
}
