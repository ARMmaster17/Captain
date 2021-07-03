package proxmox

import (
	"crypto/tls"
	"fmt"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"strings"
)

var (
	proxmoxClient *proxmox.Client
)

// Fake method wrappers for mocking in unit tests.
var proxmoxNewClientFunc = proxmox.NewClient
var proxmoxClientLoginFunc = proxmoxClient.Login
var proxmoxCreateLxcFunc = func(config *proxmox.ConfigLxc, vmr *proxmox.VmRef) error {
	return config.CreateLxc(vmr, proxmoxClient)
}

func NewClient(apiUrl string, forceSSL bool, taskTimeout int) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !forceSSL,
	}
	var err error
	proxmoxClient, err = proxmoxNewClientFunc(apiUrl, nil, tlsConfig, taskTimeout)
	return err
}

func Login(username string, password string) error {
	if !strings.Contains(username, "@") {
		return fmt.Errorf("username should be of the format user@realm")
	}
	if len(password) == 0 {
		return fmt.Errorf("password was not provided")
	}
	return proxmoxClientLoginFunc(username, password, "")
}

func CreateLxc(config *proxmox.ConfigLxc, vmr *proxmox.VmRef) error {
	return proxmoxCreateLxcFunc(config, vmr)
}
