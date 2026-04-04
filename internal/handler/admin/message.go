package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/syarifmhidayatullah/portfolio/internal/service"
)

type AdminMessageHandler struct {
	messageSvc service.MessageService
}

func NewMessageHandler(messageSvc service.MessageService) *AdminMessageHandler {
	return &AdminMessageHandler{messageSvc: messageSvc}
}

func (h *AdminMessageHandler) List(c *gin.Context) {
	messages, err := h.messageSvc.GetAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_messages.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "admin_messages.html", gin.H{
		"title":    "Messages",
		"messages": messages,
	})
}

func (h *AdminMessageHandler) MarkRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/messages")
		return
	}

	h.messageSvc.MarkRead(id)
	c.Redirect(http.StatusFound, "/admin/messages")
}

func (h *AdminMessageHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/messages")
		return
	}

	h.messageSvc.Delete(id)
	c.Redirect(http.StatusFound, "/admin/messages")
}
