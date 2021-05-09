package ImageStore

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetProviderSpecificImageConfigurationDummy(t *testing.T) {
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	path, err := GetProviderSpecificImageConfiguration("dummy", "debian-10")
	assert.NoError(t, err)
	assert.Equal(t, "notrealpath", path)
}

func TestGetProviderSpecificImageConfigurationDummyBadPath(t *testing.T) {
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	path, err := GetProviderSpecificImageConfiguration("dummy", "fakeos")
	assert.Error(t, err)
	assert.Equal(t, "", path)
}

func TestGetProviderSpecificImageConfigurationDummyBadDriver(t *testing.T) {
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	path, err := GetProviderSpecificImageConfiguration("notrealdriver", "debian-10")
	assert.Error(t, err)
	assert.Equal(t, "", path)
}

func helperSetupConfigFile(configFile string) error {
	viper.Reset()
	_ = os.Remove("/etc/captain/config.yaml")
	input, err := ioutil.ReadFile(fmt.Sprintf("../testing/%s", configFile))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/captain/config.yaml", input, 0644)
	if err != nil {
		return err
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/captain")
	return viper.ReadInConfig()
}