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
	e.GET("/api/workspaces", handlers.HandleGetWorkspaces)
	e.POST("/api/workspace/create", handlers.HandleCreateWorkspace)
	e.PUT("/api/workspace/edit", handlers.HandleEditWorkspace)
	e.DELETE("/api/workspace/delete", handlers.HandleDeleteWorkspace)

	// Workspace operations routes
	e.POST("/api/workspace/deploy", handlers.HandleDeployWorkspace)
	e.POST("/api/workspace/redeploy", handlers.HandleRedeployWorkspace)
	e.POST("/api/workspace/stop", handlers.HandleStopWorkspace)

	// Cube routes
	e.GET("/api/cubes", handlers.HandleGetCubes)
	e.POST("/api/cube/add", handlers.HandleAddCubes)
	e.PUT("/api/cube/edit", handlers.HandleEditCube)
	e.DELETE("/api/cube/delete", handlers.HandleDeleteCube)

	// Cube operations routes
	e.GET("/api/cube", handlers.HandleGetCubeData)
	e.POST("/api/cube/deploy", handlers.HandleDeployCube)
	e.POST("/api/cube/redeploy", handlers.HandleRedeployCube)
	e.POST("/api/cube/stop", handlers.HandleStopCube)
	e.POST("/api/cube/commit", handlers.HandleCommitCube)

	// Proxy route
	e.POST("/api/proxy", handlers.HandlePostProxy)

	// Images route
	e.GET("/api/images", handlers.HandleGetImages)

	// Logs route
	e.GET("/api/logs/stream", handlers.HandleLogStream)
}
