package data

import (
	"database/sql"
	"fmt"
)

type Recipe struct {
	ROW_ID      int
	Name        string `validate:"required" form:"name"`
	URL         string `form:"url"`
	Ingredients string `form:"ingredients"`
	Steps       string `form:"steps"`
	Notes       string `form:"notes"`
}

func AllRecipes(db *sql.DB) ([]Recipe, error) {
	rows, err := db.Query("SELECT ROWID, Name, URL FROM recipes;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		recipe := Recipe{}
		if err := rows.Scan(&recipe.ROW_ID, &recipe.Name, &recipe.URL); err != nil {
			fmt.Println("Error scanning row:", err)
			return recipes, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func GetRecipeById(db *sql.DB, id string) (Recipe, error) {
	row := db.QueryRow("SELECT * from recipes WHERE rowid=?", id)
	recipe := Recipe{}
	err := row.Scan(&recipe.Name, &recipe.URL, &recipe.Ingredients, &recipe.Steps, &recipe.Notes)
	if err != nil {
		println(err.Error())
		return recipe, err
	}
	return recipe, nil
}
