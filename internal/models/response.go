package models

import (
	"time"
)
type GetCubesResponse struct {
	ContainerID	int		`json:"container_id"`
	Image         string        `json:"image"`
	ContainerName string        `json:"container_name"`
	IPAddress         string      `json:"ip_address,omitempty"`
	Status			string		`json:"status"`
}

type GetCubesByIdResponse struct {
	
	IPAddress		string		`json:"ip_address"`
	Status			string		`json:"status"`
	ContainerData *Container `json:"container_data"`
}

// WorkspaceWithContainerCounts includes workspace details and container counts
type WorkspaceWithContainerCounts struct {
    ID                int        `json:"id"`
    Name              string     `json:"name"`
    Desc              string     `json:"desc"`
    TotalContainers   int        `json:"total_containers"`
    RunningContainers int        `json:"running_containers"`
    CreatedAt         *time.Time `json:"created_at"`
}