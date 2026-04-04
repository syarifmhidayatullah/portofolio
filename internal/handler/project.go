package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syarifmhidayatullah/portfolio/internal/service"
)

type ProjectHandler struct {
	projectSvc service.ProjectService
}

func NewProjectHandler(projectSvc service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectSvc: projectSvc}
}

func (h *ProjectHandler) List(c *gin.Context) {
	projects, err := h.projectSvc.GetAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "projects.html", gin.H{
		"title":    "Projects",
		"projects": projects,
	})
}
