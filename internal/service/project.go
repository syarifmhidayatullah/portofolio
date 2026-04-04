package service

import (
	"strings"

	"github.com/google/uuid"
	"github.com/syarifmhidayatullah/portfolio/internal/model"
	"github.com/syarifmhidayatullah/portfolio/internal/repository"
)

type ProjectService interface {
	GetAll() ([]model.Project, error)
	GetFeatured() ([]model.Project, error)
	GetByID(id uuid.UUID) (*model.Project, error)
	Create(input ProjectInput) (*model.Project, error)
	Update(id uuid.UUID, input ProjectInput) (*model.Project, error)
	Delete(id uuid.UUID) error
	Count() (int64, error)
}

type ProjectInput struct {
	Title       string
	Description string
	TechStack   string // comma-separated
	LiveURL     string
	GithubURL   string
	ImageURL    string
	Featured    bool
	OrderIndex  int
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) GetAll() ([]model.Project, error) {
	return s.repo.FindAll()
}

func (s *projectService) GetFeatured() ([]model.Project, error) {
	return s.repo.FindFeatured()
}

func (s *projectService) GetByID(id uuid.UUID) (*model.Project, error) {
	return s.repo.FindByID(id)
}

func (s *projectService) Create(input ProjectInput) (*model.Project, error) {
	project := &model.Project{
		Title:       input.Title,
		Description: input.Description,
		TechStack:   parseTechStack(input.TechStack),
		LiveURL:     input.LiveURL,
		GithubURL:   input.GithubURL,
		ImageURL:    input.ImageURL,
		Featured:    input.Featured,
		OrderIndex:  input.OrderIndex,
	}
	return project, s.repo.Create(project)
}

func (s *projectService) Update(id uuid.UUID, input ProjectInput) (*model.Project, error) {
	project, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	project.Title = input.Title
	project.Description = input.Description
	project.TechStack = parseTechStack(input.TechStack)
	project.LiveURL = input.LiveURL
	project.GithubURL = input.GithubURL
	project.ImageURL = input.ImageURL
	project.Featured = input.Featured
	project.OrderIndex = input.OrderIndex

	return project, s.repo.Update(project)
}

func (s *projectService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *projectService) Count() (int64, error) {
	return s.repo.Count()
}

func parseTechStack(raw string) []string {
	var result []string
	for _, t := range strings.Split(raw, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}
