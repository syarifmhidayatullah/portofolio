package admin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	input, imageURL, err := buildProjectInput(c, "")
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin_project_form.html", gin.H{
			"title": "New Project",
			"error": err.Error(),
			"input": input,
		})
		return
	}
	input.ImageURL = imageURL

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
		"title":        "Edit Project",
		"project":      project,
		"techStackRaw": strings.Join(project.TechStack, ", "),
	})
}

func (h *AdminProjectHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/projects")
		return
	}

	existing, err := h.projectSvc.GetByID(id)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/projects")
		return
	}

	input, imageURL, err := buildProjectInput(c, existing.ImageURL)
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin_project_form.html", gin.H{
			"title":   "Edit Project",
			"error":   err.Error(),
			"project": existing,
		})
		return
	}
	input.ImageURL = imageURL

	if _, err := h.projectSvc.Update(id, input); err != nil {
		c.HTML(http.StatusInternalServerError, "admin_project_form.html", gin.H{
			"title":   "Edit Project",
			"error":   err.Error(),
			"project": existing,
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

// buildProjectInput parses the multipart form. If a file is uploaded it saves it
// to web/static/uploads/ and returns its URL. Otherwise it keeps oldImageURL.
func buildProjectInput(c *gin.Context, oldImageURL string) (service.ProjectInput, string, error) {
	orderIndex, _ := strconv.Atoi(c.PostForm("order_index"))
	input := service.ProjectInput{
		Title:       strings.TrimSpace(c.PostForm("title")),
		Description: strings.TrimSpace(c.PostForm("description")),
		TechStack:   c.PostForm("tech_stack"),
		LiveURL:     strings.TrimSpace(c.PostForm("live_url")),
		GithubURL:   strings.TrimSpace(c.PostForm("github_url")),
		Featured:    c.PostForm("featured") == "on",
		OrderIndex:  orderIndex,
	}

	file, header, err := c.Request.FormFile("image_file")
	if err == nil {
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
		if !allowed[ext] {
			return input, oldImageURL, fmt.Errorf("format gambar tidak didukung (jpg, png, webp, gif)")
		}
		if header.Size > 5<<20 {
			return input, oldImageURL, fmt.Errorf("ukuran gambar maksimal 5MB")
		}

		filename := fmt.Sprintf("%d-%s%s", time.Now().UnixMilli(), uuid.New().String()[:8], ext)
		savePath := filepath.Join("web", "static", "uploads", filename)

		if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
			return input, oldImageURL, fmt.Errorf("gagal menyimpan gambar")
		}
		dst, err := os.Create(savePath)
		if err != nil {
			return input, oldImageURL, fmt.Errorf("gagal menyimpan gambar")
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			return input, oldImageURL, fmt.Errorf("gagal menyimpan gambar")
		}

		return input, "/static/uploads/" + filename, nil
	}

	// No file uploaded — use manual URL field or keep existing
	if urlField := strings.TrimSpace(c.PostForm("image_url")); urlField != "" {
		return input, urlField, nil
	}
	return input, oldImageURL, nil
}
