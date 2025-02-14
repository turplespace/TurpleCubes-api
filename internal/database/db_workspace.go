package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/turplespace/portos/internal/models"
)

type Workspace struct {
	ID              int        `json:"id"`
	Name            string     `json:"name"`
	Desc            string     `json:"desc"`
	TotalContainers int        `json:"total_containers"`
	CreatedAt       *time.Time `json:"created_at"`
}

// CreateWorkspace to create a new workspace
func CreateWorkspace(name string, desc string) error {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	insertSQL := `INSERT INTO workspace (name, desc) VALUES (?, ?)`
	_, err = db.Exec(insertSQL, name, desc)
	return err
}

// IncContainerCount to update the container count for a specific workspace
func IncContainerCount(workspaceName string, count int) error {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	// Update the total_containers for the specified workspace
	updateSQL := `UPDATE workspace SET total_containers = total_containers + ? WHERE name = ?`
	_, err = db.Exec(updateSQL, count, workspaceName)
	if err != nil {
		return fmt.Errorf("failed to update container count: %v", err)
	}

	log.Printf("Updated container count for workspace %s to %d", workspaceName, count)
	return nil
}

// DecContainerCount to update the container count for a specific workspace
func DecContainerCount(workspaceName string, count int) error {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	// Update the total_containers for the specified workspace
	updateSQL := `UPDATE workspace SET total_containers = total_containers - ? WHERE name = ?`
	_, err = db.Exec(updateSQL, count, workspaceName)
	if err != nil {
		return fmt.Errorf("failed to update container count: %v", err)
	}

	log.Printf("Updated container count for workspace %s to %d", workspaceName, count)
	return nil
}

// GetWorkspaces to fetch all the workspaces
func GetWorkspaces() ([]Workspace, error) {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, desc, total_containers, created_at FROM workspace")
	if err != nil {
		return nil, fmt.Errorf("failed to query workspaces: %v", err)
	}
	defer rows.Close()

	var workspaces []Workspace
	// Iterate over the rows and scan the data into a Workspace struct
	for rows.Next() {
		var workspace Workspace

		err = rows.Scan(&workspace.ID, &workspace.Name, &workspace.Desc, &workspace.TotalContainers, &workspace.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workspace: %v", err)
		}
		workspaces = append(workspaces, workspace)
	}
	return workspaces, nil
}

// DeleteWorkspace to delete a workspace
func DeleteWorkspace(id int) error {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	deleteSQL := `DELETE FROM workspace WHERE id = ?`
	_, err = db.Exec(deleteSQL, id)
	if err != nil {
		return fmt.Errorf("failed to delete workspace: %v", err)
	}
	return nil
}

// EditWorkspace to edit a workspace
func EditWorkspace(id int, name string, desc string) error {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	updateSQL := `UPDATE workspace SET name = ?, desc = ?  WHERE id = ?`
	_, err = db.Exec(updateSQL, name, desc, id)
	if err != nil {
		return fmt.Errorf("failed to update workspace: %v", err)
	}

	log.Printf("Workspace %d updated successfully", id)
	return nil
}

// ListContainersInWorkspace to list all containers in a workspace
func ListContainersInWorkspace(workspaceID int) ([]models.Container, error) {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	query := `SELECT id, name, image, ports, environment_vars, cpus, memory, volumes, labels 
              FROM container WHERE workspace_id = ?`
	rows, err := db.Query(query, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query containers: %v", err)
	}
	defer rows.Close()

	var containers []models.Container
	for rows.Next() {
		var container models.Container
		var ports, envVars, volumes, labels string

		err = rows.Scan(&container.ID, &container.Name, &container.Image, &ports, &envVars, &container.ResourceLimits.CPUs, &container.ResourceLimits.Memory, &volumes, &labels)
		if err != nil {
			return nil, fmt.Errorf("failed to scan container: %v", err)
		}

		container.Ports = splitString(ports)
		container.EnvironmentVars = splitString(envVars)
		container.Volumes = stringToMap(volumes)
		container.Labels = splitString(labels)

		containers = append(containers, container)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	log.Printf("Fetched %d containers for workspace %d", len(containers), workspaceID)
	return containers, nil
}

// CountWorkspaces to count the number of workspaces
func CountWorkspaces() (int, error) {
	db_path, _ := GetPath()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return 0, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	var count int
	query := `SELECT COUNT(*) FROM workspace`
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count workspaces: %v", err)
	}
	return count, nil
}
