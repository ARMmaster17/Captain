package main

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_registerAirspaceHandlers(t *testing.T) {
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

			registerAirspaceHandlers(tArgs.router)

		})
	}
}

func Test_airspacePath(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 string
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := airspacePath(tArgs.uri)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("airspacePath got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_handleAirspaceAllGet(t *testing.T) {
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

			handleAirspaceAllGet(tArgs.c)

		})
	}
}

func Test_handleAirspaceNewPost(t *testing.T) {
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

			handleAirspaceNewPost(tArgs.c)

		})
	}
}

func Test_handleAirspaceDelete(t *testing.T) {
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

			handleAirspaceDelete(tArgs.c)

		})
	}
}
