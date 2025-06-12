package middleware

import (
	"github.com/erik-farmer/me-and-u/web/handlers"
	"github.com/gin-gonic/gin"
)

func InjectUsername() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("auth_token")
		if err == nil {
			if username, err := handlers.ParseToken(tokenStr); err == nil {
				c.Set("Username", username)
			}
		}
		c.Next()
	}
}
