package main

// Borrowed largely from https://github.com/Telmate/proxmox-api-go/blob/master/main.go

import (
	"crypto/tls"
	"fmt"
	"github.com/ARMmaster17/proxmox-api-go/proxmox"
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

func ProxmoxBuildLxc(p *Plane) (proxmox.ConfigLxc, error) {
	publicKey, err := ioutil.ReadFile("/etc/captain/builder/conf/key.pub")
	if err != nil {
		return proxmox.ConfigLxc{}, fmt.Errorf("unable to read cluster-wide public key with error: %w", err)
	}
	config := proxmox.ConfigLxc{
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
	}

	return config, nil
}

func ProxmoxCreateLxc(client *proxmox.Client, config proxmox.ConfigLxc) (int, error) {
	nextId, err := client.GetNextID(0)
	if err != nil {
		return 0, fmt.Errorf("unable to retreive next available VMID with error: %w", err)
	}
	vmr := proxmox.NewVmRef(nextId)
	vmr.SetNode("pxvh1")
	err = config.CreateLxc(vmr, client)
	if err != nil {
		return 0, fmt.Errorf("unable to create LXC container with error: %w", err)
	}
	return vmr.VmId(), nil
}
