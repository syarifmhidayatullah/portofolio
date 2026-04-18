package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Project) TableName() string { return "por_projects" }

type Project struct {
	ID          uuid.UUID      `gorm:"type:varchar(36);primaryKey"`
	Title       string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	TechStack   []string       `gorm:"serializer:json"`
	LiveURL     string
	GithubURL   string
	ImageURL    string
	Featured    bool           `gorm:"default:false"`
	OrderIndex  int            `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
