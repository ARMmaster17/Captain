package main

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_registerFlightHandlers(t *testing.T) {
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

			registerFlightHandlers(tArgs.router)

		})
	}
}

func Test_flightPath(t *testing.T) {
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

			got1 := flightPath(tArgs.uri)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("flightPath got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_handleFlightAllGet(t *testing.T) {
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

			handleFlightAllGet(tArgs.c)

		})
	}
}

func Test_handleFlightNewPost(t *testing.T) {
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

			handleFlightNewPost(tArgs.c)

		})
	}
}

func Test_handleFlightDelete(t *testing.T) {
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

			handleFlightDelete(tArgs.c)

		})
	}
}
