package admin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/syarifmhidayatullah/portfolio/internal/service"
)

type AdminProfileHandler struct {
	profileSvc    service.SiteProfileService
	experienceSvc service.ExperienceService
}

func NewProfileHandler(profileSvc service.SiteProfileService, experienceSvc service.ExperienceService) *AdminProfileHandler {
	return &AdminProfileHandler{profileSvc: profileSvc, experienceSvc: experienceSvc}
}

func (h *AdminProfileHandler) Page(c *gin.Context) {
	profile, _ := h.profileSvc.Get()
	experiences, _ := h.experienceSvc.GetAll()
	c.HTML(http.StatusOK, "admin_profile.html", gin.H{
		"title":       "Profile",
		"profile":     profile,
		"experiences": experiences,
	})
}

func (h *AdminProfileHandler) SaveProfile(c *gin.Context) {
	profile, err := h.profileSvc.Get()
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/profile")
		return
	}

	// Avatar upload → always save as avatar.jpg
	if avatar, header, err := c.Request.FormFile("avatar_file"); err == nil {
		defer avatar.Close()
		ext := strings.ToLower(filepath.Ext(header.Filename))
		allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
		if !allowed[ext] {
			renderProfilePage(c, h, "Format foto tidak didukung (jpg, png, webp)", "")
			return
		}
		if header.Size > 5<<20 {
			renderProfilePage(c, h, "Ukuran foto maksimal 5MB", "")
			return
		}
		savePath := filepath.Join("web", "static", "uploads", "avatar"+ext)
		if ext != ".jpg" && ext != ".jpeg" {
			savePath = filepath.Join("web", "static", "uploads", "avatar"+ext)
		} else {
			savePath = filepath.Join("web", "static", "uploads", "avatar.jpg")
		}
		if err := saveFile(avatar, savePath); err != nil {
			renderProfilePage(c, h, "Gagal menyimpan foto: "+err.Error(), "")
			return
		}
		profile.AvatarURL = "/static/uploads/avatar" + ext
		if ext == ".jpg" || ext == ".jpeg" {
			profile.AvatarURL = "/static/uploads/avatar.jpg"
		}
	}

	// CV upload → always save as cv.pdf
	if cv, header, err := c.Request.FormFile("cv_file"); err == nil {
		defer cv.Close()
		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext != ".pdf" {
			renderProfilePage(c, h, "CV harus berformat PDF", "")
			return
		}
		if header.Size > 10<<20 {
			renderProfilePage(c, h, "Ukuran CV maksimal 10MB", "")
			return
		}
		savePath := filepath.Join("web", "static", "uploads", "cv.pdf")
		if err := saveFile(cv, savePath); err != nil {
			renderProfilePage(c, h, "Gagal menyimpan CV: "+err.Error(), "")
			return
		}
		profile.CVURL = "/static/uploads/cv.pdf"
	}

	// Text fields
	profile.Bio1 = strings.TrimSpace(c.PostForm("bio1"))
	profile.Bio2 = strings.TrimSpace(c.PostForm("bio2"))
	profile.StatsYears, _ = strconv.Atoi(c.PostForm("stats_years"))
	profile.StatsProjects, _ = strconv.Atoi(c.PostForm("stats_projects"))

	if err := h.profileSvc.Save(profile); err != nil {
		renderProfilePage(c, h, "Gagal menyimpan: "+err.Error(), "")
		return
	}
	renderProfilePage(c, h, "", "Profile berhasil disimpan!")
}

func (h *AdminProfileHandler) CreateExperience(c *gin.Context) {
	input := buildExperienceInput(c)
	if _, err := h.experienceSvc.Create(input); err != nil {
		renderProfilePage(c, h, err.Error(), "")
		return
	}
	c.Redirect(http.StatusFound, "/admin/profile")
}

func (h *AdminProfileHandler) EditExperiencePage(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/profile")
		return
	}
	exp, err := h.experienceSvc.GetByID(id)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/profile")
		return
	}
	profile, _ := h.profileSvc.Get()
	experiences, _ := h.experienceSvc.GetAll()
	c.HTML(http.StatusOK, "admin_profile.html", gin.H{
		"title":          "Profile",
		"profile":        profile,
		"experiences":    experiences,
		"editExperience": exp,
		"techStackRaw":   strings.Join(exp.TechStack, ", "),
	})
}

func (h *AdminProfileHandler) UpdateExperience(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/profile")
		return
	}
	input := buildExperienceInput(c)
	if _, err := h.experienceSvc.Update(id, input); err != nil {
		renderProfilePage(c, h, err.Error(), "")
		return
	}
	c.Redirect(http.StatusFound, "/admin/profile")
}

func (h *AdminProfileHandler) DeleteExperience(c *gin.Context) {
	id, _ := parseUintParam(c, "id")
	h.experienceSvc.Delete(id)
	c.Redirect(http.StatusFound, "/admin/profile")
}

func buildExperienceInput(c *gin.Context) service.ExperienceInput {
	sortOrder, _ := strconv.Atoi(c.PostForm("sort_order"))
	return service.ExperienceInput{
		Role:        strings.TrimSpace(c.PostForm("role")),
		Company:     strings.TrimSpace(c.PostForm("company")),
		PeriodStart: strings.TrimSpace(c.PostForm("period_start")),
		PeriodEnd:   strings.TrimSpace(c.PostForm("period_end")),
		Description: strings.TrimSpace(c.PostForm("description")),
		TechStack:   c.PostForm("tech_stack"),
		IsCurrent:   c.PostForm("is_current") == "on",
		SortOrder:   sortOrder,
	}
}

func renderProfilePage(c *gin.Context, h *AdminProfileHandler, errMsg, successMsg string) {
	profile, _ := h.profileSvc.Get()
	experiences, _ := h.experienceSvc.GetAll()
	c.HTML(http.StatusOK, "admin_profile.html", gin.H{
		"title":       "Profile",
		"profile":     profile,
		"experiences": experiences,
		"error":       errMsg,
		"success":     successMsg,
	})
}

func saveFile(src io.Reader, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, src)
	return err
}

func parseUintParam(c *gin.Context, param string) (uint, error) {
	v, err := strconv.ParseUint(c.Param(param), 10, 64)
	return uint(v), err
}
