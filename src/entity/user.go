package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id           uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	Email        string         `gorm:"not null" json:"email"`
	Password     string         `gorm:"not null" json:"-"`
	Provider     string         `gorm:"not null" json:"provider"`
	CreatedAt    time.Time      `gorm:"not null" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"not null" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	Account      Account        `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"account"`
	RefreshToken RefreshToken   `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegister struct {
	Email     string    `json:"email" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	Firstname string    `json:"firstname" binding:"required"`
	Lastname  string    `json:"lastname" binding:"required"`
	Address   string    `json:"address" binding:"required"`
	Phone     string    `json:"phone" binding:"required"`
	Birthday  time.Time `json:"birthday" binding:"required"`
	Gender    string    `json:"gender" binding:"required"`
	Avatar    string    `json:"avatar" binding:"required"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.Id = uuid.New()

	return nil
}
