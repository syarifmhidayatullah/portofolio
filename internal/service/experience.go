package service

import (
	"errors"
	"strings"

	"github.com/syarifmhidayatullah/portfolio/internal/model"
	"github.com/syarifmhidayatullah/portfolio/internal/repository"
)

type ExperienceInput struct {
	Role        string
	Company     string
	PeriodStart string
	PeriodEnd   string
	Description string
	TechStack   string // comma-separated raw input
	IsCurrent   bool
	SortOrder   int
}

type ExperienceService interface {
	GetAll() ([]model.Experience, error)
	GetByID(id uint) (*model.Experience, error)
	Create(input ExperienceInput) (*model.Experience, error)
	Update(id uint, input ExperienceInput) (*model.Experience, error)
	Delete(id uint) error
}

type experienceService struct {
	repo repository.ExperienceRepository
}

func NewExperienceService(repo repository.ExperienceRepository) ExperienceService {
	return &experienceService{repo: repo}
}

func (s *experienceService) GetAll() ([]model.Experience, error) {
	return s.repo.FindAll()
}

func (s *experienceService) GetByID(id uint) (*model.Experience, error) {
	return s.repo.FindByID(id)
}

func (s *experienceService) Create(input ExperienceInput) (*model.Experience, error) {
	if input.Role == "" || input.Company == "" {
		return nil, errors.New("role and company are required")
	}
	e := &model.Experience{
		Role:        input.Role,
		Company:     input.Company,
		PeriodStart: input.PeriodStart,
		PeriodEnd:   input.PeriodEnd,
		Description: input.Description,
		TechStack:   splitTechStack(input.TechStack),
		IsCurrent:   input.IsCurrent,
		SortOrder:   input.SortOrder,
	}
	return e, s.repo.Create(e)
}

func (s *experienceService) Update(id uint, input ExperienceInput) (*model.Experience, error) {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("experience not found")
	}
	e.Role = input.Role
	e.Company = input.Company
	e.PeriodStart = input.PeriodStart
	e.PeriodEnd = input.PeriodEnd
	e.Description = input.Description
	e.TechStack = parseTechStack(input.TechStack)
	e.IsCurrent = input.IsCurrent
	e.SortOrder = input.SortOrder
	return e, s.repo.Update(e)
}

func (s *experienceService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func splitTechStack(raw string) []string {
	var result []string
	for _, t := range strings.Split(raw, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}
