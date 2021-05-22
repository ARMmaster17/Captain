package main

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_registerFormationHandlers(t *testing.T) {
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

			registerFormationHandlers(tArgs.router)

		})
	}
}

func Test_formationPath(t *testing.T) {
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

			got1 := formationPath(tArgs.uri)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("formationPath got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_handleFormationAllGet(t *testing.T) {
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

			handleFormationAllGet(tArgs.c)

		})
	}
}

func Test_handleFormationNewPost(t *testing.T) {
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

			handleFormationNewPost(tArgs.c)

		})
	}
}

func Test_handleFormationDelete(t *testing.T) {
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

			handleFormationDelete(tArgs.c)

		})
	}
}
