package web

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/erik-farmer/me-and-u/data"
	"github.com/erik-farmer/me-and-u/web/handlers"
	"github.com/erik-farmer/me-and-u/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"log/slog"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Env struct {
	db *sql.DB
}

//go:embed templates/*.html
var tmplFS embed.FS

func StartServer() {
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

	db, err := data.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	env := &Env{db: db}

	// Router SetUp
	releaseMode := gin.DebugMode
	if envMode := os.Getenv("RELEASE_MODE"); envMode == "PROD" {
		releaseMode = gin.ReleaseMode
	}
	gin.SetMode(releaseMode)
	router := gin.Default()
	templ := template.Must(template.New("").ParseFS(tmplFS, "templates/*.html"))
	router.SetHTMLTemplate(templ)

	router.Use(middleware.InjectUsername())

	// List
	router.GET("/", handlers.ListRecipesHandler(env.db))
	// Detail
	router.GET("/recipes/:recipe_id/", handlers.RecipeDetailHandler(env.db))
	//// New
	router.GET("/recipes/new/", middleware.RequireAuth, handlers.NewRecipeForm)
	//// ToDo: Create db entry from posted data
	router.POST("/recipes/new/", handlers.CreateRecipeFromForm(db))

	// User Auth
	router.GET("/login", handlers.LoginForm)
	router.POST("/login", handlers.LoginUser(env.db))
	router.GET("/logout", handlers.Logout)

	slog.Info(fmt.Sprintf("Application running on port %s", applicationPort))
	router.Run(applicationPort)
}
