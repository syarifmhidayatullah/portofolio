package service

import (
	"github.com/google/uuid"
	"github.com/syarifmhidayatullah/portfolio/internal/model"
	"github.com/syarifmhidayatullah/portfolio/internal/repository"
)

type MessageService interface {
	GetAll() ([]model.ContactMessage, error)
	GetByID(id uuid.UUID) (*model.ContactMessage, error)
	Submit(input ContactInput) error
	MarkRead(id uuid.UUID) error
	Delete(id uuid.UUID) error
	UnreadCount() (int64, error)
}

type ContactInput struct {
	Name    string
	Email   string
	Subject string
	Message string
}

type messageService struct {
	repo     repository.MessageRepository
	emailSvc EmailService
}

func NewMessageService(repo repository.MessageRepository, emailSvc EmailService) MessageService {
	return &messageService{repo: repo, emailSvc: emailSvc}
}

func (s *messageService) GetAll() ([]model.ContactMessage, error) {
	return s.repo.FindAll()
}

func (s *messageService) GetByID(id uuid.UUID) (*model.ContactMessage, error) {
	return s.repo.FindByID(id)
}

func (s *messageService) Submit(input ContactInput) error {
	msg := &model.ContactMessage{
		Name:    input.Name,
		Email:   input.Email,
		Subject: input.Subject,
		Message: input.Message,
	}

	if err := s.repo.Create(msg); err != nil {
		return err
	}

	// Send email notification (best-effort, don't fail the request)
	go s.emailSvc.SendContactNotification(*msg)

	return nil
}

func (s *messageService) MarkRead(id uuid.UUID) error {
	return s.repo.MarkRead(id)
}

func (s *messageService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *messageService) UnreadCount() (int64, error) {
	return s.repo.UnreadCount()
}
