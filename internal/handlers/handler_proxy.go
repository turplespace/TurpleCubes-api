package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/turplespace/portos/internal/database"
	"github.com/turplespace/portos/internal/models"
)

func HandleGetProxyByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("proxyID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid proxy ID"})
	}

	proxy, err := database.GetProxyByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get proxy"})
	}

	return c.JSON(http.StatusOK, proxy)
}

func HandleGetProxiesByCubeID(c echo.Context) error {
	cubeID, err := strconv.Atoi(c.Param("cubeID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cube ID"})
	}

	proxies, err := database.GetProxiesByCubeID(cubeID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get proxies"})
	}

	return c.JSON(http.StatusOK, proxies)
}

func HandleAddProxy(c echo.Context) error {
	var req models.AddProxyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Check if the domain already exists
	existingID, err := database.GetProxyIDByDomain(req.Domain)
	if err == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Domain already exists",
			"id":      existingID,
		})
	}

	id, err := database.AddProxy(req.CubeID, req.Domain, req.Port, req.Type, req.Default)
	if err != nil {
		log.Printf("Failed to add proxy: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add proxy"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

func HandleEditProxyByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("proxyID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid proxy ID"})
	}
	var req models.EditProxyByIDRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := database.EditProxyByID(id, req.Domain, req.Port, req.Type, req.Default); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to edit proxy"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Proxy updated successfully"})
}

func HandleDeleteProxyByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("proxyID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid proxy ID"})
	}

	if err := database.DeleteProxyByID(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete proxy"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Proxy deleted successfully"})
}

func HandleDeleteProxiesByCubeID(c echo.Context) error {
	cubeID, err := strconv.Atoi(c.Param("cubeID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cube ID"})
	}

	if err := database.DeleteProxiesByCubeID(cubeID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete proxies"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Proxies deleted successfully"})
}
