package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syarifmhidayatullah/portfolio/internal/service"
)

type DashboardHandler struct {
	postSvc    service.PostService
	projectSvc service.ProjectService
	messageSvc service.MessageService
}

func NewDashboardHandler(postSvc service.PostService, projectSvc service.ProjectService, messageSvc service.MessageService) *DashboardHandler {
	return &DashboardHandler{postSvc: postSvc, projectSvc: projectSvc, messageSvc: messageSvc}
}

func (h *DashboardHandler) Index(c *gin.Context) {
	postCount, _ := h.postSvc.Count()
	projectCount, _ := h.projectSvc.Count()
	unreadCount, _ := h.messageSvc.UnreadCount()
	recentPosts, _ := h.postSvc.Recent(5)
	recentMessages, _ := h.messageSvc.GetAll()

	// Limit messages to 5
	if len(recentMessages) > 5 {
		recentMessages = recentMessages[:5]
	}

	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"title":          "Dashboard",
		"postCount":      postCount,
		"projectCount":   projectCount,
		"unreadCount":    unreadCount,
		"recentPosts":    recentPosts,
		"recentMessages": recentMessages,
	})
}
