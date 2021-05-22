package main

import (
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
		//TODO: Add test cases
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
