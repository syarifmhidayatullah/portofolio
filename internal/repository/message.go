package repository

import (
	"github.com/google/uuid"
	"github.com/syarifmhidayatullah/portfolio/internal/model"
	"gorm.io/gorm"
)

type MessageRepository interface {
	FindAll() ([]model.ContactMessage, error)
	FindByID(id uuid.UUID) (*model.ContactMessage, error)
	Create(msg *model.ContactMessage) error
	MarkRead(id uuid.UUID) error
	Delete(id uuid.UUID) error
	UnreadCount() (int64, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) FindAll() ([]model.ContactMessage, error) {
	var messages []model.ContactMessage
	err := r.db.Order("created_at desc").Find(&messages).Error
	return messages, err
}

func (r *messageRepository) FindByID(id uuid.UUID) (*model.ContactMessage, error) {
	var msg model.ContactMessage
	err := r.db.First(&msg, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *messageRepository) Create(msg *model.ContactMessage) error {
	return r.db.Create(msg).Error
}

func (r *messageRepository) MarkRead(id uuid.UUID) error {
	return r.db.Model(&model.ContactMessage{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *messageRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.ContactMessage{}, "id = ?", id).Error
}

func (r *messageRepository) UnreadCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.ContactMessage{}).Where("is_read = false").Count(&count).Error
	return count, err
}
