package framework

import (
	"fmt"
	"github.com/ARMmaster17/Captain/Shared/config"
	"net/http"
)

const defaultAPIPort = 3000

// HTTPListenStatus Represents the various listening states of the framework HTTP server.
type HTTPListenStatus int

const (
	// HTTPStopped The server is not currently listening, or has stopped listening due to an error.
	HTTPStopped HTTPListenStatus = iota
	// HTTPListening The server is currently accepting requests.
	HTTPListening
)

// GetPort Handles the prioritization of port configurations, and returns the port that the service should listen on.
func (f *Framework) GetPort() int {
	port := config.GetAppInt("API_PORT")
	if port == 0 {
		port = defaultAPIPort
	}
	return port
}

// StartAsync Starts the framework listening loop in the background. Usually only used for testing purposes.
func (f *Framework) StartAsync() {
	go f.Start()
}

// Start Sets the HttpListenState and listens for incoming HTTP connections to be handles by the Framework instance.
func (f *Framework) Start() {
	f.HTTPState = HTTPListening
	//defer func() {
	//	f.HTTPState = HTTPStopped
	//}()
	//http.Handle("/", f.Router)
}

// StopAsync Sends a stop message to the running http server job.
func (f *Framework) StopAsync() {
	f.HTTPState = HTTPStopped
}

// RegisterHandler Adds a handler function that will be called when an http request is made on the specified path.
func (f *Framework) RegisterHandler(path string, handleFunction func(w http.ResponseWriter, r *http.Request), methods ...string) {
	if len(methods) == 0 {
		f.Router.HandleFunc(fmt.Sprintf("/%s", path), handleFunction)
	} else {
		f.Router.HandleFunc(fmt.Sprintf("/%s", path), handleFunction).Methods(methods...)
	}
}
