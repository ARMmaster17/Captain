package IPAM

import (
	"github.com/ARMmaster17/Captain/ATC/DB"
	"net"
	"reflect"
	"sync"
	"testing"

	"gorm.io/gorm"
)

func TestIPAM_getAllReservedBlocks(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) *IPAM
		inspect func(r *IPAM, t *testing.T) //inspects receiver after test run

		want1      []ReservedBlock
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "gets empty IP block",
			init: func(t *testing.T) *IPAM {
				dbt, _ := DB.ConnectToDB()
				ipam := IPAM{
					db: dbt,
					mutex: &sync.Mutex{},
				}
				ipam.performMigrations()
				return &ipam
			},
			inspect: func(r *IPAM, t *testing.T) {
				
			},
			want1: []ReservedBlock{},
			wantErr: false,
			inspectErr: func(err error, t *testing.T) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.init(t)
			got1, err := receiver.getAllReservedBlocks()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("IPAM.getAllReservedBlocks got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("IPAM.getAllReservedBlocks error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestReservedBlock_hasAvailableAddress(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *ReservedBlock
		inspect func(r *ReservedBlock, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      bool
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "addresses available in empty IP block",
			init: func(t *testing.T) *ReservedBlock {
				
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			got1, err := receiver.hasAvailableAddress(tArgs.db)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ReservedBlock.hasAvailableAddress got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("ReservedBlock.hasAvailableAddress error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestReservedBlock_reserveAddress(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *ReservedBlock
		inspect func(r *ReservedBlock, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      net.IP
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			got1, err := receiver.reserveAddress(tArgs.db)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ReservedBlock.reserveAddress got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("ReservedBlock.reserveAddress error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestReservedBlock_getNextAddress(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *ReservedBlock
		inspect func(r *ReservedBlock, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1      net.IP
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			got1, err := receiver.getNextAddress(tArgs.db)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ReservedBlock.getNextAddress got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("ReservedBlock.getNextAddress error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func TestReservedBlock_addressIsInUse(t *testing.T) {
	type args struct {
		ip            net.IP
		usedAddresses []ReservedAddress
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *ReservedBlock
		inspect func(r *ReservedBlock, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		want1 bool
	}{
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			got1 := receiver.addressIsInUse(tArgs.ip, tArgs.usedAddresses)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ReservedBlock.addressIsInUse got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
