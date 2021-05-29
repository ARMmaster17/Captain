package main

import (
	"reflect"
	"testing"

	"github.com/ARMmaster17/Captain/CaptainLib"
	"github.com/gin-gonic/gin"
)

func Test_getCaptainClient(t *testing.T) {
	tests := []struct {
		name string

		want1 *CaptainLib.CaptainClient
	}{
		{
			name: "gets a new CaptainLib client",
			want1: CaptainLib.NewCaptainClient("http://localhost:5000/"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := getCaptainClient()

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getCaptainClient got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_getURLIDParameter(t *testing.T) {
	type args struct {
		name string
		c    *gin.Context
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      int
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := getURLIDParameter(tArgs.name, tArgs.c)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getURLIDParameter got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("getURLIDParameter error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func Test_forceIntRead(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 int
	}{
		{
			name: "ReadsValidValueZero",
			args: func(t *testing.T) args {
				return args{
					input: "0",
				}
			},
			want1: 0,
		},
		{
			name: "ReadsValidValueOne",
			args: func(t *testing.T) args {
				return args{
					input: "1",
				}
			},
			want1: 1,
		},
		{
			name: "ReadsInvalidValue",
			args: func(t *testing.T) args {
				return args{
					input: "invalid",
				}
			},
			want1: 0,
		},
		{
			name: "ReadsMax32ValidValue",
			args: func(t *testing.T) args {
				return args{
					input: "2147483647",
				}
			},
			want1: 2147483647,
		},
	}

		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := forceIntRead(tArgs.input)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("forceIntRead got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_getAirspaceFromURLParameter(t *testing.T) {
	type args struct {
		c      *gin.Context
		client *CaptainLib.CaptainClient
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      CaptainLib.Airspace
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := getAirspaceFromURLParameter(tArgs.c, tArgs.client)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getAirspaceFromURLParameter got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("getAirspaceFromURLParameter error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func Test_getFlightFromURLParameter(t *testing.T) {
	type args struct {
		c      *gin.Context
		client *CaptainLib.CaptainClient
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      CaptainLib.Flight
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := getFlightFromURLParameter(tArgs.c, tArgs.client)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getFlightFromURLParameter got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("getFlightFromURLParameter error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
