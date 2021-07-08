package models

import (
	"testing"
)

func TestValidate(t *testing.T) {
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