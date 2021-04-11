package main

// Borrowed largely from https://github.com/Telmate/proxmox-api-go/blob/master/main.go

import (
	"crypto/tls"
	"fmt"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"time"
)

func ProxmoxAdapterConnect() (*proxmox.Client, error) {
	tlsConf := &tls.Config{InsecureSkipVerify: true}
	c, _ := proxmox.NewClient(os.Getenv("CAPTAIN_PROXMOX_URL"), nil, tlsConf, 300)
	err := c.Login(os.Getenv("CAPTAIN_PROXMOX_USER"), os.Getenv("CAPTAIN_PROXMOX_PASSWORD"), "")
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate with Proxmox cluster with error: %w", err)
	}
	return c, nil
}

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
			"ip": "10.1.0.200/16",
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

	nextId, err := client.GetNextID(0)
	if err != nil {
		return fmt.Errorf("unable to retreive next available VMID with error: %w", err)
	}
	vmr := proxmox.NewVmRef(nextId)
	vmr.SetNode(defaults.Proxmox.DefaultNode)
	err = config.CreateLxc(vmr, client)
	if err != nil {
		return fmt.Errorf("unable to create LXC container with error: %w", err)
	}
	p.ProxmoxIdentifier = vmr.VmId()
	return nil
}

func ProxmoxDestroyLxc(client *proxmox.Client, p *Plane) error {
	vmr, err := client.GetVmRefByName(p.getFQDN())
	if err != nil {
		return fmt.Errorf("unable to obtain reference to underlying LXC container for plane %s: %w", p.getFQDN(), err)
	}
	_, err = client.StopVm(vmr)
	if err != nil {
		return fmt.Errorf("unable to stop LXC container: %w", err)
	}
	time.Sleep(10 * time.Second)
	_, err = proxmoxOverrideDeleteVmParams(client, vmr, map[string]interface{}{
		"force": 1,
		"purge": 1,
	})
	if err != nil {
		return fmt.Errorf("unable to delete LXC container for plane %s: %w", p.getFQDN(), err)
	}
	return nil
}

func proxmoxOverrideDeleteVmParams(c *proxmox.Client, vmr *proxmox.VmRef, params map[string]interface{}) (string, error) {
	err := c.CheckVmRef(vmr)
	if err != nil {
		return "", err
	}
	fmt.Printf("VMR TYPE: %s", vmr.GetVmType())
	fmt.Printf("VMR ID: %d", vmr.VmId())
	fmt.Printf("VMR NODE: %s", vmr.Node())
	reqBody := proxmox.ParamsToBody(params)
	url := fmt.Sprintf("/nodes/%s/%s/%d", vmr.Node(), vmr.GetVmType(), vmr.VmId())
	var taskResponse map[string]interface{}
	session, err := proxmox.NewSession(os.Getenv("CAPTAIN_PROXMOX_URL"), nil, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return "", fmt.Errorf("INTERNAL ERROR: %w", err)
	}
	err = session.Login(os.Getenv("CAPTAIN_PROXMOX_USER"), os.Getenv("CAPTAIN_PROXMOX_PASSWORD"), "")
	if err != nil {
		return "", fmt.Errorf("INTERNAL ERROR: %w", err)
	}
	resp, err := session.RequestJSON("DELETE", url, nil, nil, &reqBody, &taskResponse)
	if err != nil {
		return "", fmt.Errorf("INTERNAL ERROR: %w", err)
	}
	fmt.Printf("DELETE SC: %d", resp.StatusCode)
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("INTERNAL ERROR: %w", err)
	}
	fmt.Println(string(rbody))
	return "", nil
}
