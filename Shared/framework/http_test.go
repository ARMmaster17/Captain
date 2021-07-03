package framework

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func Test_FrameworkListensOnConfiguredPort(t *testing.T) {
	framework, err := NewFramework("shared")
	require.NoError(t, err)
	require.Equal(t, "shared", framework.AppName)
	require.NoError(t, os.Setenv("CAPTAIN_SHARED_API_PORT", "1000"))
	assert.Equal(t, 1000, framework.GetPort())
}

func Test_FrameworkListensOnDefaultPort(t *testing.T) {
	framework, err := NewFramework("shared")
	require.NoError(t, err)
	require.Equal(t, "shared", framework.AppName)
	require.NoError(t, os.Setenv("CAPTAIN_SHARED_API_PORT", ""))
	assert.Equal(t, 3000, framework.GetPort())
}

func Test_FrameworkStartsWithRouterStopped(t *testing.T) {
	framework, err := NewFramework("shared")
	require.NoError(t, err)
	require.Equal(t, "shared", framework.AppName)
	require.NotNil(t, framework.Router)
	assert.Equal(t, HTTPStopped, framework.HTTPState)
}

func Test_FrameworkStartingAsyncSetsListenState(t *testing.T) {
	framework, err := NewFramework("shared")
	require.NoError(t, err)
	require.Equal(t, "shared", framework.AppName)
	require.NotNil(t, framework.Router)
	framework.StartAsync()
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, HTTPListening, framework.HTTPState)
}

func Test_FrameworkStopSetsListenState(t *testing.T) {
	framework, err := NewFramework("shared")
	require.NoError(t, err)
	require.Equal(t, "shared", framework.AppName)
	require.NotNil(t, framework.Router)
	framework.StartAsync()
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, HTTPListening, framework.HTTPState)
	framework.StopAsync()
	assert.Equal(t, HTTPStopped, framework.HTTPState)
}

func Test_FrameworkRegistersHandler(t *testing.T) {
	framework, err := NewFramework("shared")
	require.NoError(t, err)
	require.Equal(t, "shared", framework.AppName)
	require.NotNil(t, framework.Router)
	framework.RegisterHandler("test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	framework.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
