package main

// Borrowed largely from https://github.com/Telmate/proxmox-api-go/blob/master/main.go

import (
	"crypto/tls"
	"fmt"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
)

func ProxmoxAdapterConnect() (*proxmox.Client, error) {
	tlsConf := &tls.Config{InsecureSkipVerify: true}
	c, _ := proxmox.NewClient("URL", nil, tlsConf, 300)
	err := c.Login(os.Getenv("PM_USER"), os.Getenv("PM_PASS"), "")
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate with Proxmox cluster with error: %w", err)
	}
	return c, nil
}

func ProxmoxBuildLxc(db *gorm.DB, client *proxmox.Client, p *Plane) error {
	publicKey, err := ioutil.ReadFile("/etc/captain/builder/conf/key.pub")
	if err != nil {
		return fmt.Errorf("unable to read cluster-wide public key with error: %w", err)
	}

	// When https://github.com/Telmate/proxmox-api-go/pull/114 merges, this can be split up into smaller functions.
	/*config := proxmox.ConfigLxc{
		Ostemplate:         "pve-img:vztmpl/debian-10-standard_10.7-1_amd64.tar.gz",
		Arch:               "amd64",
		CMode:              "tty",
		Console:            true,
		Cores:              p.Formation.CPU,
		CPULimit:           0,
		CPUUnits:           1024,
		Description:        "Managed by the Captain stack",
		Hostname:           p.getFQDN(),
		Memory:             p.Formation.RAM,
		Nameserver:         "10.1.0.4",
		Networks:           nil,
		OnBoot:             true,
		Protection:         false,
		SearchDomain:       "firecore.lab",
		SSHPublicKeys:      string(publicKey),
		Start:              true,
		Storage:            "pve-storage",
		Swap:               p.Formation.RAM,
		Template:           false,
		Tty:                2,
		Unprivileged:       true,
	}*/
	config := proxmox.NewConfigLxc()
	config.Ostemplate = "pve-img:vztmpl/debian-10-standard_10.7-1_amd64.tar.gz"
	config.Arch = "amd64"
	config.CMode = "tty"
	config.Console = true
	config.Cores = p.Formation.CPU
	config.CPULimit = 0
	config.CPUUnits = 1024
	config.Description = "Managed by the Captain stack"
	config.Hostname = p.getFQDN()
	config.Memory = p.Formation.RAM
	config.Nameserver = "10.1.0.4 8.8.8.8"
	config.Networks = nil // TODO: Figure out how this is supposed to work
	config.OnBoot = true
	config.Protection = false
	config.SearchDomain = "firecore.lab"
	config.SSHPublicKeys = string(publicKey)
	config.Start = true
	config.Storage = "pve-storage"
	config.Swap = p.Formation.RAM
	config.Template = false
	config.Tty = 2
	config.Unprivileged = true

	nextId, err := client.GetNextID(0)
	if err != nil {
		return fmt.Errorf("unable to retreive next available VMID with error: %w", err)
	}
	vmr := proxmox.NewVmRef(nextId)
	vmr.SetNode("pxvh1")
	err = config.CreateLxc(vmr, client)
	if err != nil {
		return fmt.Errorf("unable to create LXC container with error: %w", err)
	}
	p.ProxmoxIdentifier = vmr.VmId()
	result := db.Save(p)
	if result.Error != nil {
		return fmt.Errorf("unable to save new VMID of plane in DB: %w", err)
	}
	return nil
}

func ProxmoxDestroyLxc(client *proxmox.Client, p *Plane) error {
	vmr, err := client.GetVmRefByName(p.getFQDN())
	if err != nil {
		return fmt.Errorf("unable to obtain reference to underlying LXC container for plane %s: %w", p.getFQDN(), err)
	}
	_, err = client.DeleteVm(vmr)
	// TODO: Should probably parse the exit string.
	if err != nil {
		return fmt.Errorf("unable to delete LXC container for plane %s: %w", p.getFQDN(), err)
	}
	return nil
}
