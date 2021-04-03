package captain

import (
	"errors"
	"fmt"
	"github.com/ARMmaster17/Captain/pkg/ipam"
	"github.com/ARMmaster17/Captain/pkg/proxmox"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

type AllPlaneConfig struct {
	Nameservers	[]string	`yaml:"nameservers"`
	Template	string		`yaml:"template"`
	Bridge		string		`yaml:"bridge"`
	Gateway		string		`yaml:"gateway"`
	MTU			int			`yaml:"mtu"`
	Domain		string		`yaml:"domain"`
	DiskStore	string		`yaml:"disk_store"`
	StartOnBoot	int			`yaml:"autostart"`
	CIDR		int			`yaml:"cidr"`
}

func NewAllPlaneConfig() (AllPlaneConfig, error) {
	allPlaneConfigFile, err := ioutil.ReadFile("/etc/captain/builder/conf/plane_config.yaml")
	if err != nil {
		log.Println(err)
		return AllPlaneConfig{}, errors.New("unable to import cluster-wide plane configuration")
	}
	var allPlaneConfig = AllPlaneConfig{}
	err = yaml.Unmarshal(allPlaneConfigFile, &allPlaneConfig)
	if err != nil {
		log.Println(err)
		return AllPlaneConfig{}, errors.New("unable to parse cluster-wide plane configuration")
	}
	return allPlaneConfig, nil
}

func PlaneInjectDefaults(tlc *Plane) error {
	defaultConfig, err := ioutil.ReadFile("/etc/captain/builder/conf/plane_default.yaml")
	if err != nil {
		log.Println(err)
		return errors.New("unable to read default plane config values")
	}
	var defaultPlane = Plane{}
	err = yaml.Unmarshal(defaultConfig, &defaultPlane)
	if err != nil {
		log.Println(err)
		return errors.New("unable to parse default plane config values")
	}
	if tlc.Name != "" {
		defaultPlane.Name = tlc.Name
	}
	if tlc.CPU != 0 {
		defaultPlane.CPU = tlc.CPU
	}
	if tlc.RAM != 0 {
		defaultPlane.RAM = tlc.RAM
	}
	if tlc.Storage != 0 {
		defaultPlane.Storage = tlc.Storage
	}
	return nil
}

func BuildPlaneConfig(tlc *Plane) (proxmox.MachineConfig, error) {
	err := PlaneInjectDefaults(tlc)
	if err != nil {
		log.Println(err)
		return proxmox.MachineConfig{}, errors.New("unable to build plane configuration")
	}

	allPlaneConfig, err := NewAllPlaneConfig()
	if err != nil {
		log.Println(err)
		return proxmox.MachineConfig{}, errors.New("unable to build cluster-wide plane configuration")
	}
	hostname, err := tlc.GetFQDNHostname()
	if err != nil {
		log.Println(err)
		return proxmox.MachineConfig{}, errors.New("unable to build FQDN")
	}
	ipamAPI, err := ipam.NewIPAM()
	if err != nil {
		log.Println(err)
		return proxmox.MachineConfig{}, errors.New("unable to contact IPAM API")
	}
	ip, err := ipamAPI.IPCreateFromFirstFree(ipam.Hostname(hostname))
	if err != nil {
		log.Println(err)
		return proxmox.MachineConfig{}, errors.New("unable to obtain next available IP address")
	}
	publicKey, err := ioutil.ReadFile("/etc/captain/builder/conf/key.pub")
	if err != nil {
		log.Println(err)
		return proxmox.MachineConfig{}, errors.New("unable to read cluster-wide public key")
	}
	return proxmox.MachineConfig{
		OsTemplate: allPlaneConfig.Template,
		Net0: fmt.Sprintf("name=eth0,bridge=%s,ip=" + "%s/%d" + ",gw=%s,firewall=0,mtu=%d",
			allPlaneConfig.Bridge,
			ip,
			allPlaneConfig.CIDR,
			allPlaneConfig.Gateway,
			allPlaneConfig.MTU,
		),
		Hostname: string(hostname),
		Cores: tlc.CPU,
		Memory: tlc.RAM,
		Swap: tlc.RAM,
		Nameservers: strings.Join(allPlaneConfig.Nameservers, " "),
		Storage: allPlaneConfig.DiskStore,
		RootFS: fmt.Sprintf("%s:%d", allPlaneConfig.DiskStore, tlc.Storage),
		OnBoot: 1,
		Unprivileged: 1,
		Start: 1,
		SSH: string(publicKey),
	}, nil
}