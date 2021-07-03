package longjob

import (
	"encoding/json"
	framework2 "github.com/ARMmaster17/Captain/Shared/framework"
	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_RegisterLongjobRoute(t *testing.T) {
	framework := framework2.NewFramework("shared")
	require.Equal(t, "shared", framework.AppName)
	require.NotNil(t, framework.Router)
	RegisterLongjobQueue(&framework, 1, "test", func(payload JobInput) (JobOutput, error){
		return map[string]string{
			"test_key": "test_value",
		}, nil
	})
	req, err := http.NewRequest("POST", "/api/v1/test", strings.NewReader("{}"))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	framework.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	testOutput, err := json.Marshal(map[string]interface{}{
		"jobId": 0,
	})
	require.NoError(t, err)
	assert.Equal(t, string(testOutput), rr.Body.String())
}

func Test_RegisterLongjobRouteWithNoBodyBug(t *testing.T) {
	framework := framework2.NewFramework("shared")
	require.Equal(t, "shared", framework.AppName)
	require.NotNil(t, framework.Router)
	RegisterLongjobQueue(&framework, 1, "test", func(payload JobInput) (JobOutput, error){
		return map[string]string{
			"test_key": "test_value",
		}, nil
	})
	req, err := http.NewRequest("POST", "/api/v1/test", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	framework.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	testOutput, err := json.Marshal(map[string]interface{}{
		"jobId": 0,
	})
	require.NoError(t, err)
	assert.Equal(t, string(testOutput), rr.Body.String())
}

func Test_LongjobStatus(t *testing.T) {
	framework := framework2.NewFramework("shared")
	require.Equal(t, "shared", framework.AppName)
	require.NotNil(t, framework.Router)
	RegisterLongjobQueue(&framework, 1, "test", func(payload JobInput) (JobOutput, error){
		time.Sleep(1 * time.Second)
		return map[string]string{
			"test_key": "test_value",
		}, nil
	})
	req, err := http.NewRequest("POST", "/api/v1/test/status/0", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	framework.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	testOutput, err := json.Marshal(map[string]interface{}{
		"jobId": 0,
		"done": false,
	})
	require.NoError(t, err)
	assert.Equal(t, string(testOutput), rr.Body.String())
}

func Test_LongjobResults(t *testing.T) {
	framework := framework2.NewFramework("shared")
	require.Equal(t, "shared", framework.AppName)
	require.NotNil(t, framework.Router)
	RegisterLongjobQueue(&framework, 1, "test", func(payload JobInput) (JobOutput, error){
		return map[string]string{
			"test_key": "test_value",
		}, nil
	})
	req, err := http.NewRequest("POST", "/api/v1/test", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	framework.Router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	req, err = http.NewRequest("POST", "/api/v1/test/result/0", nil)
	require.NoError(t, err)
	rr = httptest.NewRecorder()
	time.Sleep(10 * time.Millisecond)
	framework.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	testOutput, err := json.Marshal(map[string]interface{}{
		"jobId": 0,
		"done": true,
		"err": nil,
		"result": map[string]string{
			"test_key": "test_value",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, string(testOutput), rr.Body.String())
}

func Test_LongjobResultsWithIncompleteError(t *testing.T) {
	framework := framework2.NewFramework("shared")
	require.Equal(t, "shared", framework.AppName)
	require.NotNil(t, framework.Router)
	RegisterLongjobQueue(&framework, 1, "test", func(payload JobInput) (JobOutput, error){
		time.Sleep(1 * time.Second)
		return map[string]string{
			"test_key": "test_value",
		}, nil
	})
	req, err := http.NewRequest("POST", "/api/v1/test", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	framework.Router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	req, err = http.NewRequest("POST", "/api/v1/test/result/0", nil)
	require.NoError(t, err)
	rr = httptest.NewRecorder()
	framework.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	testOutput, err := json.Marshal(map[string]interface{}{
		"error": "job #0 is not complete",
	})
	require.NoError(t, err)
	assert.Equal(t, string(testOutput), rr.Body.String())
}
