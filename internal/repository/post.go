package repository

import (
	"github.com/google/uuid"
	"github.com/syarifmhidayatullah/portfolio/internal/model"
	"gorm.io/gorm"
)

type PostRepository interface {
	FindAll(onlyPublished bool) ([]model.Post, error)
	FindBySlug(slug string) (*model.Post, error)
	FindByID(id uuid.UUID) (*model.Post, error)
	Create(post *model.Post) error
	Update(post *model.Post) error
	Delete(id uuid.UUID) error
	Count() (int64, error)
	Recent(limit int) ([]model.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) FindAll(onlyPublished bool) ([]model.Post, error) {
	var posts []model.Post
	q := r.db.Order("created_at desc")
	if onlyPublished {
		q = q.Where("published = true")
	}
	err := q.Find(&posts).Error
	return posts, err
}

func (r *postRepository) FindBySlug(slug string) (*model.Post, error) {
	var post model.Post
	err := r.db.Where("slug = ? AND published = true", slug).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) FindByID(id uuid.UUID) (*model.Post, error) {
	var post model.Post
	err := r.db.First(&post, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Post{}, "id = ?", id).Error
}

func (r *postRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Post{}).Count(&count).Error
	return count, err
}

func (r *postRepository) Recent(limit int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Order("created_at desc").Limit(limit).Find(&posts).Error
	return posts, err
}
