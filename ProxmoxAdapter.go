package main

// Borrowed largely from https://github.com/Telmate/proxmox-api-go/blob/master/main.go

import (
	"crypto/tls"
	"fmt"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"gorm.io/gorm"
	"net/http"
	"os"
)

// ProxmoxAdapterConnect Connect to Proxmox using a third party library. This provider driver handles the implentation of plane creation,
// management, and deletion and the specifics of running these operations on a Proxmox cluster.
func ProxmoxAdapterConnect() (*proxmox.Client, error) {
	tlsConf := &tls.Config{InsecureSkipVerify: true}
	c, _ := proxmox.NewClient(os.Getenv("CAPTAIN_PROXMOX_URL"), nil, tlsConf, 300)
	err := c.Login(os.Getenv("CAPTAIN_PROXMOX_USER"), os.Getenv("CAPTAIN_PROXMOX_PASSWORD"), "")
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate with Proxmox cluster with error: %w", err)
	}
	return c, nil
}

// ProxmoxBuildLxc Converts from a plane structure to a JSON object that can be passed to the Proxmox API to create a new plane
// instance. Also handles loading of defaults from defaults.yaml.
func ProxmoxBuildLxc(db *gorm.DB, client *proxmox.Client, p *Plane) error {
	defaults, err := getPlaneDefaults()
	if err != nil {
		return fmt.Errorf("unable to get default plane parameters: %w", err)
	}

	config := proxmox.NewConfigLxc()
	config.Ostemplate = defaults.Proxmox.Image
	config.Arch = "amd64"
	config.CMode = "tty"
	config.Console = true
	config.Cores = p.Formation.CPU
	config.CPULimit = 0
	config.CPUUnits = 1024
	config.Description = "Managed by the Captain stack"
	config.Hostname = p.getFQDN()
	config.Memory = p.Formation.RAM
	config.Nameserver = defaults.Network.Nameservers
	config.Networks = proxmox.QemuDevices{
		0 : {
			"name": "eth0",
			"bridge": defaults.Proxmox.PublicNetwork,
			"ip": "dhcp",
			"gw": defaults.Network.Gateway,
			"firewall": "0",
			"mtu": defaults.Network.MTU,
		},
	}
	config.OnBoot = true
	config.Protection = false
	config.SearchDomain = defaults.Network.SearchDomain
	config.SSHPublicKeys = defaults.PublicKey
	config.Start = true
	config.Storage = defaults.Proxmox.DiskStorage
	config.Swap = p.Formation.RAM
	config.Template = false
	config.Tty = 2
	config.Unprivileged = true

	nextID, err := client.GetNextID(0)
	if err != nil {
		return fmt.Errorf("unable to retreive next available VMID with error: %w", err)
	}
	vmr := proxmox.NewVmRef(nextID)
	vmr.SetNode(defaults.Proxmox.DefaultNode)
	err = config.CreateLxc(vmr, client)
	if err != nil {
		return fmt.Errorf("unable to create LXC container with error: %w", err)
	}
	p.ProxmoxIdentifier = vmr.VmId()
	return nil
}

// ProxmoxDestroyLxc Handles the destruction of the underlying LXC container. At the present moment, the third-party library has a bug
// with API calls the use the DELETE HTTP method. As such, this library calls a method that overrides the broken method
// to properly delete the container.
func ProxmoxDestroyLxc(client *proxmox.Client, p *Plane) error {
	vmr, err := client.GetVmRefByName(p.getFQDN())
	if err != nil {
		return fmt.Errorf("unable to obtain reference to underlying LXC container for plane %s: %w", p.getFQDN(), err)
	}
	_, err = client.StopVm(vmr)
	if err != nil {
		return fmt.Errorf("unable to stop LXC container: %w", err)
	}
	err = proxmoxOverrideDeleteVMParams(client, vmr)
	if err != nil {
		return fmt.Errorf("unable to delete LXC container for plane %s: %w", p.getFQDN(), err)
	}
	return nil
}

// This method replaces the method with the same name in the proxmox library because there is a bug where if you pass
// an empty struct to any DELETE endpoint, the Proxmox API returns an error. This method overrides that by passing
// nil to the underlying Session object.
func proxmoxOverrideDeleteVMParams(c *proxmox.Client, vmr *proxmox.VmRef) error {
	err := c.CheckVmRef(vmr)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("/nodes/%s/%s/%d", vmr.Node(), vmr.GetVmType(), vmr.VmId())
	var taskResponse map[string]interface{}
	session, err := proxmox.NewSession(os.Getenv("CAPTAIN_PROXMOX_URL"), nil, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return fmt.Errorf("unable to connect to the Proxmox API: %w", err)
	}
	err = session.Login(os.Getenv("CAPTAIN_PROXMOX_USER"), os.Getenv("CAPTAIN_PROXMOX_PASSWORD"), "")
	if err != nil {
		return fmt.Errorf("unable to authenticate with the Proxmox API: %w", err)
	}
	resp, err := session.RequestJSON("DELETE", url, nil, nil, nil, &taskResponse)
	if err != nil {
		return fmt.Errorf("unable to send DELETE request to the Proxmox API: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the Proxmxox API returned status code %d", resp.StatusCode)
	}
	return nil
}
