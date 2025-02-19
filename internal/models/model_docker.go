package models

// Container represents the configuration for a single container in the workspace
// It includes all necessary settings like resource limits, environment variables,
// port mappings, and volume configurations.
type Container struct {
	ID              int               `json:"id"`               // Container ID
	Name            string            `json:"name"`             // Container name
	Image           string            `json:"image"`            // Docker image to use
	Ports           []string          `json:"ports"`            // Port mappings (host:container)
	EnvironmentVars []string          `json:"environment_vars"` // Environment variables
	ResourceLimits  ResourceLimits    `json:"resource_limits"`  // CPU, memory, and swap limits
	Volumes         map[string]string `json:"volumes"`          // Volume mappings (source:target)
	Labels          []string          `json:"labels"`           // Container labels
}

// ResourceLimits defines the computational resources allocated to a container
// This includes CPU cores, memory allocation, and swap space.
type ResourceLimits struct {
	CPUs   string `json:"cpus,omitempty"`   // CPU cores allocation (e.g., "1.0")
	Memory string `json:"memory,omitempty"` // Memory limit (e.g., "1G")
}
