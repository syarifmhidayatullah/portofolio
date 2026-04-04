package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContactMessage struct {
	ID        uuid.UUID      `gorm:"type:varchar(36);primaryKey"`
	Name      string         `gorm:"type:varchar(255);not null"`
	Email     string         `gorm:"type:varchar(255);not null"`
	Subject   string
	Message   string         `gorm:"type:text;not null"`
	IsRead    bool           `gorm:"default:false"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *ContactMessage) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
