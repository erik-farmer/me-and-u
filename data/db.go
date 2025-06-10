package data

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func NewDB() (*sql.DB, error) {
	url := os.Getenv("TURSO_URL")
	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, fmt.Errorf("failed to open db %s: %w", url, err)
	}
	return db, nil
}
