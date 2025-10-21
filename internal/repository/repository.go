package repository

import (
	"context"
)

// ShortURLRepository 定义了对ShortURL数据的所有操作契约
// 这是一个接口，这样我们以后可以轻松切换不同的数据库实现
type ShortURLRepository interface {
	// Create 创建新的短链接记录
	Create(ctx context.Context, shortURL *ShortURL) error

	// FindByShortCode 根据短码查找长链接
	FindByShortCode(ctx context.Context, shortCode string) (*ShortURL, error)

	// ExistsByShortCode 检查短码是否已存在
	ExistsByShortCode(ctx context.Context, shortCode string) (bool, error)

	// 未来可以在这里添加更多方法，比如：
	// Delete(ctx context.Context, shortCode string) error
	// FindByLongURL(ctx context.Context, longURL string) (*ShortURL, error)
}
