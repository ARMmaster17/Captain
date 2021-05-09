package drivers

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestProviderDriverManagerDriverLookupByYAMLTagDummy(t *testing.T) {
	driver, err := driverLookupByYAMLTag("dummy")
	assert.NoError(t, err)
	assert.Equal(t, "dummy", driver.GetYAMLTag())
	assert.Equal(t, "dummy", driver.GetCUIDPrefix())
}

func TestProviderDriverManagerDriverLookupByYAMLTagProxmoxLxc(t *testing.T) {
	driver, err := driverLookupByYAMLTag("proxmoxlxc")
	assert.NoError(t, err)
	assert.Equal(t, "proxmoxlxc", driver.GetYAMLTag())
	assert.Equal(t, "proxmox.lxc", driver.GetCUIDPrefix())
}

func TestProviderDriverManagerDriverLookupByYAMLTagInvalid(t *testing.T) {
	driver, err := driverLookupByYAMLTag("invalid")
	assert.Error(t, err)
	assert.Nil(t, driver)
}

func TestProviderDriverManagerDriverLookupByCUIDPrefixDummy(t *testing.T) {
	driver, err := driverLookupByCUIDPrefix("dummy")
	assert.NoError(t, err)
	assert.Equal(t, "dummy", driver.GetYAMLTag())
	assert.Equal(t, "dummy", driver.GetCUIDPrefix())
}

func TestProviderDriverManagerDriverLookupByCUIDPrefixProxmoxLxc(t *testing.T) {
	driver, err := driverLookupByCUIDPrefix("proxmox.lxc")
	assert.NoError(t, err)
	assert.Equal(t, "proxmoxlxc", driver.GetYAMLTag())
	assert.Equal(t, "proxmox.lxc", driver.GetCUIDPrefix())
}

func TestProviderDriverManagerDriverLookupByCUIDPrefixInvalid(t *testing.T) {
	driver, err := driverLookupByCUIDPrefix("invalid")
	assert.Error(t, err)
	assert.Nil(t, driver)
}

func TestProviderDriverManagerGetDestroyDriverDummy(t *testing.T) {
	driver, err := getDestroyDriver("dummy:0")
	assert.NoError(t, err)
	assert.Equal(t, "dummy", driver.GetYAMLTag())
	assert.Equal(t, "dummy", driver.GetCUIDPrefix())
}

func TestProviderDriverManagerGetDestroyDriverProxmoxLxc(t *testing.T) {
	driver, err := getDestroyDriver("proxmox.lxc:0")
	assert.NoError(t, err)
	assert.Equal(t, "proxmoxlxc", driver.GetYAMLTag())
	assert.Equal(t, "proxmox.lxc", driver.GetCUIDPrefix())
}

func TestProviderDriverManagerGetDestroyDriverInvalid(t *testing.T) {
	driver, err := getDestroyDriver("invalid:0")
	assert.Error(t, err)
	assert.Nil(t, driver)
}

func TestProviderDriverManagerGetActiveBuildDriverDummy(t *testing.T) {
	err := helperSetupConfigFile("config_dummy_only.yaml")
	assert.NoError(t, err)
	driver, err := getActiveBuildDriver()
	assert.NoError(t, err)
	assert.NotNil(t, driver)
	assert.Equal(t, "dummy", driver.GetYAMLTag())
	assert.Equal(t, "dummy", driver.GetCUIDPrefix())
}

func helperSetupConfigFile(configFile string) error {
	viper.Reset()
	_ = os.Remove("/etc/captain/config.yaml")
	input, err := ioutil.ReadFile(fmt.Sprintf("testing/%s", configFile))
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