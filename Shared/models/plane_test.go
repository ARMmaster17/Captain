package models

import (
	"github.com/ARMmaster17/Captain/Shared/framework"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
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

func Test_PlaneCreate(t *testing.T) {
	f, err := framework.NewFramework("test")
	require.NoError(t, err)
	mockDB := HelperInitMockDB(t, f)
	plane := Plane{}
	plane.CreatedAt = time.Now()
	plane.UpdatedAt = time.Now()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("INSERT INTO `planes` (`created_at`,`updated_at`,`deleted_at`,`name`,`cpu`,`ram`,`disk`) VALUES (?,?,?,?,?,?,?)")).
		WithArgs(plane.CreatedAt, plane.UpdatedAt, plane.DeletedAt, plane.Name, plane.CPU, plane.RAM, plane.Disk).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mockDB.ExpectCommit()
	mockDB.ExpectRollback()
	err = plane.Create(f.DB)
	assert.NoError(t, err)
}

func Test_PlaneGetByID(t *testing.T) {
	f, err := framework.NewFramework("test")
	require.NoError(t, err)
	mockDB := HelperInitMockDB(t, f)
	mockDB.NewRows([]string{"id"}).AddRow("0")
	plane := Plane{}
	plane.CreatedAt = time.Now()
	plane.UpdatedAt = time.Now()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `planes` WHERE `planes`.`id` = ? AND `planes`.`deleted_at` IS NULL ORDER BY `planes`.`id` LIMIT 1")).
		WithArgs(0).WillReturnRows(mockDB.NewRows([]string{"id"}).AddRow("0"))
	mockDB.ExpectCommit()
	mockDB.ExpectRollback()
	err = plane.GetByID(f.DB, 0)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), plane.ID)
}

func Test_PlaneUpdate(t *testing.T) {
	f, err := framework.NewFramework("test")
	require.NoError(t, err)
	mockDB := HelperInitMockDB(t, f)
	plane := Plane{}
	plane.CreatedAt = time.Now()
	plane.UpdatedAt = time.Now()
	plane.ID = 1
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("UPDATE `planes` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`name`=?,`cpu`=?,`ram`=?,`disk`=? WHERE `id` = ?")).
		WithArgs(plane.CreatedAt, sqlmock.AnyArg(), plane.DeletedAt, plane.Name, plane.CPU, plane.RAM, plane.Disk, plane.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mockDB.ExpectCommit()
	mockDB.ExpectRollback()
	err = plane.Update(f.DB)
	require.NoError(t, err)
}

func Test_PlaneDelete(t *testing.T) {
	f, err := framework.NewFramework("test")
	require.NoError(t, err)
	mockDB := HelperInitMockDB(t, f)
	plane := Plane{}
	plane.CreatedAt = time.Now()
	plane.UpdatedAt = time.Now()
	plane.ID = 1
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("UPDATE `planes` SET `deleted_at`=? WHERE `planes`.`id` = ? AND `planes`.`deleted_at` IS NULL")).
		WithArgs(sqlmock.AnyArg(), plane.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mockDB.ExpectCommit()
	mockDB.ExpectRollback()
	err = plane.Delete(f.DB)
	require.NoError(t, err)
}
