package database

import (
    "database/sql"
    "fmt"
    "log"
    "strings"

    _ "github.com/mattn/go-sqlite3"
    "github.com/turplespace/portos/internal/models"
)

func InsertWorkspaceAndCubes(workspaceID int, cubes []models.Container) error {
    db_path, _ := GetPath()
    db, err := sql.Open("sqlite3", db_path)
    if err != nil {
        return fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    // Check if the workspace exists
    var exists bool
    err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM workspace WHERE id = ?)`, workspaceID).Scan(&exists)
    if err != nil {
        return fmt.Errorf("failed to check workspace existence: %v", err)
    }
    if !exists {
        return fmt.Errorf("workspace with ID %d does not exist", workspaceID)
    }

    // Insert the cubes
    for _, cube := range cubes {
        _, err := db.Exec(`INSERT INTO container (workspace_id, name, image, ports, environment_vars, cpus, memory, volumes, labels) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
            workspaceID, cube.Name, cube.Image, strings.Join(cube.Ports, ","), strings.Join(cube.EnvironmentVars, ","),
            cube.ResourceLimits.CPUs, cube.ResourceLimits.Memory, mapToString(cube.Volumes), strings.Join(cube.Labels, ","))
        if err != nil {
            return fmt.Errorf("failed to insert cube: %v", err)
        }
    }

    log.Println("Workspace and cubes inserted successfully!")
    return nil
}


// Helper function to convert map to string
func mapToString(m map[string]string) string {
    var sb strings.Builder
    for k, v := range m {
        sb.WriteString(fmt.Sprintf("%s:%s,", k, v))
    }
    result := sb.String()
    if len(result) > 0 {
        result = result[:len(result)-1] // Remove trailing comma
    }
    return result
}

func GetCubeData(cubeID int) (*models.Container, error) {
    db_path, _ := GetPath()
    db, err := sql.Open("sqlite3", db_path)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    query := `SELECT name, image, ports, environment_vars, cpus, memory, volumes, labels 
              FROM container WHERE id = ?`
    row := db.QueryRow(query, cubeID)

    var cube models.Container
    var ports, envVars, volumes, labels string

    err = row.Scan(&cube.Name, &cube.Image, &ports, &envVars, &cube.ResourceLimits.CPUs, &cube.ResourceLimits.Memory, &volumes, &labels)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("cube with ID %d not found", cubeID)
        }
        return nil, fmt.Errorf("failed to query cube data: %v", err)
    }

    cube.Ports = splitString(ports)
    cube.EnvironmentVars = splitString(envVars)
    cube.Volumes = stringToMap(volumes)
    cube.Labels = splitString(labels)

    return &cube, nil
}

func splitString(s string) []string {
    if s == "" {
        return nil
    }
    return strings.Split(s, ",")
}

func stringToMap(s string) map[string]string {
    m := make(map[string]string)
    pairs := strings.Split(s, ",")
    for _, pair := range pairs {
        kv := strings.Split(pair, ":")
        if len(kv) == 2 {
            m[kv[0]] = kv[1]
        }
    }
    return m
}

func UpdateCube(cubeID int, updatedCube models.Container) error {
    db_path, _ := GetPath()
    db, err := sql.Open("sqlite3", db_path)
    if err != nil {
        return fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    query := `UPDATE container SET name = ?, image = ?, ports = ?, environment_vars = ?, cpus = ?, memory = ?, volumes = ?, labels = ? WHERE id = ?`
    _, err = db.Exec(query, updatedCube.Name, updatedCube.Image, strings.Join(updatedCube.Ports, ","), strings.Join(updatedCube.EnvironmentVars, ","),
        updatedCube.ResourceLimits.CPUs, updatedCube.ResourceLimits.Memory, mapToString(updatedCube.Volumes), strings.Join(updatedCube.Labels, ","), cubeID)
    if err != nil {
        return fmt.Errorf("failed to update cube: %v", err)
    }

    log.Println("Cube updated successfully!")
    return nil
}

func DeleteCube(cubeID int) error {
    db_path, _ := GetPath()
    db, err := sql.Open("sqlite3", db_path)
    if err != nil {
        return fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    query := `DELETE FROM container WHERE id = ?`
    _, err = db.Exec(query, cubeID)
    if err != nil {
        return fmt.Errorf("failed to delete cube: %v", err)
    }

    log.Println("Cube deleted successfully!")
    return nil
}

func ListCubes(workspaceID int) ([]models.Container, error) {
    db_path, _ := GetPath()
    db, err := sql.Open("sqlite3", db_path)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    query := `SELECT id, name, image, ports FROM container WHERE workspace_id = ?`
    rows, err := db.Query(query, workspaceID)
    if err != nil {
        return nil, fmt.Errorf("failed to query cubes: %v", err)
    }
    defer rows.Close()

    var cubes []models.Container
    for rows.Next() {
        var cube models.Container
        var ports, envVars, volumes, labels string

        err = rows.Scan(&cube.ID, &cube.Name, &cube.Image, &ports)
        if err != nil {
            return nil, fmt.Errorf("failed to scan cube: %v", err)
        }

        cube.Ports = splitString(ports)
        cube.EnvironmentVars = splitString(envVars)
        cube.Volumes = stringToMap(volumes)
        cube.Labels = splitString(labels)

        cubes = append(cubes, cube)
    }

    return cubes, nil
}

func CountContainersByWorkspaceID(workspaceID int) (int, error) {
    dbPath, _ := GetPath()
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return 0, fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    var count int
    query := `SELECT COUNT(*) FROM container WHERE workspace_id = ?`
    err = db.QueryRow(query, workspaceID).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to count containers for workspace %d: %v", workspaceID, err)
    }

    return count, nil
}

func CountCubes() (int, error) {
    db_path, _ := GetPath()
    db, err := sql.Open("sqlite3", db_path)
    if err != nil {
        return 0, fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    var count int
    query := `SELECT COUNT(*) FROM container`
    err = db.QueryRow(query).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to count cubes: %v", err)
    }
    return count, nil
}
// DeleteContainersByWorkspaceID deletes all containers associated with a specific workspace by its ID
func DeleteContainersByWorkspaceID(workspaceID int) error {
    db_path, _ := GetPath()
    db, err := sql.Open("sqlite3", db_path)
    if err != nil {
        return fmt.Errorf("failed to open database: %v", err)
    }
    defer db.Close()

    query := `DELETE FROM container WHERE workspace_id = ?`
    _, err = db.Exec(query, workspaceID)
    if err != nil {
        return fmt.Errorf("failed to delete containers for workspace %d: %v", workspaceID, err)
    }

    log.Printf("Deleted containers for workspace %d successfully!", workspaceID)
    return nil
}