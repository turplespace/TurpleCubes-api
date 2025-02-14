package models

// ProxyRequest is the request body for the proxy service
type ProxyRequest struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	Subdomain string `json:"subdomain"`
}

type EditCubeRequest struct {
	CubeID      int       `json:"cube_id"`
	UpdatedCube Container `json:"updated_cube"`
}

type AddCubesRequest struct {
	WorkspaceID int         `json:"workspace_id"`
	Cubes       []Container `json:"cubes"`
}

type EditWorkspaceRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}
