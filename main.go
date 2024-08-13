package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/erik-farmer/me-and-u/data"
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
	router.GET("/", func(c *gin.Context) {
		recipes, err := data.AllRecipes(env.db)
		if err != nil {
			c.String(http.StatusInternalServerError, "Unable to retrieve Recipes")
		}

		c.HTML(http.StatusOK, "recipe_list.html", gin.H{
			"recipes": recipes,
		})
	})
	// Detail
	router.GET("/recipes/:recipe_id/", func(c *gin.Context) {
		recipe_id := c.Param("recipe_id")
		fmt.Fprintf(os.Stdout, "Router searching for recipe %s\n", recipe_id)
		recipe, err := data.GetRecipeById(env.db, recipe_id)
		if err != nil {
			println(err.Error())
			c.String(http.StatusNotFound, "Unable to retrieve Recipe with provided ID")
		}

		c.HTML(http.StatusOK, "recipe_detail.html", gin.H{
			"recipe": recipe,
		})
	})
	// New
	router.GET("/recipes/new/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Recipe Form Will Go Here",
		})
	})

	router.Run(":8080")
}
