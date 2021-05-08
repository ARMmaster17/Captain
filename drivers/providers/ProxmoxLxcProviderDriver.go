package providers

import (
	"crypto/tls"
	"fmt"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"os"
)

type ProxmoxLxcProviderDriver struct {
	client *proxmox.Client
}

func (d *ProxmoxLxcProviderDriver) Connect() error {
	tlsConf := &tls.Config{InsecureSkipVerify: true}
	d.client, _ = proxmox.NewClient(os.Getenv("CAPTAIN_PROXMOX_URL"), nil, tlsConf, 300)
	return d.client.Login(os.Getenv("CAPTAIN_PROXMOX_USER"), os.Getenv("CAPTAIN_PROXMOX_PASSWORD"), "")
}

func (d *ProxmoxLxcProviderDriver) BuildPlane(p *GenericPlane) error {
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

func (d *ProxmoxLxcProviderDriver) DestroyPlane(p *GenericPlane) error {

}