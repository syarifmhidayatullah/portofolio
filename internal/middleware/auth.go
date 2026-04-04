package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const SessionUserKey = "user_id"

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get(SessionUserKey)
		if userID == nil {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}
