package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	Id         uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	UserId     uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Token      string         `gorm:"not null" json:"token"`
	ExpiryDate time.Time      `gorm:"not null" json:"expires_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (refreshToken *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	refreshToken.Id = uuid.New()

	return nil
}
