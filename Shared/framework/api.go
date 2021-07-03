package framework

import (
	"encoding/json"
	"fmt"
	"github.com/ARMmaster17/Captain/Shared/metadata"
	"github.com/rs/zerolog/log"
	"net/http"
)

// LatestAPIVersion is the newest version of the API supported by this instance of the framework.
const LatestAPIVersion int = 1

// RegisterAPIHandler Adds an http route with the API path prefix and version added in.
func (f *Framework) RegisterAPIHandler(version int, path string, handleFunction func(w http.ResponseWriter, r *http.Request), methods ...string) {
	f.RegisterHandler(fmt.Sprintf("api/v%d/%s", version, path), handleFunction, methods...)
}

// APIRespondWithJSON Creates a JSON response with the given payload and writes it to the HTTP stream.
func APIRespondWithJSON(code int, w http.ResponseWriter, payload interface{}) {
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

// RegisterCommonAPIRoutes Adds routes for common paths such as /version.
func (f *Framework) RegisterCommonAPIRoutes() {
	for i := 1; i <= LatestAPIVersion; i++ {
		f.RegisterAPIHandler(i, "version", handleCommonAPIApplicationVersion)
	}
}

// handleCommonAPIApplicationVersion Handles requests for the application version. Wrapper for metadata module.
func handleCommonAPIApplicationVersion(w http.ResponseWriter, r *http.Request) {
	APIRespondWithJSON(http.StatusOK, w, map[string]string{
		"version": metadata.GetCaptainVersion(),
	})
}
