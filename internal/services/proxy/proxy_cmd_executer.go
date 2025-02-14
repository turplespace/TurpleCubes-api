package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

/*
StopAndRemoveContainer stops and removes the Docker container with the given name.
evey time restating the turplecube the proxy server will be removed and restarted to about container_name already exist error from docker demon
*/
func StopAndRemoveContainer() error {
	containerName := "turplecube-proxy"
	// Stop the container
	cmd := exec.Command("docker", "stop", containerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop container %s: %v", containerName, err)
	}

	// Remove the container
	cmd = exec.Command("docker", "rm", containerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove container %s: %v", containerName, err)
	}

	return nil
}

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

// RestartNginxService restarts the Nginx service inside the Docker container.
func RestartNginxService() error {
	containerName := "turplecube-proxy"

	// Restart the Nginx service inside the container
	cmd := exec.Command("docker", "exec", containerName, "nginx", "-s", "reload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
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
