package main

import (
	"log"
	"net/http"

	"github.com/turplespace/portos/internal/database"
	"github.com/turplespace/portos/internal/routes"
	"github.com/turplespace/portos/internal/services"
	"github.com/turplespace/portos/internal/services/proxy"
)

func main() {
	logService := services.GetLogService()
	routes.SetupRoutes()
	database.Init()
	err := proxy.RemoveDataInFolder()
	if err != nil {
		log.Fatalf("Failed to remove data in folder: %v", err)
	}
	err = proxy.CreateFolderIfNotExists()
	if err != nil {
		log.Fatalf("Failed to create folder: %v", err)
	}
	err = proxy.RunNginxDockerContainer()
	if err != nil {
		log.Fatalf("Failed to restart Nginx service: %v", err)
	}
	log.Print("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	logService.Info("Custom info message")
}
