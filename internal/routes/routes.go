package routes

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/turplespace/portos/internal/handlers"
)

func SetupRoutes(e *echo.Echo) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	fmt.Println(ex)
	path := fmt.Sprintf("%s_web", ex)

	// Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	}))

	// Static files
	e.Static("/", path)

	// Workspace routes
	workspaceGroup := e.Group("/api/workspace")
	workspaceGroup.GET("", handlers.HandleGetWorkspaces)
	workspaceGroup.POST("", handlers.HandleCreateWorkspace)
	workspaceGroup.PUT("/:workspaceID", handlers.HandleEditWorkspace)
	workspaceGroup.DELETE("/:workspaceID", handlers.HandleDeleteWorkspace)
	workspaceGroup.GET("/:workspaceID", handlers.HandleGetWorkspaceData)
	workspaceGroup.POST("/:workspaceID/deploy", handlers.HandleDeployWorkspace)
	workspaceGroup.POST("/:workspaceID/redeploy", handlers.HandleRedeployWorkspace)
	workspaceGroup.POST("/:workspaceID/stop", handlers.HandleStopWorkspace)

	// Cube routes
	cubeGroup := e.Group("/api/cube")

	cubeGroup.POST("", handlers.HandleAddCubes)
	cubeGroup.PUT("/:cubeID", handlers.HandleEditCube)
	cubeGroup.DELETE("/:cubeID", handlers.HandleDeleteCube)
	cubeGroup.GET("/:cubeID", handlers.HandleGetCubeData)
	cubeGroup.POST("/:cubeID/deploy", handlers.HandleDeployCube)
	cubeGroup.POST("/:cubeID/redeploy", handlers.HandleRedeployCube)
	cubeGroup.POST("/:cubeID/stop", handlers.HandleStopCube)
	cubeGroup.POST("/:cubeID/commit", handlers.HandleCommitCube)

	// Proxy route
	proxyGroup := e.Group("/api/proxy")
	proxyGroup.POST("", handlers.HandleAddProxy)
	proxyGroup.GET("/:proxyID", handlers.HandleGetProxyByID)
	proxyGroup.PUT("/:proxyID", handlers.HandleEditProxyByID)
	proxyGroup.DELETE("/:proxyID", handlers.HandleDeleteProxyByID)

	proxyGroup.POST("/:proxyID/deploy", handlers.HandlePostStartProxy)

	proxyGroup.GET("/by-cube/:cubeID", handlers.HandleGetProxiesByCubeID)
	proxyGroup.DELETE("/by-cube/:cubeID", handlers.HandleDeleteProxiesByCubeID)
	// Images route
	e.GET("/api/repo/local", handlers.HandleGetImages)

	// Logs route
	e.GET("/api/logs/stream", handlers.HandleLogStream)
}
