package db

import (
	"database/sql"
	"fmt"
	"github.com/tursodatabase/go-libsql"
	"os"
	"path/filepath"
)

func InitDB() *sql.DB {
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
