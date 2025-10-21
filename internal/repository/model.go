package repository

import (
	"time"
)

// ShortURL 表示数据库中的短链接模型
// 使用GORM标签来定义数据库字段属性
type ShortURL struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ShortCode string    `gorm:"size:10;not null;uniqueIndex" json:"short_code"` // 唯一索引加速查询
	LongURL   string    `gorm:"type:text;not null" json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定数据库中对应的表名
func (ShortURL) TableName() string {
	return "short_urls"
}
