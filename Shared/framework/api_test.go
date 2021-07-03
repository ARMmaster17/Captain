package framework

import (
	"encoding/json"
	"github.com/ARMmaster17/Captain/Shared/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_FrameworkRegistersApiRoute(t *testing.T) {
	framework := NewFramework("shared")
	require.Equal(t, "shared", framework.AppName)
	require.NotNil(t, framework.Router)
	framework.RegisterAPIHandler(1, "test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req, err := http.NewRequest("GET", "/api/v1/test", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	framework.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_FrameworkRegistersCommonVersionRoute(t *testing.T) {
	framework := NewFramework("shared")
	require.Equal(t, "shared", framework.AppName)
	require.NotNil(t, framework.Router)
	framework.RegisterCommonAPIRoutes()
	req, err := http.NewRequest("GET", "/api/v1/version", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	framework.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	testOutput, err := json.Marshal(map[string]string{
		"version": metadata.GetCaptainVersion(),
	})
	require.NoError(t, err)
	assert.Equal(t, string(testOutput), rr.Body.String())
}

func Test_RespondWithJson(t *testing.T) {
	framework := NewFramework("shared")
	require.Equal(t, "shared", framework.AppName)
	require.NotNil(t, framework.Router)
	framework.RegisterAPIHandler(1, "test", func(w http.ResponseWriter, r *http.Request) {
		APIRespondWithJSON(http.StatusOK, w, map[string]string{
			"test_key": "test_value",
		})
	})
	req, err := http.NewRequest("GET", "/api/v1/test", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	framework.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	testOutput, err := json.Marshal(map[string]string{
		"test_key": "test_value",
	})
	require.NoError(t, err)
	assert.Equal(t, string(testOutput), rr.Body.String())
}

func Test_FrameworkRegisterPostRoute(t *testing.T) {
	framework := NewFramework("shared")
	require.Equal(t, "shared", framework.AppName)
	require.NotNil(t, framework.Router)
	framework.RegisterAPIHandler(1, "test", func(w http.ResponseWriter, r *http.Request) {
		APIRespondWithJSON(http.StatusOK, w, map[string]string{
			"test_key": "test_value",
		})
	}, "POST")
	req, err := http.NewRequest("POST", "/api/v1/test", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	framework.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	testOutput, err := json.Marshal(map[string]string{
		"test_key": "test_value",
	})
	require.NoError(t, err)
	assert.Equal(t, string(testOutput), rr.Body.String())
}
