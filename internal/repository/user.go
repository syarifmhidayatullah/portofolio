package repository

import (
	"github.com/google/uuid"
	"github.com/syarifmhidayatullah/portfolio/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindByID(id uuid.UUID) (*model.User, error)
	Create(user *model.User) error
	Exists() bool
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Exists() bool {
	var count int64
	r.db.Model(&model.User{}).Count(&count)
	return count > 0
}
