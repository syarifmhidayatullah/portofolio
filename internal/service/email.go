package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/syarifmhidayatullah/portfolio/config"
	"github.com/syarifmhidayatullah/portfolio/internal/model"
)

type EmailService interface {
	SendContactNotification(msg model.ContactMessage) error
}

type emailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) EmailService {
	return &emailService{cfg: cfg}
}

func (s *emailService) SendContactNotification(msg model.ContactMessage) error {
	subject := fmt.Sprintf("[Portfolio] New message from %s", msg.Name)
	body := fmt.Sprintf(
		"You have a new contact message:\n\nFrom: %s <%s>\nSubject: %s\n\nMessage:\n%s",
		msg.Name, msg.Email, msg.Subject, msg.Message,
	)

	to := s.cfg.NotifyEmail
	if to == "" {
		to = s.cfg.AdminEmail
	}

	switch s.cfg.EmailDriver {
	case "resend":
		return s.sendViaResend(to, subject, body)
	default:
		return s.sendViaSMTP(to, subject, body)
	}
}

func (s *emailService) sendViaSMTP(to, subject, body string) error {
	if s.cfg.SMTPUser == "" {
		log.Println("email: SMTP not configured, skipping")
		return nil
	}

	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)

	msg := strings.Join([]string{
		"From: " + s.cfg.SMTPFrom,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=utf-8",
		"",
		body,
	}, "\r\n")

	return smtp.SendMail(addr, auth, s.cfg.SMTPFrom, []string{to}, []byte(msg))
}

func (s *emailService) sendViaResend(to, subject, body string) error {
	if s.cfg.ResendAPIKey == "" {
		log.Println("email: Resend API key not configured, skipping")
		return nil
	}

	payload := map[string]interface{}{
		"from":    s.cfg.SMTPFrom,
		"to":      []string{to},
		"subject": subject,
		"text":    body,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.cfg.ResendAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("resend API error: %s", resp.Status)
	}
	return nil
}
