package utils

import (
	"log"
	"task_management/src/entity"

	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&entity.User{}, &entity.Account{}, &entity.RefreshToken{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
