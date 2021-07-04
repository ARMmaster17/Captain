package framework

import (
	"context"
	"fmt"
	"github.com/ARMmaster17/Captain/Shared/config"
	"github.com/ARMmaster17/Captain/Shared/metadata"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

const defaultAPIPort = 3000

var (
	srv *http.Server
)

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
	log.Debug().Msgf("captain %s %s is starting up", config.ApplicationName, metadata.GetCaptainVersion())
	f.HTTPState = HTTPListening
	log.Debug().Msgf("captain %s is now listening at 0.0.0.0:%d", config.ApplicationName, f.GetPort())
	srv = &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%d", f.GetPort()),
		Handler:           f.Router,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		f.HTTPState = HTTPStopped
		log.Fatal().Err(err).Stack().Msg("http router stopped unexpectedly")
	}
}

// StopAsync Sends a stop message to the running http server job.
func (f *Framework) StopAsync() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		return err
	}
	f.HTTPState = HTTPStopped
	return nil
}

// RegisterHandler Adds a handler function that will be called when an http request is made on the specified path.
func (f *Framework) RegisterHandler(path string, handleFunction func(w http.ResponseWriter, r *http.Request), methods ...string) {
	if len(methods) == 0 {
		f.Router.HandleFunc(fmt.Sprintf("/%s", path), handleFunction)
	} else {
		f.Router.HandleFunc(fmt.Sprintf("/%s", path), handleFunction).Methods(methods...)
	}
}
