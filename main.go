package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/erik-farmer/me-and-u/db"
)

func main() {
	db := db.InitDB()
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World! you have hit the root")
	})

	mux.HandleFunc("GET /recipes/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/recipe_list.html"))
		tmpl.Execute(w, nil)
	})

	mux.HandleFunc("GET /recipe/new/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/new_recipe.html"))
		tmpl.Execute(w, nil)
	})

	mux.HandleFunc("POST /recipe/new/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Validating that we hit POST")
		r.ParseForm()
		fmt.Printf("NAME => %s\n", r.FormValue("name"))

		tmpl := template.Must(template.ParseFiles("templates/new_recipe.html"))
		tmpl.Execute(w, nil)
	})

	http.ListenAndServe(":8080", mux)
}
