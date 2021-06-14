package imagestore

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestGetProviderSpecificImageConfiguration(t *testing.T) {
	type args struct {
		driverYamlTag string
		imagetype     string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "ImageStoreValidInputs",
			args: func(t *testing.T) args {
				return args{
					driverYamlTag: "dummy",
					imagetype:     "debian-10",
				}
			},
			want1:   "notrealpath",
			wantErr: false,
		},
		{
			name: "ImageStoreBadPath",
			args: func(t *testing.T) args {
				return args{
					driverYamlTag: "dummy",
					imagetype:     "fakeos",
				}
			},
			want1:   "",
			wantErr: true,
		},
		{
			name: "ImageStoreBadDriver",
			args: func(t *testing.T) args {
				return args{
					driverYamlTag: "notrealdriver",
					imagetype:     "debian-10",
				}
			},
			want1:   "",
			wantErr: true,
		},
		{
			name: "ImageStoreEmptyDriver",
			args: func(t *testing.T) args {
				return args{
					driverYamlTag: "",
					imagetype:     "debian-10",
				}
			},
			want1:   "",
			wantErr: true,
		},
		{
			name: "ImageStoreEmptyPath",
			args: func(t *testing.T) args {
				return args{
					driverYamlTag: "dummy",
					imagetype:     "",
				}
			},
			want1:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)
			helperSetupConfigFile("config_dummy_only.yaml")
			got1, err := GetProviderSpecificImageConfiguration(tArgs.driverYamlTag, tArgs.imagetype)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetProviderSpecificImageConfiguration got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetProviderSpecificImageConfiguration error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func helperSetupConfigFile(configFile string) error {
	viper.Reset()
	_ = os.Remove("/etc/captain/atc/config.yaml")
	input, err := ioutil.ReadFile(fmt.Sprintf("../testing/%s", configFile))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/captain/atc/config.yaml", input, 0644)
	if err != nil {
		return err
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/captain/atc")
	return viper.ReadInConfig()
}
