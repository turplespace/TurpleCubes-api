package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/turplespace/portos/internal/database"
	"github.com/turplespace/portos/internal/services/docker"
	"github.com/turplespace/portos/internal/services/proxy"
)

// HandlePostProxy function receives ID in request body, fetches the proxy data from the database, and generates a proxy configuration
func HandlePostStartProxy(c echo.Context) error {

	proxyIDStr := c.Param("proxyID")
	if proxyIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing proxy ID"})
	}

	proxyID, err := strconv.Atoi(proxyIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid proxy ID"})
	}
	proxyData, err := database.GetProxyByID(proxyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch proxy data"})
	}

	// Get cube data from the database
	container, err := database.GetCubeData(proxyData.CubeID)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to get cube data: %v", err)})
	}
	ipAddress, err := docker.GetContainerIPAddress(container.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to get container IP address: %v", err)})
	}
	if err := proxy.GenerateNginxProxyConfig(ipAddress, proxyData.Port, proxyData.Domain); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate proxy config"})
	}

	if err := proxy.RestartNginxService(); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to restart Nginx service"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Proxy configuration generated successfully"})
}
