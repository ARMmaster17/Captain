package Shared

type BuildPlaneRequest struct {
	PlaneID string `json:"planeId"`
}

type BuildPlaneResponse struct {
	V string `json:"v"`
	Err string `json:"err,omitempty"`
}

type ProvisionPlaneRequest struct {
	S string `json:"s"`
}

type ProvisionPlaneResponse struct {
	V string `json:"v"`
	Err string `json:"err,omitempty"`
}
