package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	authenticator "github.com/erik-farmer/me-and-u/auth"
	"github.com/erik-farmer/me-and-u/handlers"
	"github.com/erik-farmer/me-and-u/middleware"
)

// Create a custom Env struct which holds a connection pool.
type Env struct {
	db *sql.DB
}

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file :(")
	}
	// DB setup
	primaryUrl := os.Getenv("TURSO_URL")
	authToken := os.Getenv("TURSO_TOKEN")
	url := primaryUrl + "?authToken=" + authToken

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open db err: %s", err)
		os.Exit(1)
	}

	env := &Env{db: db}

	// SetUp
	store := cookie.NewStore([]byte("secret"))
	router := gin.Default()
	router.Use(sessions.Sessions("auth-session", store))
	router.LoadHTMLGlob("templates/*")

	// Authenticator
	auth, err := authenticator.New()
	if err != nil {
		log.Printf("Failed to initialize the authenticator: %v", err)
	}

	// List
	router.GET("/", handlers.ListRecipesHandler(env.db))
	// Detail
	router.GET("/recipes/:recipe_id/", handlers.RecipeDetailHandler(env.db))
	// New
	router.GET("/recipes/new/", middleware.IsAuthenticated, handlers.NewRecipeForm)
	// ToDo: Create db entry from posted data
	router.POST("/recipes/new/", handlers.CreateRecipeFromForm(db))

	// User Auth
	router.GET("/login", handlers.LoginUser(auth))
	router.GET("/login_callback", handlers.LoginCallback(auth))
	router.GET("/logout", handlers.Logout)

	router.Run(":8080")
}
