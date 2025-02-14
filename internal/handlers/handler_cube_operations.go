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
	"github.com/turplespace/portos/internal/services/repositories"
)

func HandleDeployCube(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DEPLOY-CUBE] Starting deployment request at %s", time.Now().UTC().Format(time.RFC3339))

	if r.Method != http.MethodPost {
		log.Printf("[DEPLOY-CUBE] Error: Invalid method %s used instead of POST", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get cube ID from query parameters
	cubeIDStr := r.URL.Query().Get("cube_id")
	if cubeIDStr == "" {
		log.Printf("[DEPLOY-CUBE] Error: No cube ID provided in request")
		http.Error(w, "Missing cube ID", http.StatusBadRequest)
		return
	}

	cubeID, err := strconv.Atoi(cubeIDStr)
	if err != nil {
		log.Printf("[DEPLOY-CUBE] Error: Invalid cube ID format: %s - %v", cubeIDStr, err)
		http.Error(w, "Invalid cube ID", http.StatusBadRequest)
		return
	}
	log.Printf("[DEPLOY-CUBE] Processing deployment for cube ID: %d", cubeID)

	// Get cube data from the database
	container, err := database.GetCubeData(cubeID)
	if err != nil {
		log.Printf("[DEPLOY-CUBE] Database error while fetching cube data: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get cube data: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("[DEPLOY-CUBE] Successfully retrieved cube data for ID: %d", cubeID)

	// Start the container using the retrieved cube data
	err = docker.StartContainer(*container)
	if err != nil {
		log.Printf("[DEPLOY-CUBE] Docker error while starting container: %v", err)
		http.Error(w, fmt.Sprintf("Failed to deploy cube: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("[DEPLOY-CUBE] Successfully started container for cube ID: %d", cubeID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Cube deployed successfully"})
	log.Printf("[DEPLOY-CUBE] Deployment completed successfully for cube ID: %d", cubeID)
}

func HandleRedeployCube(w http.ResponseWriter, r *http.Request) {
	log.Printf("[REDEPLOY-CUBE] Starting redeployment request at %s", time.Now().UTC().Format(time.RFC3339))

	if r.Method != http.MethodPost {
		log.Printf("[REDEPLOY-CUBE] Error: Invalid method %s used instead of POST", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cubeIDStr := r.URL.Query().Get("cube_id")
	if cubeIDStr == "" {
		log.Printf("[REDEPLOY-CUBE] Error: No cube ID provided in request")
		http.Error(w, "Missing cube ID", http.StatusBadRequest)
		return
	}

	cubeID, err := strconv.Atoi(cubeIDStr)
	if err != nil {
		log.Printf("[REDEPLOY-CUBE] Error: Invalid cube ID format: %s - %v", cubeIDStr, err)
		http.Error(w, "Invalid cube ID", http.StatusBadRequest)
		return
	}
	log.Printf("[REDEPLOY-CUBE] Processing redeployment for cube ID: %d", cubeID)

	container, err := database.GetCubeData(cubeID)
	if err != nil {
		log.Printf("[REDEPLOY-CUBE] Database error while fetching cube data: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get cube data: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("[REDEPLOY-CUBE] Successfully retrieved cube data for ID: %d", cubeID)

	err = docker.RestartContainer(container.Name)
	if err != nil {
		log.Printf("[REDEPLOY-CUBE] Docker error while restarting container: %v", err)
		http.Error(w, fmt.Sprintf("Failed to redeploy cube: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("[REDEPLOY-CUBE] Successfully restarted container for cube ID: %d", cubeID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Cube redeployed successfully"})
	log.Printf("[REDEPLOY-CUBE] Redeployment completed successfully for cube ID: %d", cubeID)
}

func HandleStopCube(w http.ResponseWriter, r *http.Request) {
	log.Printf("[STOP-CUBE] Starting stop request at %s", time.Now().UTC().Format(time.RFC3339))

	if r.Method != http.MethodPost {
		log.Printf("[STOP-CUBE] Error: Invalid method %s used instead of POST", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cubeIDStr := r.URL.Query().Get("cube_id")
	if cubeIDStr == "" {
		log.Printf("[STOP-CUBE] Error: No cube ID provided in request")
		http.Error(w, "Missing cube ID", http.StatusBadRequest)
		return
	}

	cubeID, err := strconv.Atoi(cubeIDStr)
	if err != nil {
		log.Printf("[STOP-CUBE] Error: Invalid cube ID format: %s - %v", cubeIDStr, err)
		http.Error(w, "Invalid cube ID", http.StatusBadRequest)
		return
	}
	log.Printf("[STOP-CUBE] Processing stop request for cube ID: %d", cubeID)

	container, err := database.GetCubeData(cubeID)
	if err != nil {
		log.Printf("[STOP-CUBE] Database error while fetching cube data: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get cube data: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("[STOP-CUBE] Successfully retrieved cube data for ID: %d", cubeID)

	err = docker.StopContainer(container.Name)
	if err != nil {
		log.Printf("[STOP-CUBE] Docker error while stopping container: %v", err)
		http.Error(w, fmt.Sprintf("Failed to stop cube: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("[STOP-CUBE] Successfully stopped container for cube ID: %d", cubeID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Cube stopped successfully"})
	log.Printf("[STOP-CUBE] Stop operation completed successfully for cube ID: %d", cubeID)
}

func HandleCommitCube(w http.ResponseWriter, r *http.Request) {
	log.Printf("[*] Starting commit request at %s", time.Now().UTC().Format(time.RFC3339))

	if r.Method != http.MethodPost {
		log.Printf("[*] Error: Invalid method %s used instead of POST", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cubeIDStr := r.URL.Query().Get("cube_id")
	if cubeIDStr == "" {
		log.Printf("[*] Error: No cube ID provided in request")
		http.Error(w, "Missing cube ID", http.StatusBadRequest)
		return
	}

	cubeID, err := strconv.Atoi(cubeIDStr)
	if err != nil {
		log.Printf("[*] Error: Invalid cube ID format: %s - %v", cubeIDStr, err)
		http.Error(w, "Invalid cube ID", http.StatusBadRequest)
		return
	}

	newImage := r.URL.Query().Get("image")
	tag := r.URL.Query().Get("tag")

	if newImage == "" || tag == "" {
		log.Printf("[*] Error: Missing image name or tag in request")
		http.Error(w, "Missing new image or tag", http.StatusBadRequest)
		return
	}
	log.Printf("[*] Processing commit for cube ID: %d with image: %s and tag: %s", cubeID, newImage, tag)

	container, err := database.GetCubeData(cubeID)
	if err != nil {
		log.Printf("[*] Database error while fetching cube data: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get cube data: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("[*] Successfully retrieved cube data for ID: %d", cubeID)

	err = docker.CommitContainer(container.Name, newImage, tag)
	if err != nil {
		log.Printf("[*] Docker error while committing container: %v", err)
		http.Error(w, fmt.Sprintf("Failed to commit cube: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("[*] Successfully committed container for cube ID: %d", cubeID)
	repositories.AppendImages(models.Image{Image: newImage, Tag: tag, PulledOn: time.Now().UTC().Format(time.RFC3339)})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Cube committed successfully"})
	log.Printf("[*] Commit operation completed successfully for cube ID: %d", cubeID)
}
