package CaptainLib

import (
	"github.com/jarcoal/httpmock"
	"reflect"
	"testing"
)

func TestCaptainClient_GetAllAirspaces(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) *CaptainClient
		inspect func(r *CaptainClient, t *testing.T) //inspects receiver after test run

		want1      []Airspace
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "AirspaceGetAll",
			init: func(t *testing.T) *CaptainClient {
				return &CaptainClient{BaseUrl: "http://localhost:5000/"}
			},
			want1: []Airspace{
				{
					ID:        0,
					NetName:   "testNetName",
					HumanName: "testHumanName",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.init(t)
			helperRegisterAllAirspaceMocks()
			got1, err := receiver.GetAllAirspaces()
			helperDeregisterMocks()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CaptainClient.GetAllAirspaces got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("CaptainClient.GetAllAirspaces error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestCaptainClient_GetAirspaceByID(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *CaptainClient
		inspect func(r *CaptainClient, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      Airspace
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "AirspaceGetOne",
			init: func(t *testing.T) *CaptainClient {
				return &CaptainClient{BaseUrl: "http://localhost:5000/"}
			},
			args: func(t *testing.T) args {
				return args{
					id: 0,
				}
			},
			want1: Airspace{
				ID:        0,
				NetName:   "testNetName",
				HumanName: "testHumanName",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			helperRegisterAllAirspaceMocks()
			got1, err := receiver.GetAirspaceByID(tArgs.id)
			helperDeregisterMocks()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CaptainClient.GetAirspaceByID got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("CaptainClient.GetAirspaceByID error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestCaptainClient_CreateAirspace(t *testing.T) {
	type args struct {
		humanName string
		netName   string
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *CaptainClient
		inspect func(r *CaptainClient, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      Airspace
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "AirspaceCreateOne",
			init: func(t *testing.T) *CaptainClient {
				return &CaptainClient{BaseUrl: "http://localhost:5000/"}
			},
			args: func(t *testing.T) args {
				return args{
					humanName: "testHumanName",
					netName:   "testNetName",
				}
			},
			want1: Airspace{
				ID:        0,
				HumanName: "testHumanName",
				NetName:   "testNetName",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			helperRegisterAllAirspaceMocks()
			got1, err := receiver.CreateAirspace(tArgs.humanName, tArgs.netName)
			helperDeregisterMocks()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CaptainClient.CreateAirspace got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("CaptainClient.CreateAirspace error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func helperRegisterAllAirspaceMocks() {
	httpmock.Activate()
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:5000/airspaces",
		httpmock.NewStringResponder(
			200,
			`[{"ID":0,"NetName":"testNetName","HumanName":"testHumanName"}]`))
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:5000/airspace/0",
		httpmock.NewStringResponder(
			200,
			`{"ID":0,"NetName":"testNetName","HumanName":"testHumanName"}`))
	httpmock.RegisterResponder(
		"POST",
		"http://localhost:5000/airspace",
		httpmock.NewStringResponder(
			201,
			`{"ID":0,"NetName":"testNetName","HumanName":"testHumanName"}`))
	httpmock.RegisterResponder(
		"PUT",
		"http://localhost:5000/airspace/0",
		httpmock.NewStringResponder(
			201,
			``))
}

func helperDeregisterMocks() {
	httpmock.DeactivateAndReset()
}
