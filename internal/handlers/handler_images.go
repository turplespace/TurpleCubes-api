package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/turplespace/portos/internal/models"
	"github.com/turplespace/portos/internal/services/repositories"
)

type ImagesResponse struct {
	CustomImages      []models.Image `json:"custom_images"`
	TotalCustomImages int            `json:"total_custom_images"`
}

func HandleGetImages(c echo.Context) error {
	log.Printf("[*] Starting get images request at %s", time.Now().UTC().Format(time.RFC3339))

	// Read the images from the JSON file
	images, err := repositories.ReadImages()
	if err != nil {
		log.Printf("[*] Error: Unable to read images from file: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to open images.json file"})
	}
	log.Printf("[*] Successfully read images from repository")

	totalCustomImages := len(images.CustomImages)

	log.Printf("[*] Image Custom counts: %d", totalCustomImages)

	// Construct the response structure
	response := ImagesResponse{
		CustomImages:      images.CustomImages,
		TotalCustomImages: totalCustomImages,
	}

	// Send the JSON response
	return c.JSON(http.StatusOK, response)
}
