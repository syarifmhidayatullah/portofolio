package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syarifmhidayatullah/portfolio/internal/service"
)

type HomeHandler struct {
	postSvc       service.PostService
	projectSvc    service.ProjectService
	profileSvc    service.SiteProfileService
	experienceSvc service.ExperienceService
}

func NewHomeHandler(postSvc service.PostService, projectSvc service.ProjectService, profileSvc service.SiteProfileService, experienceSvc service.ExperienceService) *HomeHandler {
	return &HomeHandler{postSvc: postSvc, projectSvc: projectSvc, profileSvc: profileSvc, experienceSvc: experienceSvc}
}

func (h *HomeHandler) Index(c *gin.Context) {
	posts, _ := h.postSvc.Recent(3)
	featured, _ := h.projectSvc.GetFeatured()
	profile, _ := h.profileSvc.Get()
	experiences, _ := h.experienceSvc.GetAll()

	c.HTML(http.StatusOK, "home.html", gin.H{
		"title":       "Home",
		"posts":       posts,
		"projects":    featured,
		"profile":     profile,
		"experiences": experiences,
	})
}
