package main

import (
	"fmt"
	"html/template"
	"net/http"
)


func main() {
    db := init_db()
    defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World! you have hit the root")
	})

	mux.HandleFunc("GET /html", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/recipe_list.html"))
		title := "Hi Bby"
		data := map[string]string{
			"title": title,
		}
		tmpl.Execute(w, data)
	})

	mux.HandleFunc("GET /recipes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "You are looking for recipe: %s", id)
	})

	http.ListenAndServe(":8080", mux)
}
