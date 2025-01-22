package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/turplespace/portos/internal/database"
	"github.com/turplespace/portos/internal/models"
	"github.com/turplespace/portos/internal/services/docker"
)

func HandleGetCubeData(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GET-CUBE-DATA] Starting get cube data request")

	if r.Method != http.MethodGet {
		log.Printf("[GET-CUBE-DATA] Error: Invalid method %s used instead of GET",
			"", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cubeIDStr := r.URL.Query().Get("cube_id")
	if cubeIDStr == "" {
		log.Printf("[GET-CUBE-DATA] Error: No cube ID provided in request", "")
		http.Error(w, "Missing cube ID", http.StatusBadRequest)
		return
	}

	cubeID, err := strconv.Atoi(cubeIDStr)
	if err != nil {
		log.Printf("[GET-CUBE-DATA] Error: Invalid cube ID format: %s - %v",
			"", cubeIDStr, err)
		http.Error(w, "Invalid cube ID", http.StatusBadRequest)
		return
	}

	var getCubesByIdResponse models.GetCubesByIdResponse
	cube, err := database.GetCubeData(cubeID)
	if err != nil {
		log.Printf("[GET-CUBE-DATA] Database error while fetching cube data: %v",
			"", err)
		http.Error(w, fmt.Sprintf("Failed to get cube data: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("[GET-CUBE-DATA] Successfully retrieved cube data for ID: %d",
		"", cubeID)

	
	status, err := docker.GetContainerStatus(cube.Name)
	if err != nil {
		log.Printf("[GET-CUBE-DATA] Warning: Unable to get container status: %v",
			"", err)
		status = "unknown"
	}

	ipAddress, err:= docker.GetContainerIPAddress(cube.Name)
	getCubesByIdResponse.IPAddress = ipAddress
	getCubesByIdResponse.Status = status
	getCubesByIdResponse.ContainerData = cube

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getCubesByIdResponse)
	log.Printf("[GET-CUBE-DATA] Successfully sent cube data response for ID: %d",
		"", cubeID)
}

func HandleGetCubes(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GET-CUBES] Starting get cubes request",
		time.Now().UTC().Format(time.RFC3339), "")

	if r.Method != http.MethodGet {
		log.Printf("[GET-CUBES] Error: Invalid method %s used instead of GET",
			"", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workspaceIDStr := r.URL.Query().Get("workspace_id")
	if workspaceIDStr == "" {
		log.Printf("[GET-CUBES] Error: Missing workspace ID in request", "")
		http.Error(w, "Missing workspace ID", http.StatusBadRequest)
		return
	}

	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		log.Printf("[GET-CUBES] Error: Invalid workspace ID format: %s - %v",
			"", workspaceIDStr, err)
		http.Error(w, "Invalid workspace ID", http.StatusBadRequest)
		return
	}

	cubes, err := database.ListCubes(workspaceID)
	if err != nil {
		log.Printf("[GET-CUBES] Database error while fetching cubes: %v",
			"", err)
		http.Error(w, fmt.Sprintf("Failed to get cubes: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("[GET-CUBES] Successfully retrieved %d cubes for workspace ID: %d",
		"", len(cubes), workspaceID)

	var cubesResponse []models.GetCubesResponse
	for _, cube := range cubes {
		status, err := docker.GetContainerStatus(cube.Name)
		if err != nil {
			log.Printf("[GET-CUBES] Warning: Unable to get status for container %s: %v",
				"", cube.Name, err)
			status = "unknown"
		}
		ipAddress, err:= docker.GetContainerIPAddress(cube.Name)
		if err != nil {
			log.Printf("[GET-CUBES] Warning: Unable to get IP address for container %s: %v",
				"", cube.Name, err)
			ipAddress = "unknown"
		}
		cubesResponse = append(cubesResponse, models.GetCubesResponse{
			ContainerID:   cube.ID,
			Image:         cube.Image,
			ContainerName: cube.Name,
			IPAddress:         ipAddress,
			Status:        status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cubesResponse)
	log.Printf("[GET-CUBES] Successfully sent response with %d cubes",
		"", len(cubesResponse))
}

func HandleAddCubes(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ADD-CUBES] Starting add cubes request",
		time.Now().UTC().Format(time.RFC3339), "")

	if r.Method != http.MethodPost {
		log.Printf("[ADD-CUBES] Error: Invalid method %s used instead of POST",
			"", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		WorkspaceID int                `json:"workspace_id"`
		Cubes       []models.Container `json:"cubes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[ADD-CUBES] Error: Invalid request body - %v",
			"", err)
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}
	log.Printf("[ADD-CUBES] Attempting to add %d cubes to workspace %d",
		"", len(req.Cubes), req.WorkspaceID)

	err := database.InsertWorkspaceAndCubes(req.WorkspaceID, req.Cubes)
	if err != nil {
		log.Printf("[ADD-CUBES] Database error while inserting cubes: %v",
			"", err)
		http.Error(w, fmt.Sprintf("Failed to insert cubes: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("[ADD-CUBES] Successfully added %d cubes to workspace %d",
		"", len(req.Cubes), req.WorkspaceID)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cubes added successfully"})
}

func HandleEditCube(w http.ResponseWriter, r *http.Request) {
	log.Printf("[EDIT-CUBE] Starting edit cube request",
		time.Now().UTC().Format(time.RFC3339), "")

	if r.Method != http.MethodPut {
		log.Printf("[EDIT-CUBE] Error: Invalid method %s used instead of PUT",
			"", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CubeID      int              `json:"cube_id"`
		UpdatedCube models.Container `json:"updated_cube"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[EDIT-CUBE] Error: Invalid request body - %v",
			"", err)
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("[EDIT-CUBE] Attempting to update cube ID: %d",
		"", req.CubeID)

	err := database.UpdateCube(req.CubeID, req.UpdatedCube)
	if err != nil {
		log.Printf("[EDIT-CUBE] Database error while updating cube: %v",
			"", err)
		http.Error(w, fmt.Sprintf("Failed to update cube: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("[EDIT-CUBE] Successfully updated cube ID: %d",
		"", req.CubeID)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cube updated successfully"})
}

func HandleDeleteCube(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DELETE-CUBE] Starting delete cube request",
		time.Now().UTC().Format(time.RFC3339), "")

	if r.Method != http.MethodDelete {
		log.Printf("[DELETE-CUBE] Error: Invalid method %s used instead of DELETE",
			"", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cubeIDStr := r.URL.Query().Get("cube_id")
	if cubeIDStr == "" {
		log.Printf("[DELETE-CUBE] Error: Missing cube ID in request", "")
		http.Error(w, "Missing cube ID", http.StatusBadRequest)
		return
	}

	cubeID, err := strconv.Atoi(cubeIDStr)
	if err != nil {
		log.Printf("[DELETE-CUBE] Error: Invalid cube ID format: %s - %v",
			"", cubeIDStr, err)
		http.Error(w, "Invalid cube ID", http.StatusBadRequest)
		return
	}

	cube, err := database.GetCubeData(cubeID)
	if err != nil {
		log.Printf("[DELETE-CUBE] Database error while fetching cube data: %v",
			"", err)
		http.Error(w, fmt.Sprintf("Failed to get cube data: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("[DELETE-CUBE] Retrieved cube data for deletion, container name: %s",
		"", cube.Name)

	err = docker.StopContainer(cube.Name)
	if err != nil {
		log.Printf("[DELETE-CUBE] Warning: Error stopping container %s: %v",
			"", cube.Name, err)
	}

	err = database.DeleteCube(cubeID)
	if err != nil {
		log.Printf("[DELETE-CUBE] Database error while deleting cube: %v",
			"", err)
		http.Error(w, fmt.Sprintf("Failed to delete cube: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("[DELETE-CUBE] Successfully deleted cube ID: %d",
		"", cubeID)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cube deleted successfully"})
}
