package middleware

import (
	"github.com/erik-farmer/me-and-u/web/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAuth(ctx *gin.Context) {
	tokenStr, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/login")
		ctx.Abort()
		return
	}

	username, err := handlers.ParseToken(tokenStr)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/login")
		ctx.Abort()
		return
	}

	// Optional: inject username into context here too
	ctx.Set("Username", username)
	ctx.Next()
}
