package CaptainLib

import "net/http"

type CaptainClient struct {
	BaseUrl string
	client *http.Client
}

func GetVersion() string {
	return "v0.0.0"
}

func NewCaptainClient(baseUrl string) *CaptainClient {
	return &CaptainClient{
		BaseUrl: baseUrl,
		client:  &http.Client{},
	}
}