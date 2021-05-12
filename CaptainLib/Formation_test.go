package CaptainLib

import (
	"github.com/jarcoal/httpmock"
	"reflect"
	"testing"
)

func TestCaptainClient_GetAllFormations(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) *CaptainClient
		inspect func(r *CaptainClient, t *testing.T) //inspects receiver after test run

		want1      []Formation
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "FormationGetAll",
			init: func(t *testing.T) *CaptainClient {
				return &CaptainClient{BaseUrl: "http://localhost:5000/"}
			},
			want1: []Formation{
				{
					ID: 0,
					FlightID: 0,
					Name: "testFormation",
					CPU: 1,
					RAM: 128,
					Disk: 8,
					BaseName: "test",
					Domain: "example.com",
					TargetCount: 0,
				},
			},
			wantErr: false,
		},
	}

		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.init(t)
			helperRegisterAllFormationMocks()
			got1, err := receiver.GetAllFormations()
			helperDeregisterMocks()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CaptainClient.GetAllFormations got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("CaptainClient.GetAllFormations error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestCaptainClient_GetFormationByID(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *CaptainClient
		inspect func(r *CaptainClient, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      Formation
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "FormationGetOneByIDValid",
			init: func(t *testing.T) *CaptainClient {
				return &CaptainClient{BaseUrl: "http://localhost:5000/"}
			},
			args: func(t *testing.T) args {
				return args{
					id: 0,
				}
			},
			want1: Formation{
				ID: 0,
				FlightID: 0,
				Name: "testFormation",
				CPU: 1,
				RAM: 128,
				Disk: 8,
				BaseName: "test",
				Domain: "example.com",
				TargetCount: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			helperRegisterAllFormationMocks()
			got1, err := receiver.GetFormationByID(tArgs.id)
			helperDeregisterMocks()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CaptainClient.GetFormationByID got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("CaptainClient.GetFormationByID error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestCaptainClient_CreateFormation(t *testing.T) {
	type args struct {
		name        string
		flightID    int
		CPU         int
		RAM         int
		disk        int
		baseName    string
		domain      string
		targetCount int
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *CaptainClient
		inspect func(r *CaptainClient, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      Formation
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "FormationCreateValid",
			init: func(t *testing.T) *CaptainClient {
				return &CaptainClient{BaseUrl: "http://localhost:5000/"}
			},
			args: func(t *testing.T) args {
				return args{
					flightID: 0,
					name: "testFormation",
					CPU: 1,
					RAM: 128,
					disk: 8,
					baseName: "test",
					domain: "example.com",
					targetCount: 0,
				}
			},
			want1: Formation{
				ID: 0,
				FlightID: 0,
				Name: "testFormation",
				CPU: 1,
				RAM: 128,
				Disk: 8,
				BaseName: "test",
				Domain: "example.com",
				TargetCount: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			helperRegisterAllFormationMocks()
			got1, err := receiver.CreateFormation(tArgs.name, tArgs.flightID, tArgs.CPU, tArgs.RAM, tArgs.disk, tArgs.baseName, tArgs.domain, tArgs.targetCount)
			helperDeregisterMocks()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CaptainClient.CreateFormation got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("CaptainClient.CreateFormation error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestCaptainClient_UpdateFormation(t *testing.T) {
	type args struct {
		formation Formation
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
			err := receiver.UpdateFormation(tArgs.formation)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("CaptainClient.UpdateFormation error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestCaptainClient_DeleteFormation(t *testing.T) {
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
			err := receiver.DeleteFormation(tArgs.id)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("CaptainClient.DeleteFormation error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func helperRegisterAllFormationMocks() {
	httpmock.Activate()
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:5000/formations",
		httpmock.NewStringResponder(
			200,
			`[{"ID":0,"FlightID":0,"Name":"testFormation","CPU":1,"RAM":128,"Disk":8,"BaseName":"test","Domain":"example.com","TargetCount":0}]`))
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:5000/formation/0",
		httpmock.NewStringResponder(
			200,
			`{"ID":0,"FlightID":0,"Name":"testFormation","CPU":1,"RAM":128,"Disk":8,"BaseName":"test","Domain":"example.com","TargetCount":0}`))
	httpmock.RegisterResponder(
		"POST",
		"http://localhost:5000/formation",
		httpmock.NewStringResponder(
			201,
			`{"ID":0,"FlightID":0,"Name":"testFormation","CPU":1,"RAM":128,"Disk":8,"BaseName":"test","Domain":"example.com","TargetCount":0}`))
	httpmock.RegisterResponder(
		"PUT",
		"http://localhost:5000/formation/0",
		httpmock.NewStringResponder(
			201,
			``))
}
