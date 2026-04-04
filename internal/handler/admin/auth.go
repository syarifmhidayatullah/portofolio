package admin

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/syarifmhidayatullah/portfolio/config"
	"github.com/syarifmhidayatullah/portfolio/internal/middleware"
	"github.com/syarifmhidayatullah/portfolio/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthHandler(userRepo repository.UserRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, cfg: cfg}
}

func (h *AuthHandler) LoginPage(c *gin.Context) {
	// Redirect if already logged in
	session := sessions.Default(c)
	if session.Get(middleware.SessionUserKey) != nil {
		c.Redirect(http.StatusFound, "/admin")
		return
	}

	c.HTML(http.StatusOK, "admin_login.html", gin.H{
		"title": "Admin Login",
	})
}

func (h *AuthHandler) LoginSubmit(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := h.userRepo.FindByEmail(email)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "admin_login.html", gin.H{
			"title": "Admin Login",
			"error": "Invalid email or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "admin_login.html", gin.H{
			"title": "Admin Login",
			"error": "Invalid email or password",
		})
		return
	}

	session := sessions.Default(c)
	session.Set(middleware.SessionUserKey, user.ID.String())
	session.Save()

	c.Redirect(http.StatusFound, "/admin")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/admin/login")
}
