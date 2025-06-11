package handlers

import (
	"database/sql"
	"fmt"
	"github.com/erik-farmer/me-and-u/data"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func LoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		user, err := data.GetUserByUsername(db, username)
		if err != nil {
			c.String(http.StatusUnauthorized, "Invalid username or password")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			c.String(http.StatusUnauthorized, "Invalid username or password")
			return
		}

		// ToDo: Create JWT and set it in a cookie

		c.String(http.StatusOK, fmt.Sprintf("Welcome, %s!", user.Username))
	}
}
