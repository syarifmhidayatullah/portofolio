package model

import "time"

func (SiteProfile) TableName() string { return "por_site_profile" }

type SiteProfile struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Bio1         string    `gorm:"type:text"`
	Bio2         string    `gorm:"type:text"`
	StatsYears   int       `gorm:"default:0"`
	StatsProjects int      `gorm:"default:0"`
	AvatarURL    string    `gorm:"type:text"`
	CVURL        string    `gorm:"type:text"`
	UpdatedAt    time.Time
}
