package handlers

import (
	"database/sql"
	"fmt"
	"github.com/erik-farmer/me-and-u/data"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func ListRecipesHandler(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		recipes, err := data.AllRecipes(db)
		if err != nil {
			c.String(http.StatusInternalServerError, "Unable to retrieve Recipes")
		}

		session := sessions.Default(c)
		profile := session.Get("profile")
		c.HTML(http.StatusOK, "recipe_list.html", gin.H{
			"recipes": recipes,
			"profile": profile,
		})
	}

	return fn
}

func RecipeDetailHandler(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		recipe_id := c.Param("recipe_id")
		fmt.Fprintf(os.Stdout, "Router searching for recipe %s\n", recipe_id)
		recipe, err := data.GetRecipeById(db, recipe_id)
		if err != nil {
			println(err.Error())
			c.String(http.StatusNotFound, "Unable to retrieve Recipe with provided ID")
		}

		c.HTML(http.StatusOK, "recipe_detail.html", gin.H{
			"recipe": recipe,
		})
	}

	return fn
}

func NewRecipeForm(c *gin.Context) {
	c.HTML(http.StatusOK, "new_recipe.html", gin.H{})
}

func CreateRecipeFromForm(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		recipe := data.Recipe{}
		c.Bind(&recipe)

		// Insert our form data into the DB
		stmt := "INSERT INTO recipes (name, url, ingredients, steps, notes) VALUES(?,?,?,?,?);"
		result, err := db.Exec(stmt, recipe.Name, recipe.URL, recipe.Ingredients, recipe.Steps, recipe.Notes)
		if err != nil {
			fmt.Printf("There was an error executing: \n%s\n", stmt)
			fmt.Printf("Error: \n%s\n", err)
		}

		// Get the rowid for the redirect to the detail view and do said redirect
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			log.Fatalf("Failed to get last insert id: %v", err)
		}
		log.Printf("Inserted record with rowid: %d", lastInsertID)
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/recipes/%d", lastInsertID))
	}

	return fn
}
