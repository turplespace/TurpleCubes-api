package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RunNginxDockerContainer runs an Nginx Docker container with the specified configuration folder mounted as a volume.
func RunNginxDockerContainer() error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	folder := fmt.Sprintf("%s_proxy", ex)

	// Ensure the folder exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return fmt.Errorf("configuration folder does not exist: %s", folder)
	}

	// Run the Nginx Docker container with the configuration folder mounted
	cmd := exec.Command("docker", "run", "--name", "turplecube-proxy", "-v", fmt.Sprintf("%s:/etc/nginx/conf.d", folder), "-p", "80:80", "-d", "nginx")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RestartNginxService restarts the Nginx service.
func RestartNginxService() error {
	cmd := exec.Command("systemctl", "restart", "nginx")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RemoveDataInFolder removes all data in the specified folder.
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

	// Remove all files in the folder
	dir, err := os.Open(folder)
	if err != nil {
		return err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

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
