package repository

import (
	"github.com/syarifmhidayatullah/portfolio/internal/model"
	"gorm.io/gorm"
)

type SiteProfileRepository interface {
	GetOrCreate() (*model.SiteProfile, error)
	Save(p *model.SiteProfile) error
}

type siteProfileRepository struct{ db *gorm.DB }

func NewSiteProfileRepository(db *gorm.DB) SiteProfileRepository {
	return &siteProfileRepository{db: db}
}

func (r *siteProfileRepository) GetOrCreate() (*model.SiteProfile, error) {
	var p model.SiteProfile
	err := r.db.First(&p).Error
	if err == gorm.ErrRecordNotFound {
		p = model.SiteProfile{}
		if err := r.db.Create(&p).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *siteProfileRepository) Save(p *model.SiteProfile) error {
	return r.db.Save(p).Error
}
