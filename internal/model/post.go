package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Post) TableName() string { return "por_posts" }

type Post struct {
	ID          uuid.UUID      `gorm:"type:varchar(36);primaryKey"`
	Title       string         `gorm:"type:varchar(255);not null"`
	Slug        string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	Excerpt     string
	Content     string         `gorm:"type:text;not null"`
	CoverImage  string
	Published   bool           `gorm:"default:false"`
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
