package longjob

import (
	"encoding/json"
	"fmt"
	framework2 "github.com/ARMmaster17/Captain/Shared/framework"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// RegisterLongjobQueue Creates a new Queue and registers all Longjob routes with the framework. The API path can be
// used to POST new jobs, and two other routes /status/:id and /result/:id will be created for post-completion data.
func RegisterLongjobQueue(framework *framework2.Framework, apiVersion int, path string, jobFunction JobFunction) {
	jobQueue := NewQueue(path, 1, jobFunction)
	registerLongjobEnqueueRoute(framework, apiVersion, path, jobQueue)
	registerLongjobStatusRoute(framework, apiVersion, path, jobQueue)
	registerLongjobResultRoute(framework, apiVersion, path, jobQueue)
}

// registerLongjobEnqueueRoute Registers the primary route for submitting jobs to the queue.
func registerLongjobEnqueueRoute(framework *framework2.Framework, apiVersion int, path string, jobQueue Queue) {
	framework.RegisterAPIHandler(apiVersion, path, func(w http.ResponseWriter, r *http.Request) {
		var input map[string]interface{}
		if r.Body == nil {
			input = nil
		} else {
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&input)
			if err != nil {
				framework2.APIRespondWithJSON(http.StatusInternalServerError, w, map[string]string{
					"error": err.Error(),
				})
				return
			}
		}
		jobID := jobQueue.Enqueue(input)
		framework2.APIRespondWithJSON(http.StatusOK, w, map[string]uint64{
			"jobID": jobID,
		})
	}, "POST")
}

// registerLongjobStatusRoute Registers the route for checking the status on a previously-submitted job.
func registerLongjobStatusRoute(framework *framework2.Framework, apiVersion int, path string, jobQueue Queue) {
	framework.RegisterAPIHandler(apiVersion, fmt.Sprintf("%s/status/{jobid:[0-9]+}", path), func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		jobID, err := strconv.ParseUint(vars["jobid"], 10, 64)
		if err != nil {
			framework2.APIRespondWithJSON(http.StatusInternalServerError, w, map[string]string{
				"error": err.Error(),
			})
			return
		}
		framework2.APIRespondWithJSON(http.StatusOK, w, map[string]interface{}{
			"jobID": jobID,
			"done":  jobQueue.IsJobDone(jobID),
		})
	})
}

// registerLongjobResultRoute Registers the routes used for getting results from a previous job run.
func registerLongjobResultRoute(framework *framework2.Framework, apiVersion int, path string, jobQueue Queue) {
	framework.RegisterAPIHandler(apiVersion, fmt.Sprintf("%s/result/{jobid:[0-9]+}", path), func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		jobID, err := strconv.ParseUint(vars["jobid"], 10, 64)
		if err != nil {
			framework2.APIRespondWithJSON(http.StatusInternalServerError, w, map[string]string{
				"error": err.Error(),
			})
			return
		}
		if !jobQueue.IsJobDone(jobID) {
			framework2.APIRespondWithJSON(http.StatusInternalServerError, w, map[string]string{
				"error": fmt.Sprintf("job #%d is not complete", jobID),
			})
			return
		}
		result, err := jobQueue.GetResult(jobID)
		framework2.APIRespondWithJSON(http.StatusOK, w, map[string]interface{}{
			"jobID":  jobID,
			"done":   true,
			"result": result,
			"err":    err,
		})
	})
}
