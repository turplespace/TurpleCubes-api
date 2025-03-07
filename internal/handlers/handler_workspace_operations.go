package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/turplespace/portos/internal/database"
	"github.com/turplespace/portos/internal/services/docker"
)

// HandleDeployWorkspace function receives workspace_id in query params and deploys the workspace
func HandleDeployWorkspace(c echo.Context) error {
	// Get workspace ID from query parameters
	workspaceIDStr := c.Param("workspaceID")
	if workspaceIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing workspace ID"})
	}

	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		log.Printf("Failed to convert workspace ID to integer: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid workspace ID"})
	}

	// Get all containers in the workspace
	containers, err := database.ListContainersInWorkspace(workspaceID)
	if err != nil {
		log.Printf("Failed to list containers: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to list containers: %v", err)})
	}

	// Start each container
	for _, container := range containers {
		err := docker.StartContainer(container)
		if err != nil {
			log.Printf("Failed to deploy container %s: %v", container.Name, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to deploy container %s: %v", container.Name, err)})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Workspace deployed successfully"})
}

func HandleRedeployWorkspace(c echo.Context) error {
	// Get workspace ID from query parameters
	workspaceIDStr := c.Param("workspaceID")
	if workspaceIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing workspace ID"})
	}

	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid workspace ID"})
	}

	// Get all containers in the workspace
	containers, err := database.ListContainersInWorkspace(workspaceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to list containers: %v", err)})
	}

	// Redeploy each container
	for _, container := range containers {
		err := docker.RestartContainer(container.Name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to redeploy container %s: %v", container.Name, err)})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Workspace redeployed successfully"})
}

func HandleStopWorkspace(c echo.Context) error {
	// Get workspace ID from query parameters
	workspaceIDStr := c.Param("workspaceID")
	if workspaceIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing workspace ID"})
	}

	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid workspace ID"})
	}

	// Get all containers in the workspace
	containers, err := database.ListContainersInWorkspace(workspaceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to list containers: %v", err)})
	}

	// Stop each container
	for _, container := range containers {
		err := docker.StopContainer(container.Name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to stop container %s: %v", container.Name, err)})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Workspace stopped successfully"})
}
