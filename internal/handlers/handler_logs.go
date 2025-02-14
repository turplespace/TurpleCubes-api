package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/turplespace/portos/internal/services"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// HandleLogStream handles WebSocket connections for real-time log streaming
func HandleLogStream(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not upgrade connection"})
	}
	defer conn.Close()

	logService := services.GetLogService()
	msgChan := logService.Subscribe()
	defer logService.Unsubscribe(msgChan)

	// Handle WebSocket connection closure
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket client disconnected: %v", err)
				return
			}
		}
	}()

	// Keep connection alive and send messages
	for message := range msgChan {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Error writing to WebSocket: %v", err)
			break
		}
	}

	return nil
}
