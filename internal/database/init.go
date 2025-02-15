package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Init() {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Enable foreign key constraints
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal(err)
	}

	// Create the workspace table
	createWorkspaceTableSQL := `CREATE TABLE IF NOT EXISTS workspace (
        "id" INTEGER PRIMARY KEY AUTOINCREMENT,
        "name" TEXT UNIQUE,
        "desc" TEXT,
        "total_containers" INTEGER DEFAULT 0,
        "created_at" DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
	_, err = db.Exec(createWorkspaceTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Create the container table
	createContainerTableSQL := `CREATE TABLE IF NOT EXISTS container (
        "id" INTEGER PRIMARY KEY AUTOINCREMENT,
        "workspace_id" INTEGER,
        "name" TEXT,
        "image" TEXT,
        "ports" TEXT,
        "environment_vars" TEXT,
        "cpus" TEXT,
        "memory" TEXT,
        "volumes" TEXT,
        "labels" TEXT,
        "created_at" DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(workspace_id) REFERENCES workspace(id) ON DELETE CASCADE
    );`
	_, err = db.Exec(createContainerTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Create the proxy table
	createProxyTableSQL := `CREATE TABLE IF NOT EXISTS proxy (
        "id" INTEGER PRIMARY KEY AUTOINCREMENT,
        "cube_id" INTEGER,
        "domain" TEXT,
        "port" INTEGER,
        "type" TEXT,
        "default" BOOLEAN,
        "created_at" DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(cube_id) REFERENCES container(id) ON DELETE CASCADE
    );`
	_, err = db.Exec(createProxyTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Tables created successfully!")
}

func GetPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}
	return fmt.Sprintf("%s.db", ex), nil
}
