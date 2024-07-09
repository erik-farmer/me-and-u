package handlers

import (
	"html/template"
	"net/http"
)

func NewRecipe(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/new_recipe.html"))
	tmpl.Execute(w, nil)
}
