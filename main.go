package main

import (
	"database/sql"
	"fmt"
	"github.com/erik-farmer/me-and-u/handlers"
	"github.com/gin-gonic/gin"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// Create a custom Env struct which holds a connection pool.
type Env struct {
	db *sql.DB
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// List
	router.GET("/", handlers.ListRecipesHandler(env.db))
	// Detail
	router.GET("/recipes/:recipe_id/", handlers.RecipeDetailHandler(env.db))
	// New
	router.GET("/recipes/new/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Recipe Form Will Go Here",
		})
	})

	router.Run(":8080")
}
