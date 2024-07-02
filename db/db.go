package db

import (
	"fmt"
	"github.com/tursodatabase/go-libsql"
	"os"
	"path/filepath"
)

func MakeDbDir() string {
	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	return dir
}

func MakeDbConnector(dir string) *libsql.Connector {
	dbName := "local.db"
	primaryUrl := os.Getenv("TURSO_URL")
	authToken := os.Getenv("TURSO_TOKEN")

	dbPath := filepath.Join(dir, dbName)
	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
	)
	if err != nil {
		fmt.Println("Error creating connector:", err)
		os.Exit(1)
	}

	return connector
}
