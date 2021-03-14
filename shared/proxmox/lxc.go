package proxmox

import (
	"errors"
	"fmt"
	"github.com/ARMmaster17/Captain/shared/ipam"
	"github.com/tidwall/gjson"
	"log"
)

type LXC struct {
	VMID string `json:"vmid"`
}

// Creates an LXC container with the specified parameters. Waits until
// task returns complete before returning.
func (p *Proxmox) LXCCreate(config MachineConfig) (*LXC, error) {
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
	vmid, err := p.getNextVmid()
	if err != nil {
		log.Println(err)
		return &LXC{}, &Task{}, errors.New("unable to obtain next available VMID")
	}
	config.VMID = vmid
	node, err := p.GetLowestRAMUtilizationNode()
	if err != nil {
		log.Println("unable to find available Proxmox node")
	}
	body, err := p.doPostUrlEncode("nodes/" + node.Name + "/lxc", config)
	if err != nil {
		log.Println(err)
		return &LXC{}, &Task{}, errors.New("failed to create new LXC container")
	}
	task, err := NewTask(body, node.Name)
	if err != nil {
		log.Println(err)
		return &LXC{}, &Task{}, errors.New("failed to obtain status of LXC container creation")
	}
	return &LXC{
		VMID: vmid,
	}, &task, nil
}

func (l *LXC) Destroy(p *Proxmox) error {
	node, err := p.FindNodeWithVMID(l.VMID)
	if err != nil {
		log.Println(err)
		return errors.New("container to be destroyed does not exist")
	}
	_, err = p.doDelete(fmt.Sprintf("nodes/%s/lxc/%s", node.Name, l.VMID))
	if err != nil {
		log.Println(err)
		return errors.New("unable to destroy LXC container")
	}
	return nil
}

func (l *LXC) Stop(p *Proxmox) error {
	node, err := p.FindNodeWithVMID(l.VMID)
	if err != nil {
		log.Println(err)
		return errors.New("container to be destroyed does not exist")
	}
	body, err := p.doPost(fmt.Sprintf("nodes/%s/lxc/%s/status/shutdown", node.Name, l.VMID), nil)
	if err != nil {
		log.Println(err)
		return errors.New("unable to stop LXC container")
	}
	task, err := NewTask(node.Name, body)
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
	nodeList, err := p.getNodes()
	if err != nil {
		log.Println(err)
		return LXC{}, errors.New("unable to get list of nodes")
	}
	for _, node := range nodeList {
		url := fmt.Sprintf("nodes/%s/lxc", node.Name)
		body, err := p.doGet(url)
		if err != nil {
			log.Println(err)
			return LXC{}, errors.New("unable to query API for hostname/VMID association")
		}
		vmid := gjson.Get(body, "data.#(name==\""+string(hostname)+"\").vmid").String()
		if vmid != "" {
			return LXC{
				VMID: vmid,
			}, nil
		}
	}
	return LXC{}, errors.New("no VMID found for hostname " + string(hostname))
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