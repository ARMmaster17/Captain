package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type defaultplane struct {
	PublicKey string `yaml:"publickey"`
	Network defaultplanenetwork `yaml:"network"`
	Proxmox defaultplaneproxmox `yaml:"proxmox"`
}

type defaultplanenetwork struct {
	Nameservers string `yaml:"nameservers"`
	SearchDomain string `yaml:"searchdomain"`
	Gateway string `yaml:"gateway"`
	MTU string `yaml:"mtu"`
}

type defaultplaneproxmox struct {
	Image string `yaml:"image"`
	PublicNetwork string `yaml:"publicnetwork"`
	DiskStorage string `yaml:"diskstorage"'`
	DefaultNode string `yaml:"defaultnode"`
}

func getPlaneDefaults() (*defaultplane, error) {
	defaultFile, err := ioutil.ReadFile("defaults.yaml")
	if err != nil {
		return nil, fmt.Errorf("unable to read defaults.yaml: %w", err)
	}
	var defaults defaultplane
	err = yaml.Unmarshal(defaultFile, &defaults)
	if err != nil {
		return nil, fmt.Errorf("unable to parse YAML in defaults.yaml: %w", err)
	}
	return &defaults, nil
}
