package proxmox

import (
	"crypto/tls"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_NewProxmoxClient(t *testing.T) {
	proxmoxNewClientFunc = func(apiUrl string, hclient *http.Client, tls *tls.Config, taskTimeout int) (client *proxmox.Client, err error) {
		assert.Equal(t, "http://test.example.com", apiUrl)
		assert.Equal(t, true, tls.InsecureSkipVerify)
		assert.Equal(t, 300, taskTimeout)
		return nil, nil
	}
	err := NewClient("http://test.example.com", false, 300)
	assert.NoError(t, err)
	assert.Nil(t, proxmoxClient)
}

func Test_ProxmoxClientLogin(t *testing.T) {
	proxmoxClientLoginFunc = func(username string, password string, otp string) error {
		assert.Equal(t, "username@realm", username)
		assert.Equal(t, "password", password)
		assert.Equal(t, "", otp)
		return nil
	}
	err := Login("username@realm", "password")
	assert.NoError(t, err)
}

func Test_ProxmoxClientFailsInvalidUsername(t *testing.T) {
	proxmoxClientLoginFunc = func(username string, password string, otp string) error {
		assert.Fail(t, "this method should not be invoked")
		return nil
	}
	err := Login("invalid_user", "password")
	assert.Error(t, err)
}

func Test_ProxmoxClientFailsInvalidPassword(t *testing.T) {
	proxmoxClientLoginFunc = func(username string, password string, otp string) error {
		assert.Fail(t, "this method should not be invoked")
		return nil
	}
	err := Login("username@realsm", "")
	assert.Error(t, err)
}

func Test_ProxmoxConfigCreateLxc(t *testing.T) {
	proxmoxCreateLxcFunc = func(config *proxmox.ConfigLxc, vmr *proxmox.VmRef) error {
		assert.Equal(t, proxmox.NewConfigLxc(), *config)
		assert.Nil(t, vmr)
		return nil
	}
	testConfig := proxmox.NewConfigLxc()
	err := CreateLxc(&testConfig, nil)
	assert.NoError(t, err)
}

func HelperMockProxmoxFunctions() (*int, *int, *int) {
	var loginFuncCounter int = 0
	proxmoxClientLoginFunc = func(username string, password string, otp string) error {
		loginFuncCounter++
		return nil
	}
	var newClientFuncCounter int = 0
	proxmoxNewClientFunc = func(apiUrl string, hclient *http.Client, tls *tls.Config, taskTimeout int) (client *proxmox.Client, err error) {
		newClientFuncCounter++
		return nil, nil
	}
	var createLxcFuncCounter int = 0
	proxmoxCreateLxcFunc = func(config *proxmox.ConfigLxc, vmr *proxmox.VmRef) error {
		createLxcFuncCounter++
		return nil
	}
	return &loginFuncCounter, &newClientFuncCounter, &createLxcFuncCounter
}
