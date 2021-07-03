package longjob

import (
	"encoding/json"
	framework2 "github.com/ARMmaster17/Captain/Shared/framework"
	"net/http"
)

func RegisterLongjobQueue(framework *framework2.Framework, apiVersion int, path string, jobFunction func(payload *json.Decoder) (interface{}, error)) {
	framework.RegisterApiHandler(apiVersion, path, func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		result, err := jobFunction(decoder)
		if err != nil {
			framework2.ApiRespondWithJson(http.StatusInternalServerError, w, map[string]string{
				"error": err.Error(),
			})
		} else {
			framework2.ApiRespondWithJson(http.StatusOK, w, result)
		}
	})
}
