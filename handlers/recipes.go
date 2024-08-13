package handlers

import (
	"database/sql"
	"fmt"
	"github.com/erik-farmer/me-and-u/data"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
)

func ListRecipesHandler(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		recipes, err := data.AllRecipes(db)
		if err != nil {
			c.String(http.StatusInternalServerError, "Unable to retrieve Recipes")
		}

		c.HTML(http.StatusOK, "recipe_list.html", gin.H{
			"recipes": recipes,
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

func NewRecipe(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/new_recipe.html"))
	tmpl.Execute(w, nil)
}
