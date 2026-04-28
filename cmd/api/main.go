package main

import (
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/syarifmhidayatullah/portfolio/config"
	"github.com/syarifmhidayatullah/portfolio/internal/handler"
	adminHandler "github.com/syarifmhidayatullah/portfolio/internal/handler/admin"
	"github.com/syarifmhidayatullah/portfolio/internal/middleware"
	"github.com/syarifmhidayatullah/portfolio/internal/model"
	"github.com/syarifmhidayatullah/portfolio/internal/repository"
	"github.com/syarifmhidayatullah/portfolio/internal/service"
)

func main() {
	cfg := config.Load()

	// Database
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Project{},
		&model.ContactMessage{},
		&model.SiteProfile{},
		&model.Experience{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// Seed admin user
	seedAdmin(db, cfg)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	profileRepo := repository.NewSiteProfileRepository(db)
	experienceRepo := repository.NewExperienceRepository(db)

	// Services
	emailSvc := service.NewEmailService(cfg)
	postSvc := service.NewPostService(postRepo)
	projectSvc := service.NewProjectService(projectRepo)
	messageSvc := service.NewMessageService(messageRepo, emailSvc)
	profileSvc := service.NewSiteProfileService(profileRepo)
	experienceSvc := service.NewExperienceService(experienceRepo)

	// Handlers
	homeH := handler.NewHomeHandler(postSvc, projectSvc, profileSvc, experienceSvc)
	seoH := handler.NewSEOHandler(postSvc, projectSvc, cfg.AppURL)
	blogH := handler.NewBlogHandler(postSvc)
	projectH := handler.NewProjectHandler(projectSvc)
	contactH := handler.NewContactHandler(messageSvc)

	adminAuthH := adminHandler.NewAuthHandler(userRepo, cfg)
	adminDashH := adminHandler.NewDashboardHandler(postSvc, projectSvc, messageSvc)
	adminPostH := adminHandler.NewPostHandler(postSvc)
	adminProjectH := adminHandler.NewProjectHandler(projectSvc)
	adminMessageH := adminHandler.NewMessageHandler(messageSvc)
	adminProfileH := adminHandler.NewProfileHandler(profileSvc, experienceSvc)

	// Router
	r := gin.Default()

	// Templates
	r.SetHTMLTemplate(loadTemplates())
	r.Static("/static", "./web/static")

	// Session
	store := cookie.NewStore([]byte(cfg.SessionSecret))
	r.Use(sessions.Sessions("portfolio_sess", store))

	// Public routes
	r.GET("/", homeH.Index)
	r.GET("/blog", blogH.List)
	r.GET("/blog/:slug", blogH.Detail)
	r.GET("/projects", projectH.List)
	r.POST("/contact", contactH.Submit)
	r.GET("/robots.txt", seoH.Robots)
	r.GET("/sitemap.xml", seoH.Sitemap)

	// Admin auth (public)
	r.GET("/admin/login", adminAuthH.LoginPage)
	r.POST("/admin/login", adminAuthH.LoginSubmit)
	r.POST("/admin/logout", adminAuthH.Logout)

	// Admin protected routes
	adm := r.Group("/admin", middleware.AuthRequired())
	adm.GET("", adminDashH.Index)

	adm.GET("/posts", adminPostH.List)
	adm.GET("/posts/new", adminPostH.New)
	adm.POST("/posts", adminPostH.Create)
	adm.GET("/posts/:id/edit", adminPostH.Edit)
	adm.POST("/posts/:id", adminPostH.Update)
	adm.POST("/posts/:id/delete", adminPostH.Delete)
	adm.POST("/posts/:id/toggle-publish", adminPostH.TogglePublish)

	adm.GET("/projects", adminProjectH.List)
	adm.GET("/projects/new", adminProjectH.New)
	adm.POST("/projects", adminProjectH.Create)
	adm.GET("/projects/:id/edit", adminProjectH.Edit)
	adm.POST("/projects/:id", adminProjectH.Update)
	adm.POST("/projects/:id/delete", adminProjectH.Delete)

	adm.GET("/messages", adminMessageH.List)
	adm.POST("/messages/:id/read", adminMessageH.MarkRead)
	adm.POST("/messages/:id/delete", adminMessageH.Delete)

	adm.GET("/profile", adminProfileH.Page)
	adm.POST("/profile", adminProfileH.SaveProfile)
	adm.POST("/profile/experience", adminProfileH.CreateExperience)
	adm.GET("/profile/experience/:id/edit", adminProfileH.EditExperiencePage)
	adm.POST("/profile/experience/:id", adminProfileH.UpdateExperience)
	adm.POST("/profile/experience/:id/delete", adminProfileH.DeleteExperience)

	log.Printf("Server starting on %s", cfg.Port)
	if err := r.Run(cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func loadTemplates() *template.Template {
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"currentYear": func() int {
			return time.Now().Year()
		},
		"readingTime": func(content string) int {
			words := len(strings.Fields(content))
			mins := words / 200
			if mins < 1 {
				return 1
			}
			return mins
		},
		"joinStrings": func(ss []string, sep string) string {
			result := ""
			for i, s := range ss {
				if i > 0 {
					result += sep
				}
				result += s
			}
			return result
		},
	}

	tmpl := template.Must(
		template.New("").Funcs(funcMap).ParseGlob("web/templates/partials/*.html"),
	)
	template.Must(tmpl.ParseGlob("web/templates/*.html"))
	template.Must(tmpl.ParseGlob("web/templates/admin/*.html"))

	return tmpl
}

func seedAdmin(db *gorm.DB, cfg *config.Config) {
	userRepo := repository.NewUserRepository(db)
	if userRepo.Exists() {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash admin password: %v", err)
		return
	}

	user := &model.User{
		Email:        cfg.AdminEmail,
		PasswordHash: string(hash),
	}

	if err := userRepo.Create(user); err != nil {
		log.Printf("failed to seed admin: %v", err)
		return
	}

	log.Printf("Admin user created: %s", cfg.AdminEmail)
}
