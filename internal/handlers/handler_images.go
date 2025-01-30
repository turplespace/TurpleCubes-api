package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/turplespace/portos/internal/models"
	"github.com/turplespace/portos/internal/services/repositories"
)

type ImagesResponse struct {
	RepoImages        []models.Image `json:"repo_images"`
	CustomImages      []models.Image `json:"custom_images"`
	TotalImages       int            `json:"total_images"`
	TotalRepoImages   int            `json:"total_repo_images"`
	TotalCustomImages int            `json:"total_custom_images"`
}

func HandleGetImages(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GET-IMAGES] Starting get images request at %s", time.Now().UTC().Format(time.RFC3339))

	// Read the images from the JSON file
	images, err := repositories.ReadImages()
	if err != nil {
		log.Printf("[GET-IMAGES] Error: Unable to read images from file: %v", err)
		http.Error(w, "Unable to open images.json file", http.StatusInternalServerError)
		return
	}
	log.Printf("[GET-IMAGES] Successfully read images from repository")

	// Calculate the totals
	totalRepoImages := len(images.RepoImages)
	totalCustomImages := len(images.CustomImages)
	totalImages := totalRepoImages + totalCustomImages

	log.Printf("[GET-IMAGES] Image counts - Repo: %d, Custom: %d, Total: %d",
		totalRepoImages, totalCustomImages, totalImages)

	// Construct the response structure
	response := ImagesResponse{
		RepoImages:        images.RepoImages,
		CustomImages:      images.CustomImages,
		TotalImages:       totalImages,
		TotalRepoImages:   totalRepoImages,
		TotalCustomImages: totalCustomImages,
	}

	// Set the response header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Send the JSON response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("[GET-IMAGES] Error: Failed to encode response: %v", err)
		http.Error(w, "Error sending response", http.StatusInternalServerError)
		return
	}
	log.Printf("[GET-IMAGES] Successfully sent images response")
}
