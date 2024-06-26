package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tursodatabase/go-libsql"
)

func init_db() *sql.DB {
	dbName := "local.db"
	dbUrl := "libsql://me-and-u-dev-efdev.turso.io"
	authToken := os.Getenv("TURSO_TOKEN")

	dir, err := os.MkdirTemp("", "libsql-*")
    if err != nil {
        fmt.Println("Error creating temporary directory:", err)
        os.Exit(1)
    }
    defer os.RemoveAll(dir)

    dbPath := filepath.Join(dir, dbName)

    connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, dbUrl,
        libsql.WithAuthToken(authToken),
    )
    if err != nil {
        fmt.Println("Error creating connector:", err)
        os.Exit(1)
    }

    db := sql.OpenDB(connector)
    return db
}

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
