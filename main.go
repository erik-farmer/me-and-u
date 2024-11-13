package main

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tursodatabase/go-libsql"

	authenticator "github.com/erik-farmer/me-and-u/auth"
	"github.com/erik-farmer/me-and-u/handlers"
	"github.com/erik-farmer/me-and-u/middleware"
)

// Create a custom Env struct which holds a connection pool.
type Env struct {
	db *sql.DB
}

//go:embed templates
var f embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file :(")
	}

	// Default to 8080 unless port specified. non-https (8000) required for local dev with auth
	applicationPort := ":8080"
	if port := os.Getenv("APPLICATION_PORT"); port != "" {
		applicationPort = port
	}
	slog.Info(fmt.Sprintf("Application running on port %s", applicationPort))

	// DB setup
	// ToDo: Clean this up into it's own module
	dbName := "local.db"
	primaryUrl := os.Getenv("TURSO_URL")
	authToken := os.Getenv("TURSO_TOKEN")

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)
	dbPath := filepath.Join(dir, dbName)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
	)
	if err != nil {
		fmt.Println("Error creating connector:", err)
		os.Exit(1)
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()
	// End db setup

	env := &Env{db: db}

	// Rouer SetUp
	releaseMode := gin.DebugMode
	if envMode := os.Getenv("RELEASE_MODE"); envMode == "PROD" {
		releaseMode = gin.ReleaseMode
	}
	gin.SetMode(releaseMode)
	store := cookie.NewStore([]byte("secret"))
	router := gin.Default()
	router.Use(sessions.Sessions("auth-session", store))
	//router.LoadHTMLGlob("./templates/*")
	templ := template.Must(template.New("").ParseFS(f, "templates/*.html"))
	router.SetHTMLTemplate(templ)

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

	router.Run(applicationPort)
}
