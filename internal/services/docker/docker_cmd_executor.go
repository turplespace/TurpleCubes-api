package docker

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/turplespace/portos/internal/models"
)

// StartContainer starts a new container
func StartContainer(container models.Container) error {
	// Check if the container already exists
	cmdCheck := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("name=%s", container.Name), "--format", "{{.ID}}")
	existingContainerID, err := cmdCheck.Output()
	if err != nil {
		return fmt.Errorf("failed to check existing container: %v", err)
	}

	if len(existingContainerID) > 0 {
		// Stop the existing container
		cmdStop := exec.Command("docker", "stop", container.Name)
		if output, err := cmdStop.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to stop existing container: %v\nOutput: %s", err, string(output))
		}

		// Remove the existing container
		cmdRemove := exec.Command("docker", "rm", container.Name)
		if output, err := cmdRemove.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to remove existing container: %v\nOutput: %s", err, string(output))
		}

		log.Printf("Existing container %s stopped and removed successfully", container.Name)
	}

	args := []string{"run", "-d"}

	// Add name
	args = append(args, "--name", container.Name)

	// Add ports
	for _, port := range container.Ports {
		args = append(args, "-p", port)
	}

	// Add environment variables
	for _, env := range container.EnvironmentVars {
		args = append(args, "-e", env)
	}

	// Add volumes
	for hostPath, containerPath := range container.Volumes {
		if strings.Contains(hostPath, "[DEFAULT]") {
			ex, err := os.Executable()
			if err != nil {
				panic(err)
			}
			fmt.Println(ex)
			path := fmt.Sprintf("%s_volumes", ex)
			fmt.Println(path)
			hostPath = strings.Replace(hostPath, "[DEFAULT]", path, 1)
		}
		args = append(args, "-v", fmt.Sprintf("%s:%s", hostPath, containerPath))
	}

	// Add labels
	for _, label := range container.Labels {
		args = append(args, "-l", label)
	}

	// Add resource limits
	if container.ResourceLimits.CPUs != "" {
		args = append(args, "--cpus", container.ResourceLimits.CPUs)
	}
	if container.ResourceLimits.Memory != "" {
		args = append(args, "--memory", container.ResourceLimits.Memory)
	}

	// Add image name
	args = append(args, container.Image)

	// Execute command
	cmdRun := exec.Command("docker", args...)
	output, err := cmdRun.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start container: %v\nOutput: %s", err, string(output))
	}

	log.Printf("Container %s started successfully: %s", container.Name, string(output))
	return nil
}

// StopContainer stops a running container
func StopContainer(containerName string) error {
	cmd := exec.Command("docker", "stop", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop container: %v\nOutput: %s", err, string(output))
	}

	log.Printf("Container %s stopped successfully: %s", containerName, string(output))
	return nil
}

// RestartContainer restarts a container
func RestartContainer(containerName string) error {
	cmd := exec.Command("docker", "restart", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to restart container: %v\nOutput: %s", err, string(output))
	}

	log.Printf("Container %s restarted successfully: %s", containerName, string(output))
	return nil
}

// CommitContainer creates a new image from a container's changes
func CommitContainer(containerName, newImageName, newImageTag string) error {
	args := []string{"commit", containerName, fmt.Sprintf("%s:%s", newImageName, newImageTag)}
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to commit container: %v\nOutput: %s", err, string(output))
	}

	log.Printf("Container %s committed successfully as %s:%s: %s",
		containerName, newImageName, newImageTag, string(output))
	return nil
}

// PrintCommand prints the docker command that would be executed (for debugging)
func PrintCommand(args []string) string {
	cmd := "docker"
	for _, arg := range args {
		if arg == "" {
			continue
		}
		if containsSpace := false; containsSpace {
			cmd += fmt.Sprintf(" %q", arg)
		} else {
			cmd += fmt.Sprintf(" %s", arg)
		}
	}
	return cmd
}
