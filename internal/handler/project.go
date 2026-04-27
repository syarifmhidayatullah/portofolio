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

	// Collect unique tech stack values for filter chips
	seen := map[string]bool{}
	var allTechs []string
	for _, p := range projects {
		for _, t := range p.TechStack {
			if !seen[t] {
				seen[t] = true
				allTechs = append(allTechs, t)
			}
		}
	}

	c.HTML(http.StatusOK, "projects.html", gin.H{
		"title":         "Projects",
		"activeNav":     "projects",
		"ogDescription": "A collection of things I've built — from production systems to side experiments.",
		"projects":      projects,
		"allTechs":      allTechs,
	})
}
