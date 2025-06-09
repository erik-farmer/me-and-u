package main

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/erik-farmer/me-and-u/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"html/template"
	"log/slog"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
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
	slog.Info(fmt.Sprintf("Application set to port %s", applicationPort))

	// DB setup
	// ToDo: Clean this up into it's own module
	url := os.Getenv("TURSO_URL")
	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	defer db.Close()
	// End db setup

	env := &Env{db: db}

	// Router SetUp
	releaseMode := gin.DebugMode
	if envMode := os.Getenv("RELEASE_MODE"); envMode == "PROD" {
		releaseMode = gin.ReleaseMode
	}
	gin.SetMode(releaseMode)
	//store := cookie.NewStore([]byte("secret"))
	router := gin.Default()
	//router.Use(sessions.Sessions("auth-session", store))
	//router.LoadHTMLGlob("./templates/*")
	templ := template.Must(template.New("").ParseFS(f, "templates/*.html"))
	router.SetHTMLTemplate(templ)

	// List
	router.GET("/", handlers.ListRecipesHandler(env.db))
	// Detail
	router.GET("/recipes/:recipe_id/", handlers.RecipeDetailHandler(env.db))
	//// New
	//router.GET("/recipes/new/", handlers.NewRecipeForm)
	//// ToDo: Create db entry from posted data
	//router.POST("/recipes/new/", handlers.CreateRecipeFromForm(db))

	// User Auth
	//router.GET("/login", handlers.LoginUser(auth))
	//router.GET("/login_callback", handlers.LoginCallback(auth))
	//router.GET("/logout", handlers.Logout)

	slog.Info(fmt.Sprintf("Application running on port %s", applicationPort))
	router.Run(applicationPort)
}
