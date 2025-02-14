package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// Function to get the status of a Docker container by name
func GetContainerStatus(containerName string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("failed to create Docker client: %v", err)
	}

	// Inspect the container to get detailed information
	containerJSON, err := cli.ContainerInspect(context.Background(), containerName)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container %s: %v", containerName, err)
	}

	// Return the status of the container
	return containerJSON.State.Status, nil
}

// GetContainersByLabel retrieves a list of containers with a specific label
func GetContainersByLabel(labelKey, labelValue string) ([]types.Container, error) {
	cli, err := client.NewClientWithOpts(
		client.WithVersion("1.45"), // Set the API version to 1.45
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	// Define the label to filter by
	labelFilter := filters.NewArgs()
	labelFilter.Add("label", fmt.Sprintf("%s=%s", labelKey, labelValue))

	// Retrieve a list of containers with the specific label
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{
		Filters: labelFilter,
	})
	if err != nil {
		return nil, err
	}

	// Return the list of filtered containers
	return containers, nil
}

// CountContainers counts the number of containers with a specific label
func CountContainersByLabel(labelKey, labelValue string) (int, error) {
	cli, err := client.NewClientWithOpts(
		client.WithVersion("1.45"), // Set the API version to 1.45
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return 0, err
	}

	// Define the label to filter by
	labelFilter := filters.NewArgs()
	labelFilter.Add("label", fmt.Sprintf("%s=%s", labelKey, labelValue))

	// Retrieve a list of containers with the specific label
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{
		Filters: labelFilter,
	})
	if err != nil {
		return 0, err
	}

	// Return the count of filtered containers
	return len(containers), nil
}

// Function to get the IP address of a Docker container by name
func GetContainerIPAddress(containerName string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("failed to create Docker client: %v", err)
	}

	// Inspect the container to get detailed information
	containerJSON, err := cli.ContainerInspect(context.Background(), containerName)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %v", err)
	}

	// Get the IP address from the container's network settings
	if containerJSON.NetworkSettings != nil && len(containerJSON.NetworkSettings.Networks) > 0 {
		for _, network := range containerJSON.NetworkSettings.Networks {
			return network.IPAddress, nil
		}
	}

	return "", fmt.Errorf("no IP address found for container: %s", containerName)
}
