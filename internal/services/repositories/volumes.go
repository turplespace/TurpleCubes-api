package repositories

import (
	"os"
	"path/filepath"
)

func GetVolumesBasePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	// Get the directory containing the executable
	exePath := filepath.Dir(ex)
	// Create the volumes directory path
	volumesPath := filepath.Join(exePath, "turplecube_volumes")

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(volumesPath, 0755); err != nil {
		return "", err
	}

	return volumesPath, nil
}
