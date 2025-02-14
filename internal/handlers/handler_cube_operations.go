package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/turplespace/portos/internal/database"
	"github.com/turplespace/portos/internal/models"
	"github.com/turplespace/portos/internal/services/docker"
	"github.com/turplespace/portos/internal/services/repositories"
)

func HandleDeployCube(c echo.Context) error {
	log.Printf("[*] Starting deployment request at %s", time.Now().UTC().Format(time.RFC3339))

	// Get cube ID from query parameters
	cubeIDStr := c.QueryParam("cube_id")
	if cubeIDStr == "" {
		log.Printf("[*] Error: No cube ID provided in request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing cube ID"})
	}

	cubeID, err := strconv.Atoi(cubeIDStr)
	if err != nil {
		log.Printf("[*] Error: Invalid cube ID format: %s - %v", cubeIDStr, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cube ID"})
	}
	log.Printf("[*] Processing deployment for cube ID: %d", cubeID)

	// Get cube data from the database
	container, err := database.GetCubeData(cubeID)
	if err != nil {
		log.Printf("[*] Database error while fetching cube data: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to get cube data: %v", err)})
	}
	log.Printf("[*] Successfully retrieved cube data for ID: %d", cubeID)

	// Start the container using the retrieved cube data
	err = docker.StartContainer(*container)
	if err != nil {
		log.Printf("[*] Docker error while starting container: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to deploy cube: %v", err)})
	}
	log.Printf("[*] Successfully started container for cube ID: %d", cubeID)

	return c.JSON(http.StatusOK, map[string]string{"message": "Cube deployed successfully"})
}

func HandleRedeployCube(c echo.Context) error {
	log.Printf("[*] Starting redeployment request at %s", time.Now().UTC().Format(time.RFC3339))

	cubeIDStr := c.QueryParam("cube_id")
	if cubeIDStr == "" {
		log.Printf("[*] Error: No cube ID provided in request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing cube ID"})
	}

	cubeID, err := strconv.Atoi(cubeIDStr)
	if err != nil {
		log.Printf("[*] Error: Invalid cube ID format: %s - %v", cubeIDStr, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cube ID"})
	}
	log.Printf("[*] Processing redeployment for cube ID: %d", cubeID)

	container, err := database.GetCubeData(cubeID)
	if err != nil {
		log.Printf("[*] Database error while fetching cube data: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to get cube data: %v", err)})
	}
	log.Printf("[*] Successfully retrieved cube data for ID: %d", cubeID)

	err = docker.RestartContainer(container.Name)
	if err != nil {
		log.Printf("[*] Docker error while restarting container: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to redeploy cube: %v", err)})
	}
	log.Printf("[*] Successfully restarted container for cube ID: %d", cubeID)

	return c.JSON(http.StatusOK, map[string]string{"message": "Cube redeployed successfully"})
}

func HandleStopCube(c echo.Context) error {
	log.Printf("[*] Starting stop request at %s", time.Now().UTC().Format(time.RFC3339))

	cubeIDStr := c.QueryParam("cube_id")
	if cubeIDStr == "" {
		log.Printf("[*] Error: No cube ID provided in request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing cube ID"})
	}

	cubeID, err := strconv.Atoi(cubeIDStr)
	if err != nil {
		log.Printf("[*] Error: Invalid cube ID format: %s - %v", cubeIDStr, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cube ID"})
	}
	log.Printf("[*] Processing stop request for cube ID: %d", cubeID)

	container, err := database.GetCubeData(cubeID)
	if err != nil {
		log.Printf("[*] Database error while fetching cube data: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to get cube data: %v", err)})
	}
	log.Printf("[*] Successfully retrieved cube data for ID: %d", cubeID)

	err = docker.StopContainer(container.Name)
	if err != nil {
		log.Printf("[*] Docker error while stopping container: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to stop cube: %v", err)})
	}
	log.Printf("[*] Successfully stopped container for cube ID: %d", cubeID)

	return c.JSON(http.StatusOK, map[string]string{"message": "Cube stopped successfully"})
}

func HandleCommitCube(c echo.Context) error {
	log.Printf("[*] Starting commit request at %s", time.Now().UTC().Format(time.RFC3339))

	cubeIDStr := c.QueryParam("cube_id")
	if cubeIDStr == "" {
		log.Printf("[*] Error: No cube ID provided in request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing cube ID"})
	}

	cubeID, err := strconv.Atoi(cubeIDStr)
	if err != nil {
		log.Printf("[*] Error: Invalid cube ID format: %s - %v", cubeIDStr, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cube ID"})
	}

	newImage := c.QueryParam("image")
	tag := c.QueryParam("tag")

	if newImage == "" || tag == "" {
		log.Printf("[*] Error: Missing image name or tag in request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing new image or tag"})
	}
	log.Printf("[*] Processing commit for cube ID: %d with image: %s and tag: %s", cubeID, newImage, tag)

	container, err := database.GetCubeData(cubeID)
	if err != nil {
		log.Printf("[*] Database error while fetching cube data: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to get cube data: %v", err)})
	}
	log.Printf("[*] Successfully retrieved cube data for ID: %d", cubeID)

	err = docker.CommitContainer(container.Name, newImage, tag)
	if err != nil {
		log.Printf("[*] Docker error while committing container: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to commit cube: %v", err)})
	}
	log.Printf("[*] Successfully committed container for cube ID: %d", cubeID)
	repositories.AppendImages(models.Image{Image: newImage, Tag: tag, PulledOn: time.Now().UTC().Format(time.RFC3339)})

	return c.JSON(http.StatusOK, map[string]string{"message": "Cube committed successfully"})
}
