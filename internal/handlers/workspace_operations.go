package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "github.com/turplespace/portos/internal/database"
    "github.com/turplespace/portos/internal/services/docker"
)

func HandleDeployWorkspace(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get workspace ID from query parameters
    workspaceIDStr := r.URL.Query().Get("workspace_id")
    if workspaceIDStr == "" {
        http.Error(w, "Missing workspace ID", http.StatusBadRequest)
        return
    }

    workspaceID, err := strconv.Atoi(workspaceIDStr)
    if err != nil {
        http.Error(w, "Invalid workspace ID", http.StatusBadRequest)
        return
    }

    // Get all containers in the workspace
    containers, err := database.ListContainersInWorkspace(workspaceID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to list containers: %v", err), http.StatusInternalServerError)
        return
    }

    // Start each container
    for _, container := range containers {
        err := docker.StartContainer(container)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to deploy container %s: %v", container.Name, err), http.StatusInternalServerError)
            return
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Workspace deployed successfully"})
}

func HandleRedeployWorkspace(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get workspace ID from query parameters
    workspaceIDStr := r.URL.Query().Get("workspace_id")
    if workspaceIDStr == "" {
        http.Error(w, "Missing workspace ID", http.StatusBadRequest)
        return
    }

    workspaceID, err := strconv.Atoi(workspaceIDStr)
    if err != nil {
        http.Error(w, "Invalid workspace ID", http.StatusBadRequest)
        return
    }

    // Get all containers in the workspace
    containers, err := database.ListContainersInWorkspace(workspaceID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to list containers: %v", err), http.StatusInternalServerError)
        return
    }

    // Redeploy each container
    for _, container := range containers {
        err := docker.RestartContainer(container.Name)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to redeploy container %s: %v", container.Name, err), http.StatusInternalServerError)
            return
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Workspace redeployed successfully"})
}

func HandleStopWorkspace(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get workspace ID from query parameters
    workspaceIDStr := r.URL.Query().Get("workspace_id")
    if workspaceIDStr == "" {
        http.Error(w, "Missing workspace ID", http.StatusBadRequest)
        return
    }

    workspaceID, err := strconv.Atoi(workspaceIDStr)
    if err != nil {
        http.Error(w, "Invalid workspace ID", http.StatusBadRequest)
        return
    }

    // Get all containers in the workspace
    containers, err := database.ListContainersInWorkspace(workspaceID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to list containers: %v", err), http.StatusInternalServerError)
        return
    }

    // Stop each container
    for _, container := range containers {
        err := docker.StopContainer(container.Name)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to stop container %s: %v", container.Name, err), http.StatusInternalServerError)
            return
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Workspace stopped successfully"})
}