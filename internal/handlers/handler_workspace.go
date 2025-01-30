package handlers

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"

    "github.com/turplespace/portos/internal/database"
    "github.com/turplespace/portos/internal/models"
    "github.com/turplespace/portos/internal/services/docker"
)

// WorkspaceResponse holds the total counts and the list of workspaces with container counts
type WorkspaceResponse struct {
    TotalWorkspaces    int                                `json:"total_workspaces"`
    TotalCubes         int                                `json:"total_cubes"`
    TotalRunningCubes  int                                `json:"total_running_cubes"`
    Workspaces         []models.WorkspaceWithContainerCounts `json:"workspaces"`
}

// HandleGetWorkspaces handles the HTTP request to get the list of workspaces with container counts
func HandleGetWorkspaces(w http.ResponseWriter, r *http.Request) {
    log.Println("[GET-WORKSPACES] Starting get workspaces request")

    if r.Method != http.MethodGet {
        log.Printf("[GET-WORKSPACES] Error: Invalid method %s used instead of GET", r.Method)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Retrieve the list of workspaces
    workspaces, err := database.GetWorkspaces()
    if err != nil {
        log.Printf("[GET-WORKSPACES] Error: Failed to get workspaces: %v", err)
        http.Error(w, fmt.Sprintf("Failed to get workspaces: %v", err), http.StatusInternalServerError)
        return
    }

    // Get the total counts
    totalWorkspaces, err := database.CountWorkspaces()
    if err != nil {
        log.Printf("[GET-WORKSPACES] Error: Failed to count workspaces: %v", err)
        http.Error(w, fmt.Sprintf("Failed to count workspaces: %v", err), http.StatusInternalServerError)
        return
    }

    totalCubes, err := database.CountCubes()
    if err != nil {
        log.Printf("[GET-WORKSPACES] Error: Failed to count cubes: %v", err)
        http.Error(w, fmt.Sprintf("Failed to count cubes: %v", err), http.StatusInternalServerError)
        return
    }

    totalRunningCubes, err := docker.CountContainersByLabel("service", "turplespace")
    if err != nil {
        log.Printf("[GET-WORKSPACES] Error: Failed to count running containers: %v", err)
        http.Error(w, fmt.Sprintf("Failed to count running containers: %v", err), http.StatusInternalServerError)
        return
    }

    // Create a slice to hold workspaces with container counts
    var workspacesWithCounts []models.WorkspaceWithContainerCounts

    // For each workspace, count the number of total and running containers
    for _, workspace := range workspaces {
        totalCount, err := database.CountContainersByWorkspaceID(workspace.ID)
        if err != nil {
            log.Printf("[GET-WORKSPACES] Error: Failed to count containers for workspace %d: %v", workspace.ID, err)
            http.Error(w, fmt.Sprintf("Failed to count containers for workspace %d: %v", workspace.ID, err), http.StatusInternalServerError)
            return
        }

        runningCount, err := docker.CountContainersByLabel("workspace_id", fmt.Sprintf("%d", workspace.ID))
        if err != nil {
            log.Printf("[GET-WORKSPACES] Error: Failed to count running containers for workspace %d: %v", workspace.ID, err)
            http.Error(w, fmt.Sprintf("Failed to count running containers for workspace %d: %v", workspace.ID, err), http.StatusInternalServerError)
            return
        }

        workspacesWithCounts = append(workspacesWithCounts, models.WorkspaceWithContainerCounts{
            ID:                workspace.ID,
            Name:              workspace.Name,
            Desc:              workspace.Desc,
            TotalContainers:   totalCount,
            RunningContainers: runningCount,
            CreatedAt:         workspace.CreatedAt,
        })
    }

    // Create the response object
    response := WorkspaceResponse{
        TotalWorkspaces:   totalWorkspaces,
        TotalCubes:        totalCubes,
        TotalRunningCubes: totalRunningCubes,
        Workspaces:        workspacesWithCounts,
    }

    // Encode the response as JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
    log.Println("[GET-WORKSPACES] Successfully sent workspaces response")
}

func HandleCreateWorkspace(w http.ResponseWriter, r *http.Request) {
    log.Println("[CREATE-WORKSPACE] Starting create workspace request")

    if r.Method != http.MethodPost {
        log.Printf("[CREATE-WORKSPACE] Error: Invalid method %s used instead of POST", r.Method)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        Name  string `json:"name"`
        Desc  string `json:"desc"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("[CREATE-WORKSPACE] Error: Invalid request body - %v", err)
        http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
        return
    }

    err := database.CreateWorkspace(req.Name, req.Desc)
    if err != nil {
        log.Printf("[CREATE-WORKSPACE] Error: Failed to create workspace: %v", err)
        http.Error(w, fmt.Sprintf("Failed to create workspace: %v", err), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Workspace created successfully"})
    log.Println("[CREATE-WORKSPACE] Successfully created workspace")
}

func HandleEditWorkspace(w http.ResponseWriter, r *http.Request) {
    log.Println("[EDIT-WORKSPACE] Starting edit workspace request")

    if r.Method != http.MethodPut {
        log.Printf("[EDIT-WORKSPACE] Error: Invalid method %s used instead of PUT", r.Method)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        ID    int    `json:"id"`
        Name  string `json:"name"`
        Desc  string `json:"desc"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("[EDIT-WORKSPACE] Error: Invalid request body - %v", err)
        http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
        return
    }

    err := database.EditWorkspace(req.ID, req.Name, req.Desc)
    if err != nil {
        log.Printf("[EDIT-WORKSPACE] Error: Failed to edit workspace: %v", err)
        http.Error(w, fmt.Sprintf("Failed to edit workspace: %v", err), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Workspace updated successfully"})
    log.Printf("[EDIT-WORKSPACE] Successfully updated workspace ID: %d", req.ID)
}

func HandleDeleteWorkspace(w http.ResponseWriter, r *http.Request) {
    log.Println("[DELETE-WORKSPACE] Starting delete workspace request")

    if r.Method != http.MethodDelete {
        log.Printf("[DELETE-WORKSPACE] Error: Invalid method %s used instead of DELETE", r.Method)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        log.Println("[DELETE-WORKSPACE] Error: Missing workspace ID")
        http.Error(w, "Missing workspace ID", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        log.Printf("[DELETE-WORKSPACE] Error: Invalid workspace ID format: %s - %v", idStr, err)
        http.Error(w, "Invalid workspace ID", http.StatusBadRequest)
        return
    }

    // Stopping Cubes
    cubes, err := database.ListCubes(id)
    if err != nil {
        log.Printf("[DELETE-WORKSPACE] Error: Failed to get cubes for workspace ID %d: %v", id, err)
        http.Error(w, fmt.Sprintf("Failed to get cubes: %v", err), http.StatusInternalServerError)
        return
    }

    for _, cube := range cubes {
        err = docker.StopContainer(cube.Name)
        if err != nil {
            log.Printf("[DELETE-WORKSPACE] Warning: Error stopping container %s: %v", cube.Name, err)
        }
    }

    // Deleting the Cubes from DB
    err = database.DeleteContainersByWorkspaceID(id)
    if err != nil {
        log.Printf("[DELETE-WORKSPACE] Error: Failed to delete containers for workspace ID %d: %v", id, err)
        http.Error(w, fmt.Sprintf("Failed to delete containers: %v", err), http.StatusInternalServerError)
        return
    }

    // Deleting Workspace from the DB
    err = database.DeleteWorkspace(id)
    if err != nil {
        log.Printf("[DELETE-WORKSPACE] Error: Failed to delete workspace ID %d: %v", id, err)
        http.Error(w, fmt.Sprintf("Failed to delete workspace: %v", err), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Workspace deleted successfully"})
    log.Printf("[DELETE-WORKSPACE] Successfully deleted workspace ID: %d", id)
}