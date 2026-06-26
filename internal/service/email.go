package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		log.Println("[EMAIL] SMTP not configured, skipping")
		return nil
	}

	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)
	log.Printf("[EMAIL] sending via SMTP to=%s from=%s addr=%s user=%s", to, s.cfg.SMTPFrom, addr, s.cfg.SMTPUser)

	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	msg := strings.Join([]string{
		"From: " + s.cfg.SMTPFrom,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=utf-8",
		"",
		body,
	}, "\r\n")

	if err := smtp.SendMail(addr, auth, s.cfg.SMTPFrom, []string{to}, []byte(msg)); err != nil {
		log.Printf("[EMAIL] SMTP send failed: %v", err)
		return err
	}
	log.Printf("[EMAIL] sent successfully to %s", to)
	return nil
}

func (s *emailService) sendViaResend(to, subject, body string) error {
	if s.cfg.ResendAPIKey == "" {
		log.Println("[EMAIL] Resend API key not configured, skipping")
		return nil
	}

	log.Printf("[EMAIL] sending via Resend to=%s from=%s subject=%q", to, s.cfg.SMTPFrom, subject)

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
		log.Printf("[EMAIL] Resend request failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		log.Printf("[EMAIL] Resend API error: status=%s body=%s", resp.Status, string(respBody))
		return fmt.Errorf("resend API error: %s: %s", resp.Status, string(respBody))
	}

	log.Printf("[EMAIL] sent successfully via Resend to %s (response: %s)", to, string(respBody))
	return nil
}
