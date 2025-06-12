package handlers

import (
	"database/sql"
	"fmt"
	"github.com/erik-farmer/me-and-u/data"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func LoginForm(c *gin.Context) {
	tokenStr, err := c.Cookie("auth_token")
	if err == nil {
		_, err := ParseToken(tokenStr)
		if err == nil {
			//Valid Token
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
	}
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

		token, err := CreateToken(username)
		c.SetCookie(
			"auth_token", token,
			86400, // maxAge (24h)
			"/",   // path
			"",    // domain (empty = current)
			true,  // secure
			true,  // httpOnly
		)

		c.Redirect(http.StatusSeeOther, "/")
	}
}

func CreateToken(username string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr string) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil {
		return "", fmt.Errorf("invalid claims")
	}

	username, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("sub is not a string")
	}

	return username, nil
}

func Logout(c *gin.Context) {
	c.SetCookie(
		"auth_token", "",
		-1,
		"/",
		"",
		true,
		true,
	)
	c.Redirect(http.StatusSeeOther, "/")
}
