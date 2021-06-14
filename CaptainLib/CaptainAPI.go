package CaptainLib

import "net/http"

// CaptainClient is a stateful client for connecting to a Captain ATC instance.
type CaptainClient struct {
	BaseUrl string
	client  *http.Client
}

// GetVersion returns the current library version (not implemented)
func GetVersion() string {
	return "v0.0.0"
}

// NewCaptainClient creates a new client instance.
func NewCaptainClient(baseUrl string) *CaptainClient {
	return &CaptainClient{
		BaseUrl: baseUrl,
		client:  &http.Client{},
	}
}
