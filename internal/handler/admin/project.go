package admin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/syarifmhidayatullah/portfolio/internal/service"
)

type AdminProjectHandler struct {
	projectSvc service.ProjectService
}

func NewProjectHandler(projectSvc service.ProjectService) *AdminProjectHandler {
	return &AdminProjectHandler{projectSvc: projectSvc}
}

func (h *AdminProjectHandler) List(c *gin.Context) {
	projects, err := h.projectSvc.GetAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_projects.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "admin_projects.html", gin.H{
		"title":    "Projects",
		"projects": projects,
	})
}

func (h *AdminProjectHandler) New(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_project_form.html", gin.H{
		"title": "New Project",
	})
}

func (h *AdminProjectHandler) Create(c *gin.Context) {
	input := buildProjectInput(c)

	if input.Title == "" {
		c.HTML(http.StatusBadRequest, "admin_project_form.html", gin.H{
			"title": "New Project",
			"error": "Title is required",
			"input": input,
		})
		return
	}

	if _, err := h.projectSvc.Create(input); err != nil {
		c.HTML(http.StatusInternalServerError, "admin_project_form.html", gin.H{
			"title": "New Project",
			"error": err.Error(),
			"input": input,
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/projects")
}

func (h *AdminProjectHandler) Edit(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/projects")
		return
	}

	project, err := h.projectSvc.GetByID(id)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/projects")
		return
	}

	c.HTML(http.StatusOK, "admin_project_form.html", gin.H{
		"title":          "Edit Project",
		"project":        project,
		"techStackRaw":   strings.Join(project.TechStack, ", "),
	})
}

func (h *AdminProjectHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/projects")
		return
	}

	input := buildProjectInput(c)

	if _, err := h.projectSvc.Update(id, input); err != nil {
		project, _ := h.projectSvc.GetByID(id)
		c.HTML(http.StatusInternalServerError, "admin_project_form.html", gin.H{
			"title":   "Edit Project",
			"error":   err.Error(),
			"project": project,
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/projects")
}

func (h *AdminProjectHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/projects")
		return
	}

	h.projectSvc.Delete(id)
	c.Redirect(http.StatusFound, "/admin/projects")
}

func buildProjectInput(c *gin.Context) service.ProjectInput {
	orderIndex, _ := strconv.Atoi(c.PostForm("order_index"))
	return service.ProjectInput{
		Title:       strings.TrimSpace(c.PostForm("title")),
		Description: strings.TrimSpace(c.PostForm("description")),
		TechStack:   c.PostForm("tech_stack"),
		LiveURL:     strings.TrimSpace(c.PostForm("live_url")),
		GithubURL:   strings.TrimSpace(c.PostForm("github_url")),
		ImageURL:    strings.TrimSpace(c.PostForm("image_url")),
		Featured:    c.PostForm("featured") == "on",
		OrderIndex:  orderIndex,
	}
}
