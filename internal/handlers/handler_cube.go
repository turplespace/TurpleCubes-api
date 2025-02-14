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

/*
HandleGetCubeData function receievs cube_id in query params and return GetCubesByIdResponse model as reponse
*/
func HandleGetCubeData(c echo.Context) error {
	log.Printf("[*] Geting cube data request")

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

	var getCubesByIdResponse models.GetCubesByIdResponse
	cube, err := database.GetCubeData(cubeID)
	if err != nil {
		log.Printf("[*] Database error while fetching cube data: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to get cube data: %v", err)})
	}
	log.Printf("[*] Successfully retrieved cube data for ID: %d", cubeID)

	status, err := docker.GetContainerStatus(cube.Name)
	if err != nil {
		log.Printf("[*] Warning: Unable to get container status: %v", err)
		status = "unknown"
	}

	ipAddress, err := docker.GetContainerIPAddress(cube.Name)
	getCubesByIdResponse.IPAddress = ipAddress
	getCubesByIdResponse.Status = status
	getCubesByIdResponse.ContainerData = cube

	return c.JSON(http.StatusOK, getCubesByIdResponse)
}

/*
HandleGetCubes function returns all the cubes in a workspace
*/
func HandleGetCubes(c echo.Context) error {
	log.Printf("[*] Starting get cubes request ")

	workspaceIDStr := c.QueryParam("workspace_id")
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

/*
HandleAddCubes function receives workspace_id and cubes in request body and add cubes to workspace
*/
func HandleAddCubes(c echo.Context) error {
	log.Printf("[*] Starting add cubes request")

	var req models.AddCubesRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("[*] Error: Invalid request body - %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid request: %v", err)})
	}
	log.Printf("[*] Attempting to add %d cubes to workspace %d", len(req.Cubes), req.WorkspaceID)

	err := database.InsertWorkspaceAndCubes(req.WorkspaceID, req.Cubes)
	if err != nil {
		log.Printf("[*] Database error while inserting cubes: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to insert cubes: %v", err)})
	}

	log.Printf("[*] Successfully added %d cubes to workspace %d", len(req.Cubes), req.WorkspaceID)
	return c.JSON(http.StatusOK, map[string]string{"message": "Cubes added successfully"})
}

/*
HandleEditCube function receives cube_id in query params and restarts the cube
*/
func HandleEditCube(c echo.Context) error {
	log.Printf("[*] Starting edit cube request")

	var req models.EditCubeRequest

	if err := c.Bind(&req); err != nil {
		log.Printf("[*] Error: Invalid request body - %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid request: %v", err)})
	}

	log.Printf("[*] Attempting to update cube ID: %d", req.CubeID)

	err := database.UpdateCube(req.CubeID, req.UpdatedCube)
	if err != nil {
		log.Printf("[*] Database error while updating cube: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to update cube: %v", err)})
	}

	log.Printf("[*] Successfully updated cube ID: %d", req.CubeID)
	return c.JSON(http.StatusOK, map[string]string{"message": "Cube updated successfully"})
}

/*
HandleDeleteCube function receives cube_id in query params and restarts the cube
*/
func HandleDeleteCube(c echo.Context) error {
	log.Printf("[*] Starting delete cube request")

	cubeIDStr := c.QueryParam("cube_id")
	if cubeIDStr == "" {
		log.Printf("[*] Error: Missing cube ID in request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing cube ID"})
	}

	cubeID, err := strconv.Atoi(cubeIDStr)
	if err != nil {
		log.Printf("[*] Error: Invalid cube ID format: %s - %v", cubeIDStr, err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cube ID"})
	}

	cube, err := database.GetCubeData(cubeID)
	if err != nil {
		log.Printf("[*] Database error while fetching cube data: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to get cube data: %v", err)})
	}
	log.Printf("[*] Retrieved cube data for deletion, container name: %s", cube.Name)

	err = docker.StopContainer(cube.Name)
	if err != nil {
		log.Printf("[*] Warning: Error stopping container %s: %v", cube.Name, err)
	}

	err = database.DeleteCube(cubeID)
	if err != nil {
		log.Printf("[*] Database error while deleting cube: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to delete cube: %v", err)})
	}

	log.Printf("[*] Successfully deleted cube ID: %d", cubeID)
	return c.JSON(http.StatusOK, map[string]string{"message": "Cube deleted successfully"})
}
