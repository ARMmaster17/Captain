package Shared

type BuildPlaneRequest struct {
	S string `json:"s"`
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
