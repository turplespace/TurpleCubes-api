package models

type CreateWorkspaceRequest struct {
	WorkspaceName string      `json:"workspace_name"`
	Desc          string      `json:"desc"`
	Containers    []Container `json:"containers"`
}
