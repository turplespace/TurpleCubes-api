package models

// ProxyRequest is the request body for the proxy service

type EditCubeRequest struct {
	CubeID      int       `json:"cube_id"`
	UpdatedCube Container `json:"updated_cube"`
}

type AddCubesRequest struct {
	WorkspaceID int         `json:"workspace_id"`
	Cubes       []Container `json:"cubes"`
}

type EditWorkspaceRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type AddProxyRequest struct {
	CubeID  int    `json:"cube_id"`
	Domain  string `json:"domain"`
	Port    int    `json:"port"`
	Type    string `json:"type"`
	Default bool   `json:"default"`
}

type EditProxyByIDRequest struct {
	Domain  string `json:"domain"`
	Port    int    `json:"port"`
	Type    string `json:"type"`
	Default bool   `json:"default"`
}
