package models

import (
	"testing"
)

func TestValidateName(t *testing.T) {
	type args struct {
		c CRUDObject
	}
	tests := []struct {
		name    string
		args    func(*Plane) args
		wantErr bool
	}{
		{
			name: "ValidName",
			args: func(plane *Plane) args {
				plane.Name = "testname"
				return args{
					c: plane,
				}
			},
			wantErr: false,
		},
		{
			name: "RejectEmpty",
			args: func(plane *Plane) args {
				plane.Name = ""
				return args{
					c: plane,
				}
			},
			wantErr: true,
		},
		{
			name: "RejectSymbol+",
			args: func(plane *Plane) args {
				plane.Name = "+"
				return args{
					c: plane,
				}
			},
			wantErr: true,
		},
		{
			name: "RejectSymbol.",
			args: func(plane *Plane) args {
				plane.Name = "."
				return args{
					c: plane,
				}
			},
			wantErr: true,
		},
		{
			name: "RejectSymbol&",
			args: func(plane *Plane) args {
				plane.Name = "&"
				return args{
					c: plane,
				}
			},
			wantErr: true,
		},
		{
			name: "RejectSymbol@",
			args: func(plane *Plane) args {
				plane.Name = "@"
				return args{
					c: plane,
				}
			},
			wantErr: true,
		},
		{
			name: "RejectSymbol%",
			args: func(plane *Plane) args {
				plane.Name = "%"
				return args{
					c: plane,
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.args(NewPlane()).c); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCPU(t *testing.T) {
	type args struct {
		c CRUDObject
	}
	tests := []struct {
		name    string
		args    func(*Plane) args
		wantErr bool
	}{
		{
			name: "ValidCPU1",
			args: func(plane *Plane) args {
				plane.CPU = 1
				return args{
					c: plane,
				}
			},
			wantErr: false,
		},
		{
			name: "InvalidCPU0",
			args: func(plane *Plane) args {
				plane.CPU = 0
				return args{
					c: plane,
				}
			},
			wantErr: true,
		},
		{
			name: "InvalidCPU8193",
			args: func(plane *Plane) args {
				plane.CPU = 8193
				return args{
					c: plane,
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.args(NewPlane()).c); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRAM(t *testing.T) {
	type args struct {
		c CRUDObject
	}
	tests := []struct {
		name    string
		args    func(*Plane) args
		wantErr bool
	}{
		{
			name: "ValidRAM1",
			args: func(plane *Plane) args {
				plane.RAM = 128
				return args{
					c: plane,
				}
			},
			wantErr: false,
		},
		{
			name: "InvalidMinimumRAM15",
			args: func(plane *Plane) args {
				plane.RAM = 15
				return args{
					c: plane,
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.args(NewPlane()).c); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDisk(t *testing.T) {
	type args struct {
		c CRUDObject
	}
	tests := []struct {
		name    string
		args    func(*Plane) args
		wantErr bool
	}{
		{
			name: "ValidDisk1",
			args: func(plane *Plane) args {
				plane.Disk = 8
				return args{
					c: plane,
				}
			},
			wantErr: false,
		},
		{
			name: "InvalidMinimumDisk0",
			args: func(plane *Plane) args {
				plane.Disk = 0
				return args{
					c: plane,
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.args(NewPlane()).c); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}