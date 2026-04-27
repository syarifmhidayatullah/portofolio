package model

import "time"

func (Experience) TableName() string { return "por_experiences" }

type Experience struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Role        string    `gorm:"type:varchar(150);not null"`
	Company     string    `gorm:"type:varchar(150);not null"`
	PeriodStart string    `gorm:"type:varchar(20)"`
	PeriodEnd   string    `gorm:"type:varchar(20)"` // empty = "Present"
	Description string    `gorm:"type:text"`
	TechStack   []string  `gorm:"serializer:json"`
	IsCurrent   bool      `gorm:"default:false"`
	SortOrder   int       `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
