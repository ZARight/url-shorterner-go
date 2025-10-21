package database

import (
	"fmt"
	"url-shortener/internal/repository"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate((&repository.ShortURL{}))
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	fmt.Println("✅ Database migration completed successfully!")
	return nil

}
