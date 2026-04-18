package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (User) TableName() string { return "por_users" }

type User struct {
	ID           uuid.UUID      `gorm:"type:varchar(36);primaryKey"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string         `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
