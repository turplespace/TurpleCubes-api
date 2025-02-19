package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/turplespace/portos/internal/database"
	"github.com/turplespace/portos/internal/models"
	"github.com/turplespace/portos/internal/services/docker"
)

// HandleGetWorkspaces handles the HTTP request to get the list of workspaces
func HandleGetWorkspaces(c echo.Context) error {
	log.Println("[*] Starting get workspaces request")

	// Retrieve the list of workspaces
	workspaces, err := database.GetWorkspaces()
	if err != nil {
		log.Printf("[*] Error: Failed to get workspaces: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to get workspaces: %v", err)})
	}

	// Get the total counts
	totalWorkspaces, err := database.CountWorkspaces()
	if err != nil {
		log.Printf("[*] Error: Failed to count workspaces: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to count workspaces: %v", err)})
	}

	totalCubes, err := database.CountCubes()
	if err != nil {
		log.Printf("[*] Error: Failed to count cubes: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to count cubes: %v", err)})
	}

	totalRunningCubes, err := docker.CountContainersByLabel("service", "turplespace")
	if err != nil {
		log.Printf("[*] Error: Failed to count running containers: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to count running containers: %v", err)})
	}

	// Create a slice to hold workspaces with container counts
	var workspacesWithCounts []models.WorkspaceWithContainerCounts

	// For each workspace, count the number of total and running containers
	for _, workspace := range workspaces {
		totalCount, err := database.CountContainersByWorkspaceID(workspace.ID)
		if err != nil {
			log.Printf("[*] Error: Failed to count containers for workspace %d: %v", workspace.ID, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to count containers for workspace %d: %v", workspace.ID, err)})
		}

		runningCount, err := docker.CountContainersByLabel("workspace_id", fmt.Sprintf("%d", workspace.ID))
		if err != nil {
			log.Printf("[*] Error: Failed to count running containers for workspace %d: %v", workspace.ID, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to count running containers for workspace %d: %v", workspace.ID, err)})
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
	response := models.WorkspaceResponse{
		TotalWorkspaces:   totalWorkspaces,
		TotalCubes:        totalCubes,
		TotalRunningCubes: totalRunningCubes,
		Workspaces:        workspacesWithCounts,
	}

	// Encode the response as JSON
	return c.JSON(http.StatusOK, response)
}

// HandleCreateWorkspace handles the HTTP request to create a new workspace
func HandleCreateWorkspace(c echo.Context) error {
	log.Println("[*] Starting create workspace request")

	var req models.CreateWorkspaceRequest

	if err := c.Bind(&req); err != nil {
		log.Printf("[*] Error: Invalid request body - %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid request: %v", err)})
	}

	id, err := database.CreateWorkspace(req.Name, req.Desc)
	if err != nil {
		log.Printf("[*] Error: Failed to create workspace: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create workspace: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]int{"id": int(id)})
}

/*
HandleEditWorkspace handles the HTTP request to edit an existing workspace
request body should contain id name and desc
*/
func HandleEditWorkspace(c echo.Context) error {
	log.Println("[*] Starting edit workspace request")
	idStr := c.Param("workspaceID")
	if idStr == "" {
		log.Println("[*] Error: Missing workspace ID")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing workspace ID"})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[*] Error: Invalid workspace ID format: %s - %v", idStr, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid workspace ID"})
	}
	var req models.EditWorkspaceRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("[*] Error: Invalid request body - %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid request: %v", err)})
	}

	err = database.EditWorkspace(id, req.Name, req.Desc)
	if err != nil {
		log.Printf("[*] Error: Failed to edit workspace: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to edit workspace: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Workspace updated successfully"})
}

/*
HandleDeleteWorkspace handles the HTTP request to delete an existing workspace
request quesy param should contain id
*/
func HandleDeleteWorkspace(c echo.Context) error {
	log.Println("[*] Starting delete workspace request")

	idStr := c.Param("workspaceID")
	if idStr == "" {
		log.Println("[*] Error: Missing workspace ID")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing workspace ID"})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[*] Error: Invalid workspace ID format: %s - %v", idStr, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid workspace ID"})
	}

	// Stopping Cubes
	cubes, err := database.ListCubes(id)
	if err != nil {
		log.Printf("[*] Error: Failed to get cubes for workspace ID %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to get cubes: %v", err)})
	}

	for _, cube := range cubes {
		err = docker.StopContainer(cube.Name)
		if err != nil {
			log.Printf("[*] Warning: Error stopping container %s: %v", cube.Name, err)
		}
	}

	// Deleting the Cubes from DB
	err = database.DeleteContainersByWorkspaceID(id)
	if err != nil {
		log.Printf("[*] Error: Failed to delete containers for workspace ID %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to delete containers: %v", err)})
	}

	// Deleting Workspace from the DB
	err = database.DeleteWorkspace(id)
	if err != nil {
		log.Printf("[*] Error: Failed to delete workspace ID %d: %v", id, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to delete workspace: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Workspace deleted successfully"})
}

/*
HandleGetCubes function returns all the cubes in a workspace
*/
func HandleGetWorkspaceData(c echo.Context) error {
	log.Printf("[*] Starting get cubes request ")

	workspaceIDStr := c.Param("workspaceID")
	if workspaceIDStr == "" {
		log.Printf("[*] Error: Missing workspace ID in request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing workspace ID"})
	}

	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		log.Printf("[*] Error: Invalid workspace ID format: %s - %v", workspaceIDStr, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid workspace ID"})
	}

	cubes, err := database.ListCubes(workspaceID)
	if err != nil {
		log.Printf("[*] Database error while fetching cubes: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to get cubes: %v", err)})
	}
	log.Printf("[*] Successfully retrieved %d cubes for workspace ID: %d", len(cubes), workspaceID)

	var cubesResponse []models.GetCubesResponse
	for _, cube := range cubes {
		status, err := docker.GetContainerStatus(cube.Name)
		if err != nil {
			log.Printf("[*] Warning: Unable to get status for container %s: %v", cube.Name, err)
			status = "unknown"
		}
		ipAddress, err := docker.GetContainerIPAddress(cube.Name)
		if err != nil {
			log.Printf("[*] Warning: Unable to get IP address for container %s: %v", cube.Name, err)
			ipAddress = "unknown"
		}
		cubesResponse = append(cubesResponse, models.GetCubesResponse{
			ContainerID:   cube.ID,
			Image:         cube.Image,
			ContainerName: cube.Name,
			IPAddress:     ipAddress,
			Status:        status,
		})
	}

	return c.JSON(http.StatusOK, cubesResponse)
}
