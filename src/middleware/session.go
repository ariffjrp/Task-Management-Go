package middleware

import (
	"encoding/gob"
	"net/http"
	"task_management/src/entity"

	"github.com/gin-gonic/gin"
	gsessions "github.com/gorilla/sessions"
)

var store *gsessions.CookieStore

func InitSessionStore(secret []byte) {

	gob.Register(&entity.User{})
	gob.Register(&entity.Account{})

	store = gsessions.NewCookieStore(secret)
	store.Options = &gsessions.Options{
		Path:     "/",
		MaxAge:   60 * 60, // 1 minute
		HttpOnly: true,
		Secure:   false,
	}
}

func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, "mysession")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get session"})
			c.Abort()
			return
		}

		c.Set("session", session)
		c.Next()
	}
}

func GetSession(c *gin.Context) *gsessions.Session {
	if session, exists := c.Get("session"); exists {
		return session.(*gsessions.Session)
	}
	return nil
}
