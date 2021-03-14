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

// Creates an LXC container with the specified parameters. Waits until
// task returns complete before returning.
func (p *Proxmox) LXCCreate(config MachineConfig) (*LXC, error) {
	err := p.Authenticate()
	if err != nil {
		log.Println(err)
		return &LXC{}, errors.New("unable to connect to Proxmox API")
	}
	lxc, task, err := p.LXCCreateAsync(config)
	if err != nil {
		log.Println(err)
		return &LXC{}, errors.New("unable to track async container creation")
	}
	err = task.WaitForTaskCompletion(p)
	if err != nil {
		log.Println(err)
		return &LXC{}, errors.New("unable to verify if task completed")
	}
	return lxc, nil
}

// Creates and LXC container with the specified parameters. Returns a task to track
// build completion while performing other tasks.
func (p *Proxmox) LXCCreateAsync(config MachineConfig) (*LXC, *Task, error) {
	err := p.Authenticate()
	if err != nil {
		log.Println(err)
		return &LXC{}, &Task{}, errors.New("unable to connect to Proxmox API")
	}
	vmid, err := p.getNextVmid()
	if err != nil {
		log.Println(err)
		return &LXC{}, &Task{}, errors.New("unable to obtain next available VMID")
	}
	config.VMID = vmid
	node := os.Getenv("PROXMOX_DEFAULT_NODE")
	body, err := p.Client.PostUrlEncode("nodes/" + node + "/lxc", config)
	if err != nil {
		log.Println(err)
		return &LXC{}, &Task{}, errors.New("failed to create new LXC container")
	}
	task, err := NewTask(body, node)
	if err != nil {
		log.Println(err)
		return &LXC{}, &Task{}, errors.New("failed to obtain status of LXC container creation")
	}
	return &LXC{
		VMID: vmid,
	}, &task, nil
}

func (l *LXC) Destroy(p *Proxmox) error {
	err := p.Authenticate()
	if err != nil {
		log.Println(err)
		return errors.New("unable to connect to Proxmox API")
	}
	node := os.Getenv("PROXMOX_DEFAULT_NODE")
	_, err = p.Client.Delete("nodes/" + node + "/lxc/" + l.VMID)
	if err != nil {
		log.Println(err)
		return errors.New("unable to destroy LXC container")
	}
	return nil
}

func (l *LXC) Stop(p *Proxmox) error {
	err := p.Authenticate()
	if err != nil {
		log.Println(err)
		return errors.New("unable to connect to Proxmox API")
	}
	node := os.Getenv("PROXMOX_DEFAULT_NODE")
	body, err := p.Client.Post(fmt.Sprintf("nodes/%s/lxc/%s/status/shutdown", node, l.VMID), nil)
	if err != nil {
		log.Println(err)
		return errors.New("unable to stop LXC container")
	}
	task, err := NewTask(node, body)
	if err != nil {
		log.Println(err)
		return errors.New("failed to obtain status of LXC container creation")
	}
	err = task.WaitForTaskCompletion(p)
	if err != nil {
		log.Println(err)
		return errors.New("failed wait for task completion")
	}
	return nil
}

func (p *Proxmox) GetLXCFromHostname(hostname ipam.Hostname) (LXC, error) {
	err := p.Authenticate()
	if err != nil {
		log.Println(err)
		return LXC{}, errors.New("unable to connect to Proxmox API")
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