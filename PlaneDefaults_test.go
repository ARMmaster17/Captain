package main

import "testing"

func TestLoadPlaneDefaults(t *testing.T) {
	defaults, err := getPlaneDefaults()
	if err != nil {
		t.Errorf("unexpected error getting plane defaults: %w", err)
		return
	}
	if helperAssertStringsEqual(t, "publickey", defaults.PublicKey, "x") {
		return
	}
}

func TestLoadPlaneDefaultsNetwork(t *testing.T) {
	defaults, err := getPlaneDefaults()
	if err != nil {
		t.Errorf("unexpected error getting plane defaults: %w", err)
		return
	}
	if helperAssertStringsEqual(t, "network.searchdomain", defaults.Network.SearchDomain, "") {
		return
	}
	if helperAssertStringsEqual(t, "network.nameservers", defaults.Network.Nameservers, "8.8.8.8 8.8.4.4") {
		return
	}
	if helperAssertStringsEqual(t, "network.gateway", defaults.Network.Gateway, "10.1.0.1") {
		return
	}
	if helperAssertStringsEqual(t, "network.mtu", defaults.Network.MTU, "1450") {
		return
	}
}

func TestLoadPlaneDefaultsProxmox(t *testing.T) {
	defaults, err := getPlaneDefaults()
	if err != nil {
		t.Errorf("unexpected error getting plane defaults: %w", err)
		return
	}
	if helperAssertStringsEqual(t, "proxmox.image", defaults.Proxmox.Image, "pve-img:vztmpl/debian-10-standard_10.7-1_amd64.tar.gz") {
		return
	}
	if helperAssertStringsEqual(t, "proxmox.publicnetwork", defaults.Proxmox.PublicNetwork, "internal") {
		return
	}
	if helperAssertStringsEqual(t, "proxmox.diskstorage", defaults.Proxmox.DiskStorage, "pve-storage") {
		return
	}
}

// Returns true if two strings are not equal and automatically fails the test.
func helperAssertStringsEqual(t *testing.T, name string, expect string, actual string) bool {
	if expect != actual {
		t.Errorf("%s: expected '%s', got '%s'", name, expect, actual)
		return true
	}
	return false
}
