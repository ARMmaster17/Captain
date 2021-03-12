package proxmox

import (
	"errors"
	"fmt"
	"github.com/ARMmaster17/Captain/shared/ipam"
	"github.com/tidwall/gjson"
	"log"
	"os"
)

type LXC struct {
	VMID string
}

func (p *Proxmox) LXCCreate(config MachineConfig) (*LXC, error) {
	p.Authenticate()
	vmid, err := p.getNextVmid()
	if err != nil {
		log.Println(err)
		return &LXC{}, errors.New("unable to obtain next available VMID")
	}
	config.VMID = vmid
	node := os.Getenv("PROXMOX_DEFAULT_NODE")
	_, err = p.Client.PostUrlEncode("nodes/" + node + "/lxc", config)
	if err != nil {
		log.Println(err)
		return &LXC{}, errors.New("failed to create new LXC container")
	}
	return &LXC{
		VMID: vmid,
	}, nil
}

func (l *LXC) Destroy(p *Proxmox) error {
	p.Authenticate()
	node := os.Getenv("PROXMOX_DEFAULT_NODE")
	_, err := p.Client.Delete("nodes/" + node + "/lxc/" + l.VMID)
	if err != nil {
		log.Println(err)
		return errors.New("unable to destroy LXC container")
	}
	return nil
}

func (l *LXC) Stop(p *Proxmox) error {
	p.Authenticate()
	node := os.Getenv("PROXMOX_DEFAULT_NODE")
	_, err := p.Client.Post(fmt.Sprintf("nodes/%s/lxc/%s/status/shutdown", node, l.VMID), nil)
	if err != nil {
		log.Println(err)
		return errors.New("unable to stop LXC container")
	}
	return nil
}

func (p *Proxmox) GetLXCFromHostname(hostname ipam.Hostname) (LXC, error) {
	err := p.Authenticate()
	if err != nil {
		log.Println(err)
		return LXC{}, errors.New("unable to query API for hostname/VMID association")
	}
	url := fmt.Sprintf("nodes/%s/lxc", os.Getenv("PROXMOX_DEFAULT_NODE"))
	body, err := p.Client.Get(url)
	if err != nil {
		log.Println(err)
		return LXC{}, errors.New("unable to query API for hostname/VMID association")
	}
	vmid := gjson.Get(string(body), "data.#(name==\"" + string(hostname) + "\").vmid").String()
	if vmid == "" {
		return LXC{}, errors.New("no VMID found for hostname " + string(hostname))
	}
	return LXC{
		VMID: vmid,
	}, nil
}

type MachineConfig struct {
	VMID			string	`schema:"vmid"`
	OsTemplate		string	`schema:"ostemplate"`
	Net0			string	`schema:"net0"`
	Hostname		string	`schema:"hostname"`
	Cores			int		`schema:"cores"`
	Memory			int		`schema:"memory"`
	Swap			int		`schema:"swap"`
	Nameservers		string	`schema:"nameserver"`
	Storage			string	`schema:"storage"`
	RootFS			string	`schema:"rootfs"`
	OnBoot			int		`schema:"onboot"`
	Unprivileged	int		`schema:"unprivileged"`
	Start			int		`schema:"start"`
	SSH				string	`schema:"ssh-public-keys"`
}

type NetworkConfig struct {
	Name		string	`schema:"name"`
	Bridge		string	`schema:"bridge"`
	IP			string	`schema:"ip"`
	Gateway		string	`schema:"gw"`
	Firewall	int		`schema:"firewall"`
	MTU			int		`schema:"mtu"`
}