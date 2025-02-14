package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/turplespace/portos/internal/database"
	"github.com/turplespace/portos/internal/services/docker"
)

func HandleDeployWorkspace(c echo.Context) error {
	// Get workspace ID from query parameters
	workspaceIDStr := c.QueryParam("workspace_id")
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

	// Start each container
	for _, container := range containers {
		err := docker.StartContainer(container)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to deploy container %s: %v", container.Name, err)})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Workspace deployed successfully"})
}

func HandleRedeployWorkspace(c echo.Context) error {
	// Get workspace ID from query parameters
	workspaceIDStr := c.QueryParam("workspace_id")
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
	workspaceIDStr := c.QueryParam("workspace_id")
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
