package framework

import (
	"encoding/json"
	"fmt"
	"github.com/ARMmaster17/Captain/Shared/metadata"
	"github.com/rs/zerolog/log"
	"net/http"
)

const LatestApiVersion int = 1

// RegisterApiHandler Adds an http route with the API path prefix and version added in.
func (f *Framework) RegisterApiHandler(version int, path string, handleFunction func(w http.ResponseWriter, r *http.Request), methods ...string) {
	f.RegisterHandler(fmt.Sprintf("api/v%d/%s", version, path), handleFunction, methods...)
}

// ApiRespondWithJson Creates a JSON response with the given payload and writes it to the HTTP stream.
func ApiRespondWithJson(code int, w http.ResponseWriter, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Stack().Msg("unable to convert data to json")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Error().Err(err).Stack().Msg("unable to write data to http stream")
		return
	}
}

// RegisterCommonApiRoutes Adds routes for common paths such as /version.
func (f *Framework) RegisterCommonApiRoutes() {
	for i := 1; i <= LatestApiVersion; i++ {
		f.RegisterApiHandler(i, "version", handleCommonApiApplicationVersion)
	}
}

// handleCommonApiApplicationVersion Handles requests for the application version. Wrapper for metadata module.
func handleCommonApiApplicationVersion(w http.ResponseWriter, r *http.Request) {
	ApiRespondWithJson(http.StatusOK, w, map[string]string{
		"version": metadata.GetCaptainVersion(),
	})
}