package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

type Plane struct {
	Name	string	`yaml:"name"`
	CPU		int		`yaml:"cpu"`
	RAM		int		`yaml:"ram"`
	Storage	int		`yaml:"storage"`
}

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

func planeDefaultOverlapBuild(tlc Plane) (Plane, error) {
	defaultConfig, err := ioutil.ReadFile("./conf/plane_default.yaml")
	if err != nil {
		log.Println(err)
		return Plane{}, errors.New("unable to read default plane config values")
	}
	var defaultPlane = Plane{}
	err = yaml.Unmarshal(defaultConfig, &defaultPlane)
	if err != nil {
		log.Println(err)
		return Plane{}, errors.New("unable to parse default plane config values")
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
	allPlaneConfigFile, err := ioutil.ReadFile("./conf/plane_config.yaml")
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

func buildPlaneConfig(tlc Plane) (MachineConfig, error) {
	planeConfig, err := planeDefaultOverlapBuild(tlc)
	if err != nil {
		log.Println(err)
		return MachineConfig{}, errors.New("unable to build plane configuration")
	}

	allPlaneConfig, err := getAllPlaneConfig()
	if err != nil {
		log.Println(err)
		return MachineConfig{}, errors.New("unable to build cluster-wide plane configuration")
	}
	hostname := fmt.Sprintf("%s.%s", planeConfig.Name, allPlaneConfig.Domain)
	ip, err := getIPAddress(hostname)
	publicKey, err := ioutil.ReadFile("./conf/key.pub")
	return MachineConfig{
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

func makePlane(plane Plane) (string, error) {
	machineConfig, err := buildPlaneConfig(plane)
	if err != nil {
		log.Println(err)
		return "", errors.New("an error occurred while building the plane configuration")
	}
	vmid, err := createLxcContainer(machineConfig)
	if err != nil {
		log.Println(err)
		return "", errors.New("an error occurred while building the plane")
	}
	return vmid, nil
}
