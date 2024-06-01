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

func AddOAuth2IdToUsers(db *gorm.DB) error {
	if !db.Migrator().HasColumn(&entity.User{}, "o_auth2_id") {
		err := db.Migrator().AddColumn(&entity.User{}, "o_auth2_id")
		if err != nil {
			return err
		}
	}
	return nil
}
