package longjob

import (
	"encoding/json"
	framework2 "github.com/ARMmaster17/Captain/Shared/framework"
	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_RegisterLongjobRoute(t *testing.T) {
	framework := framework2.NewFramework("shared")
	require.Equal(t, "shared", framework.AppName)
	require.NotNil(t, framework.Router)
	RegisterLongjobQueue(&framework, 1, "test", func(payload *json.Decoder) (interface{}, error){
		return map[string]string{
			"test_key": "test_value",
		}, nil
	})
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
