package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/syarifmhidayatullah/portfolio/internal/service"
)

type ContactHandler struct {
	messageSvc service.MessageService
}

func NewContactHandler(messageSvc service.MessageService) *ContactHandler {
	return &ContactHandler{messageSvc: messageSvc}
}

func (h *ContactHandler) Submit(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.TrimSpace(c.PostForm("email"))
	subject := strings.TrimSpace(c.PostForm("subject"))
	message := strings.TrimSpace(c.PostForm("message"))

	if name == "" || email == "" || message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, email, and message are required"})
		return
	}

	err := h.messageSvc.Submit(service.ContactInput{
		Name:    name,
		Email:   email,
		Subject: subject,
		Message: message,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message. Please try again."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your message has been sent! I'll get back to you soon."})
}
