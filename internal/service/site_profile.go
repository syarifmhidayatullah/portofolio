package service

import (
	"github.com/syarifmhidayatullah/portfolio/internal/model"
	"github.com/syarifmhidayatullah/portfolio/internal/repository"
)

type SiteProfileService interface {
	Get() (*model.SiteProfile, error)
	Save(p *model.SiteProfile) error
}

type siteProfileService struct {
	repo repository.SiteProfileRepository
}

func NewSiteProfileService(repo repository.SiteProfileRepository) SiteProfileService {
	return &siteProfileService{repo: repo}
}

func (s *siteProfileService) Get() (*model.SiteProfile, error) {
	return s.repo.GetOrCreate()
}

func (s *siteProfileService) Save(p *model.SiteProfile) error {
	return s.repo.Save(p)
}
