package main

import (
	"errors"
	"fmt"
	"github.com/ARMmaster17/Captain/shared/ampq"
	"github.com/ARMmaster17/Captain/shared/ipam"
	"github.com/ARMmaster17/Captain/shared/proxmox"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
	"time"
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

func planeDefaultOverlapBuild(tlc ampq.Plane) (ampq.Plane, error) {
	defaultConfig, err := ioutil.ReadFile("./conf/plane_default.yaml")
	if err != nil {
		log.Println(err)
		return ampq.Plane{}, errors.New("unable to read default plane config values")
	}
	var defaultPlane = ampq.Plane{}
	err = yaml.Unmarshal(defaultConfig, &defaultPlane)
	if err != nil {
		log.Println(err)
		return ampq.Plane{}, errors.New("unable to parse default plane config values")
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
	return defaultPlane, nil
}

func getAllPlaneConfig() (AllPlaneConfig, error) {
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

func getFQDNHostname(name string) (string, error) {
	allPlaneConfig, err := getAllPlaneConfig()
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to build cluster-wide plane configuration")
	}
	return fmt.Sprintf("%s.%s", name, allPlaneConfig.Domain), nil
}

func buildPlaneConfig(tlc ampq.Plane) (proxmox.MachineConfig, error) {
	planeConfig, err := planeDefaultOverlapBuild(tlc)
	if err != nil {
		log.Println(err)
		return proxmox.MachineConfig{}, errors.New("unable to build plane configuration")
	}

	allPlaneConfig, err := getAllPlaneConfig()
	if err != nil {
		log.Println(err)
		return proxmox.MachineConfig{}, errors.New("unable to build cluster-wide plane configuration")
	}
	hostname, err := getFQDNHostname(tlc.Name)
	if err != nil {
		log.Println(err)
		return proxmox.MachineConfig{}, errors.New("unable to build FQDN")
	}
	ip, err := ipam.GetIPAddress(hostname)
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
		Hostname: hostname,
		Cores: planeConfig.CPU,
		Memory: planeConfig.RAM,
		Swap: planeConfig.RAM,
		Nameservers: strings.Join(allPlaneConfig.Nameservers, " "),
		Storage: allPlaneConfig.DiskStore,
		RootFS: fmt.Sprintf("%s:%d", allPlaneConfig.DiskStore, planeConfig.Storage),
		OnBoot: 1,
		Unprivileged: 1,
		Start: 1,
		SSH: string(publicKey),
	}, nil
}

func makePlane(plane ampq.Plane) (string, error) {
	machineConfig, err := buildPlaneConfig(plane)
	if err != nil {
		log.Println(err)
		return "", errors.New("an error occurred while building the plane configuration")
	}
	vmid, err := proxmox.CreateLxcContainer(machineConfig)
	if err != nil {
		log.Println(err)
		return "", errors.New("an error occurred while building the plane")
	}
	return vmid, nil
}

func destroyPlane(plane ampq.Plane) error {
	hostname, err := getFQDNHostname(plane.Name)
	if err != nil {
		log.Println(err)
		return errors.New("unable to build FQDN")
	}
	err = ipam.ReleaseIPAddress(hostname)
	if err != nil {
		log.Println(err)
		return errors.New("unable to release IP address")
	}
	vmid, err := proxmox.GetVmidFromHostname(hostname)
	err = proxmox.StopLxcContainer(vmid)
	if err != nil {
		log.Println(err)
		return errors.New("unable to stop container")
	}
	time.Sleep(30 * time.Second)
	err = proxmox.DestroyLxcContainer(vmid)
	if err != nil {
		log.Println(err)
		return errors.New("unable to destroy container")
	}
	return nil
}
