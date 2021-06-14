package CaptainLib

import (
	"github.com/jarcoal/httpmock"
	"reflect"
	"testing"
)

func TestCaptainClient_GetAllFlights(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) *CaptainClient
		inspect func(r *CaptainClient, t *testing.T) //inspects receiver after test run

		want1      []Flight
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "FlightGetAll",
			init: func(t *testing.T) *CaptainClient {
				return &CaptainClient{BaseUrl: "http://localhost:5000/"}
			},
			want1: []Flight{
				{
					ID:         0,
					AirspaceID: 0,
					Name:       "testFlight",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.init(t)
			helperRegisterAllFlightMocks()
			got1, err := receiver.GetAllFlights()
			helperDeregisterMocks()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CaptainClient.GetAllFlights got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("CaptainClient.GetAllFlights error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestCaptainClient_GetFlightByID(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *CaptainClient
		inspect func(r *CaptainClient, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      Flight
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "FlightGetByIDValid",
			init: func(t *testing.T) *CaptainClient {
				return &CaptainClient{BaseUrl: "http://localhost:5000/"}
			},
			args: func(t *testing.T) args {
				return args{
					id: 0,
				}
			},
			want1: Flight{
				ID:         0,
				AirspaceID: 0,
				Name:       "testFlight",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			helperRegisterAllFlightMocks()
			got1, err := receiver.GetFlightByID(tArgs.id)
			helperDeregisterMocks()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CaptainClient.GetFlightByID got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("CaptainClient.GetFlightByID error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestCaptainClient_CreateFlight(t *testing.T) {
	type args struct {
		name       string
		airspaceID int
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *CaptainClient
		inspect func(r *CaptainClient, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      Flight
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "FlightCreateValid",
			init: func(t *testing.T) *CaptainClient {
				return &CaptainClient{BaseUrl: "http://localhost:5000/"}
			},
			args: func(t *testing.T) args {
				return args{
					name:       "testFlight",
					airspaceID: 0,
				}
			},
			want1: Flight{
				ID:         0,
				AirspaceID: 0,
				Name:       "testFlight",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			helperRegisterAllFlightMocks()
			got1, err := receiver.CreateFlight(tArgs.name, tArgs.airspaceID)
			helperDeregisterMocks()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CaptainClient.CreateFlight got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("CaptainClient.CreateFlight error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestCaptainClient_UpdateFlight(t *testing.T) {
	type args struct {
		flight Flight
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *CaptainClient
		inspect func(r *CaptainClient, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			helperRegisterAllFlightMocks()
			err := receiver.UpdateFlight(tArgs.flight)
			helperDeregisterMocks()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("CaptainClient.UpdateFlight error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestCaptainClient_DeleteFlight(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *CaptainClient
		inspect func(r *CaptainClient, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			helperRegisterAllFlightMocks()
			err := receiver.DeleteFlight(tArgs.id)
			helperDeregisterMocks()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("CaptainClient.DeleteFlight error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func helperRegisterAllFlightMocks() {
	httpmock.Activate()
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:5000/flights",
		httpmock.NewStringResponder(
			200,
			`[{"ID":0,"AirspaceID":0,"Name":"testFlight"}]`))
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:5000/flight/0",
		httpmock.NewStringResponder(
			200,
			`{"ID":0,"AirspaceID":0,"Name":"testFlight"}`))
	httpmock.RegisterResponder(
		"POST",
		"http://localhost:5000/flight",
		httpmock.NewStringResponder(
			201,
			`{"ID":0,"AirspaceID":0,"Name":"testFlight"}`))
	httpmock.RegisterResponder(
		"PUT",
		"http://localhost:5000/flight/0",
		httpmock.NewStringResponder(
			201,
			``))
}
