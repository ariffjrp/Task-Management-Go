package utils

import (
	"task_management/src/entity"
	"gorm.io/gorm"
	"log"
)

func MigrateDB(db *gorm.DB) {
	// Migrate the schema
	err := db.AutoMigrate(&entity.User{}, &entity.Account{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}