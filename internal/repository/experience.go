package repository

import (
	"github.com/syarifmhidayatullah/portfolio/internal/model"
	"gorm.io/gorm"
)

type ExperienceRepository interface {
	FindAll() ([]model.Experience, error)
	FindByID(id uint) (*model.Experience, error)
	Create(e *model.Experience) error
	Update(e *model.Experience) error
	Delete(id uint) error
}

type experienceRepository struct{ db *gorm.DB }

func NewExperienceRepository(db *gorm.DB) ExperienceRepository {
	return &experienceRepository{db: db}
}

func (r *experienceRepository) FindAll() ([]model.Experience, error) {
	var list []model.Experience
	err := r.db.Order("sort_order asc, created_at asc").Find(&list).Error
	return list, err
}

func (r *experienceRepository) FindByID(id uint) (*model.Experience, error) {
	var e model.Experience
	err := r.db.First(&e, id).Error
	return &e, err
}

func (r *experienceRepository) Create(e *model.Experience) error {
	return r.db.Create(e).Error
}

func (r *experienceRepository) Update(e *model.Experience) error {
	return r.db.Save(e).Error
}

func (r *experienceRepository) Delete(id uint) error {
	return r.db.Delete(&model.Experience{}, id).Error
}
