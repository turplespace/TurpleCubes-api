package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/turplespace/portos/internal/database"
	"github.com/turplespace/portos/internal/routes"
	"github.com/turplespace/portos/internal/services"
	"github.com/turplespace/portos/internal/services/proxy"
)

func main() {
	e := echo.New()

	logService := services.GetLogService()
	routes.SetupRoutes(e)
	database.Init()
	err := proxy.RemoveDataInFolder()
	if err != nil {
		log.Fatalf("Failed to remove data in folder: %v", err)
	}
	err = proxy.CreateFolderIfNotExists()
	if err != nil {
		log.Fatalf("Failed to create folder: %v", err)
	}

	log.Print("Server starting on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	logService.Info("Custom info message")
}
