package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_registerRootHandlers(t *testing.T) {
	type args struct {
		router *gin.Engine
	}
	tests := []struct {
		name string
		args func(t *testing.T) args
	}{
		{
			name: "RegistersWithValidRouter",
			args: func(t *testing.T) args {
				return args{
					router: gin.Default(),
				}
			},
		},
	}

		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			registerRootHandlers(tArgs.router)

		})
	}
}

func Test_handleRootGet(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args func(t *testing.T) args
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			handleRootGet(tArgs.c)

		})
	}
}

func helperPerformTestRequest(method string, url string, body io.Reader) *httptest.ResponseRecorder {
	r := initRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body)
	r.ServeHTTP(w, req)
	return w
}
