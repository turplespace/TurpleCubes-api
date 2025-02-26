package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RestartNginxService restarts the Nginx service inside the Docker container.
func RestartNginxService() error {
	containerName := "turplecubes-proxy"

	// Restart the Nginx service inside the container
	cmd := exec.Command("docker", "exec", containerName, "nginx", "-s", "reload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error: ", err)
		return fmt.Errorf("failed to restart Nginx service inside container %s: %v", containerName, err)
	}

	return nil
}

/*
RemoveDataInFolder removes all data in the specified folder.
the proxy data will be deleted while restarting the turplecubes
to avoid domain not found error in ngxin
*/
func RemoveDataInFolder() error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	folder := fmt.Sprintf("%s_proxy", ex)

	// Ensure the folder exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return fmt.Errorf("configuration folder does not exist: %s", folder)
	}

	// Open the folder
	dir, err := os.Open(folder)
	if err != nil {
		return err
	}
	defer dir.Close()

	// Read all files in the folder
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}
	// Remove all files in the folder
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(folder, name))
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateFolderIfNotExists creates the folder if it doesn't exist.
func CreateFolderIfNotExists() error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	folder := fmt.Sprintf("%s_proxy", ex)

	// Ensure the folder exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return err
		}
	}

	return nil
}
