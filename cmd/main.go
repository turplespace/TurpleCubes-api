package main

import (
	"log"
	"net/http"

	"github.com/turplespace/portos/internal/database"
	"github.com/turplespace/portos/internal/routes"
	"github.com/turplespace/portos/internal/services"
)

func main() {
	logService := services.GetLogService()
	routes.SetupRoutes()
	database.Init()
	log.Print("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	logService.Info("Custom info message")
}
