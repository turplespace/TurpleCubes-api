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

// WorkspaceResponse holds the total counts and the list of workspaces with container counts
type WorkspaceResponse struct {
	TotalWorkspaces   int                                   `json:"total_workspaces"`
	TotalCubes        int                                   `json:"total_cubes"`
	TotalRunningCubes int                                   `json:"total_running_cubes"`
	Workspaces        []models.WorkspaceWithContainerCounts `json:"workspaces"`
}

// HandleGetWorkspaces handles the HTTP request to get the list of workspaces with container counts
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
	response := WorkspaceResponse{
		TotalWorkspaces:   totalWorkspaces,
		TotalCubes:        totalCubes,
		TotalRunningCubes: totalRunningCubes,
		Workspaces:        workspacesWithCounts,
	}

	// Encode the response as JSON
	return c.JSON(http.StatusOK, response)
}

func HandleCreateWorkspace(c echo.Context) error {
	log.Println("[*] Starting create workspace request")

	var req struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}

	if err := c.Bind(&req); err != nil {
		log.Printf("[*] Error: Invalid request body - %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid request: %v", err)})
	}

	err := database.CreateWorkspace(req.Name, req.Desc)
	if err != nil {
		log.Printf("[*] Error: Failed to create workspace: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create workspace: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Workspace created successfully"})
}

func HandleEditWorkspace(c echo.Context) error {
	log.Println("[*] Starting edit workspace request")

	var req struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Desc string `json:"desc"`
	}

	if err := c.Bind(&req); err != nil {
		log.Printf("[*] Error: Invalid request body - %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid request: %v", err)})
	}

	err := database.EditWorkspace(req.ID, req.Name, req.Desc)
	if err != nil {
		log.Printf("[*] Error: Failed to edit workspace: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to edit workspace: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Workspace updated successfully"})
}

func HandleDeleteWorkspace(c echo.Context) error {
	log.Println("[*] Starting delete workspace request")

	idStr := c.QueryParam("id")
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
