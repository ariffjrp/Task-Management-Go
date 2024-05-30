package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	Id        uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	UserId    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Firstname string         `json:"firstname"`
	Lastname  string         `json:"lastname"`
	Address   string         `json:"address"`
	Phone     string         `json:"phone"`
	Birthday  time.Time      `json:"birthday"`
	Gender    string         `json:"gender"`
	Avatar    string         `json:"avatar"`
	CreatedAt time.Time      `gorm:"not null" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"not null" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (account *Account) BeforeCreate(tx *gorm.DB) (err error) {
	account.Id = uuid.New()

	return nil
}
