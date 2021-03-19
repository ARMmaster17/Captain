package captain

import (
	"errors"
	"fmt"
	"github.com/ARMmaster17/Captain/ipam"
	"github.com/ARMmaster17/Captain/proxmox"
	"log"
)

const MaxCpu = 8192
const MaxRam = 3072000
const MaxStorage = 8000

type Plane struct {
	Name	string	`yaml:"name"`
	CPU		int		`yaml:"cpu"`
	RAM		int		`yaml:"ram"`
	Storage	int		`yaml:"storage"`
}

func NewPlane(name string, cpu int, ram int, storage int) (Plane, error) {
	// TODO: Validate name for invalid characters
	if cpu <= 0 || cpu > MaxCpu {
		return Plane{}, errors.New("invalid CPU parameter " + string(cpu))
	}
	if ram <= 0 || ram > MaxRam {
		return Plane{}, errors.New("invalid RAM parameter " + string(ram))
	}
	if storage <= 0 || storage > MaxStorage {
		return Plane{}, errors.New("invalid Storage parameter " + string(storage))
	}
	return Plane{
		Name: name,
		CPU: cpu,
		RAM: ram,
		Storage: storage,
	}, nil
}

func (p *Plane) Create(machineConfig proxmox.MachineConfig) (string, error) {
	proxmoxAPI, err := proxmox.NewProxmox()
	if err != nil {
		return "", errors.New("unable to contact Proxmox API")
	}
	lxc, err := proxmoxAPI.LXCCreate(machineConfig)
	if err != nil {
		return "", errors.New("unable to create LXC container")
	}
	return lxc.VMID, nil
}

func (p *Plane) Destroy() error {
	ipamAPI, err := ipam.NewIPAM()
	if err != nil {
		return errors.New("unable to contact IPAM API")
	}
	hostname, err := p.GetFQDNHostname()
	if err != nil {
		log.Println(err)
		return errors.New("unable to build FQDN")
	}
	h := ipam.Hostname(hostname)
	err = h.Delete(ipamAPI)
	if err != nil {
		log.Println(err)
		return errors.New("unable to release IP address")
	}
	proxmoxAPI, err := proxmox.NewProxmox()
	if err != nil {
		log.Println(err)
		return errors.New("unable to contact Proxmox API")
	}
	lxc, err := proxmoxAPI.GetLXCFromHostname(h)
	if err != nil {
		log.Println(err)
		return errors.New("unable to lookup VMID from hostname")
	}
	err = lxc.Stop(proxmoxAPI)
	if err != nil {
		log.Println(err)
		return errors.New("unable to stop container")
	}
	err = lxc.Destroy(proxmoxAPI)
	if err != nil {
		log.Println(err)
		return errors.New("unable to destroy container")
	}
	return nil
}

func (p *Plane) GetFQDNHostname() (ipam.Hostname, error) {
	//allPlaneConfig, err := getAllPlaneConfig()
	//if err != nil {
	//	log.Println(err)
	//	return "", errors.New("unable to build cluster-wide plane configuration")
	//}
	// TODO: Fix somehow
	return ipam.Hostname(fmt.Sprintf("%s.%s", p.Name, /*allPlaneConfig.Domain*/ "firecore.lab")), nil
}
