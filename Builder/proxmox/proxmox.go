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

// NewClient initializes a new Proxmox client object (does not attempt to connect to remote resources).
func NewClient(apiURL string, forceSSL bool, taskTimeout int) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !forceSSL,
	}
	var err error
	proxmoxClient, err = proxmoxNewClientFunc(apiURL, nil, tlsConfig, taskTimeout)
	return err
}

// Login authenticates with the specified Proxmox cluster with the provided credentials. Does not support
// OTP login at the moment.
func Login(username string, password string) error {
	if !strings.Contains(username, "@") {
		return fmt.Errorf("username should be of the format user@realm")
	}
	if len(password) == 0 {
		return fmt.Errorf("password was not provided")
	}
	return proxmoxClientLoginFunc(username, password, "")
}

// CreateLxc creates a container with the given parameters. This method is only a wrapper for the underlying Proxmox
// library, and does not perform any kind of balancing or validation.
func CreateLxc(config *proxmox.ConfigLxc, vmr *proxmox.VmRef) error {
	return proxmoxCreateLxcFunc(config, vmr)
}
