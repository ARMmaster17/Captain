package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Top-level structure that defines valid YAML sections of a defaults.yaml file as configured by the user. This file
// contains values that should be the same for all planes, eliminating the need to enter them each time a formation
// is created.
type defaultplane struct {
	// The SSH public key used by Captain for connecting to planes for provisioning, and monitoring if SSH monitoring
	// is configured.
	PublicKey string `yaml:"publickey"`
	Network defaultplanenetwork `yaml:"network"`
	Proxmox defaultplaneproxmox `yaml:"proxmox"`
}

// Defines public network settings for a plane. In the future when SDN isolation is implemented, these settings will
// only affect the public-facing network, and does not affect the network used by Captain for provisioning and
// monitoring planes.
type defaultplanenetwork struct {
	// Space-separated IPv4 addresses to use as nameservers. Keep in mind that some Linux distributions may only use
	// the first 1, 2, or 3 values of this field.
	Nameservers string `yaml:"nameservers"`
	// Search domain to use in the DNS settings of new planes.
	SearchDomain string `yaml:"searchdomain"`
	// IPv4 address to use as the internet gateway for all new planes.
	Gateway string `yaml:"gateway"`
	// The MTU for each network adapter in a plane. Keep in mind that when using Proxmox SDN, this value must be 50 bytes
	// lower than your actual network MTU to account for VLAN header overhead. If you are using Proxmox SDN, generally
	// this value should be 1450.
	MTU string `yaml:"mtu"`
}

// Defines settings that are unique to Proxmox. Only planes that use the ProxmoxVM or ProxmoxLxc provider driver will
// use these settings.
type defaultplaneproxmox struct {
	// The fully qualified name of the image to use for all planes.
	Image string `yaml:"image"`
	PublicNetwork string `yaml:"publicnetwork"`
	// Name of storage device to use for storing the VM/container disk image. This should be the same on all nodes, or
	// (preferrably) a network share on CephFS or NFS should be used.
	DiskStorage string `yaml:"diskstorage"'`
	// The default node to run containers/VMs on. In the future this will be replaced with a better load-balancing
	// solution.
	DefaultNode string `yaml:"defaultnode"`
}

// Loads the plane defaults from a file. In the future, this may also respond to environment variables an inject them
// into the final struct.
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
