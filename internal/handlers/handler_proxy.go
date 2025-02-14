package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/turplespace/portos/internal/models"
	"github.com/turplespace/portos/internal/services/proxy"
)

// HandlePostProxy function receives IP, Port and Subdomain in request body and generates a proxy configuration
func HandlePostProxy(c echo.Context) error {
	var req models.ProxyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if req.IP == "" || req.Port == 0 || req.Subdomain == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
	}

	if err := proxy.GenerateNginxProxyConfig(req.IP, req.Port, req.Subdomain); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate proxy config"})
	}

	if err := proxy.RestartNginxService(); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to restart Nginx service"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Proxy configuration generated successfully"})
}
