package main

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

func TestGenerateConfigFile(t *testing.T) {
	_ = os.Remove("/etc/captain/atc/config.yaml")
	_, err := os.Create("/etc/captain/atc/config.yaml")
	assert.NoError(t, err)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/captain/atc")
	err = generateConfigFile()
	assert.NoError(t, err)
	info, err := os.Stat("/etc/captain/atc/config.yaml")
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Greater(t, info.Size(), int64(0))
}

func TestInitLogging(t *testing.T) {
	initLogging()
}

func Test_getLogLevel(t *testing.T) {
	type args struct {
		level int
	}
	tests := []struct {
		name string
		args func(t *testing.T) args
		want1 zerolog.Level
	}{
		{
			name: "TestLevel5",
			args: func(t *testing.T) args {
				return args{
					level: 5,
				}
			},
			want1: zerolog.PanicLevel,
		},
		{
			name: "TestLevel4",
			args: func(t *testing.T) args {
				return args{
					level: 4,
				}
			},
			want1: zerolog.FatalLevel,
		},
		{
			name: "TestLevel3",
			args: func(t *testing.T) args {
				return args{
					level: 3,
				}
			},
			want1: zerolog.ErrorLevel,
		},
		{
			name: "TestLevel2",
			args: func(t *testing.T) args {
				return args{
					level: 2,
				}
			},
			want1: zerolog.WarnLevel,
		},
		{
			name: "TestLevel1",
			args: func(t *testing.T) args {
				return args{
					level: 1,
				}
			},
			want1: zerolog.InfoLevel,
		},
		{
			name: "TestLevel0",
			args: func(t *testing.T) args {
				return args{
					level: 0,
				}
			},
			want1: zerolog.DebugLevel,
		},
		{
			name: "TestLevel-1",
			args: func(t *testing.T) args {
				return args{
					level: -1,
				}
			},
			want1: zerolog.TraceLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := getLogLevel(tArgs.level)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getCaptainClient got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
