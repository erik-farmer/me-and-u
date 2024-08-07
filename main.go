package main

import (
	"database/sql"
	"fmt"
	"github.com/erik-farmer/me-and-u/handlers"
	"github.com/erik-farmer/me-and-u/middleware"
	"html/template"
	"net/http"
	"os"

	"github.com/erik-farmer/me-and-u/data"
	"github.com/erik-farmer/me-and-u/db"
)

func main() {
	dir := db.MakeDbDir()
	defer os.RemoveAll(dir)

	connector := db.MakeDbConnector(dir)
	defer connector.Close()

	// ToDo: https://www.alexedwards.net/blog/organising-database-access
	// Method 2 is probably fine and good practice
	database := sql.OpenDB(connector)
	defer database.Close()

	mux := http.NewServeMux()

	// ToDo: Add a router? Maybe if they start to get complicated
	// ToDo: make handlers it's own package?

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		rows, err := database.Query("SELECT rowid, name, url FROM recipes")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
			os.Exit(1)
		}
		defer rows.Close()

		var recipes []data.Recipe
		for rows.Next() {
			var recipe data.Recipe

			if err := rows.Scan(&recipe.ROW_ID, &recipe.Name, &recipe.URL); err != nil {
				fmt.Println("Error scanning row:", err)
				return
			}
			recipes = append(recipes, recipe)
		}

		template_data := map[string]interface{}{
			"recipes": recipes,
		}
		tmpl := template.Must(template.ParseFiles("templates/recipe_list.html"))
		tmpl.Execute(w, template_data)
	})

	mux.HandleFunc("GET /recipe/{id}/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/recipe_detail.html"))
		idString := r.PathValue("id")

		row := database.QueryRow("SELECT * from recipes WHERE rowid=?", idString)
		recipe := data.Recipe{}
		error := row.Scan(&recipe.Name, &recipe.URL, &recipe.Ingredients, &recipe.Steps, &recipe.Notes)
		if error != nil {
			println(error.Error())
			// ToDo 404 here
		}

		template_data := map[string]data.Recipe{"recipe": recipe}
		tmpl.Execute(w, template_data)
	})

	mux.HandleFunc("GET /recipe/new/", handlers.NewRecipe)

	mux.HandleFunc("POST /recipe/new/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Validating that we hit POST")
		r.ParseForm()

		recipe := data.Recipe{
			Name:        r.FormValue("name"),
			URL:         r.FormValue("url"),
			Ingredients: r.FormValue("ingredients"),
			Steps:       r.FormValue("steps"),
			Notes:       r.FormValue("notes"),
		}

		stmt := "INSERT INTO recipes (name, url, ingredients, steps, notes) VALUES(?,?,?,?,?);"
		_, err := database.Exec(stmt, recipe.Name, recipe.URL, recipe.Ingredients, recipe.Steps, recipe.Notes)
		if err != nil {
			fmt.Printf("There was an error executing: \n%s\n", stmt)
			fmt.Printf("Error: \n%s\n", err)
		}

		tmpl := template.Must(template.ParseFiles("templates/new_recipe.html"))
		tmpl.Execute(w, nil)
	})

	middlewareStack := middleware.CreateStackMiddleware(
		middleware.Logging,
		middleware.Authorize,
	)
	server := http.Server{
		Addr:    ":8080",
		Handler: middlewareStack(mux),
	}
	server.ListenAndServe()
}
