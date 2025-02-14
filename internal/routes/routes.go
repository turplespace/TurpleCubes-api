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
	workspaceGroup.GET("/workspaces", handlers.HandleGetWorkspaces)
	workspaceGroup.POST("/create", handlers.HandleCreateWorkspace)
	workspaceGroup.PUT("/edit", handlers.HandleEditWorkspace)
	workspaceGroup.DELETE("/delete", handlers.HandleDeleteWorkspace)
	workspaceGroup.POST("/deploy", handlers.HandleDeployWorkspace)
	workspaceGroup.POST("/redeploy", handlers.HandleRedeployWorkspace)
	workspaceGroup.POST("/stop", handlers.HandleStopWorkspace)

	// Cube routes
	cubeGroup := e.Group("/api/cube")
	cubeGroup.GET("/cubes", handlers.HandleGetCubes)
	cubeGroup.POST("/add", handlers.HandleAddCubes)
	cubeGroup.PUT("/edit", handlers.HandleEditCube)
	cubeGroup.DELETE("/delete", handlers.HandleDeleteCube)
	cubeGroup.GET("", handlers.HandleGetCubeData)
	cubeGroup.POST("/deploy", handlers.HandleDeployCube)
	cubeGroup.POST("/redeploy", handlers.HandleRedeployCube)
	cubeGroup.POST("/stop", handlers.HandleStopCube)
	cubeGroup.POST("/commit", handlers.HandleCommitCube)

	// Proxy route
	e.POST("/api/proxy", handlers.HandlePostProxy)

	// Images route
	e.GET("/api/images", handlers.HandleGetImages)

	// Logs route
	e.GET("/api/logs/stream", handlers.HandleLogStream)
}
