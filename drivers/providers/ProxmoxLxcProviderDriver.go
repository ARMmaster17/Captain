package providers

import (
	"crypto/tls"
	"fmt"
	"github.com/ARMmaster17/Captain/ImageStore"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

// ProxmoxLxcProviderDriver is an implementation of ProviderDriver that interfaces with the LXC capabilities of a
// Proxmox cluster.
type ProxmoxLxcProviderDriver struct {
	client *proxmox.Client
}

// Connect reads any Proxmox connection parameters defined in the environment and performs authentication with the
// designated Proxmox cluster.
func (d ProxmoxLxcProviderDriver) Connect() error {
	tlsConf := &tls.Config{InsecureSkipVerify: !viper.GetBool(d.getConfigItemPath("forcessl"))}
	d.client, _ = proxmox.NewClient(viper.GetString(d.getConfigItemPath("url")), nil, tlsConf, 300)
	return d.client.Login(os.Getenv("CAPTAIN_PROXMOX_USER"), os.Getenv("CAPTAIN_PROXMOX_PASSWORD"), "")
}

// BuildPlane converts the given GenericPlane into a data structure that is usable by the Proxmox API and submits
// the container for provisioning.
func (d ProxmoxLxcProviderDriver) BuildPlane(p *GenericPlane) (string, error) {
	config := proxmox.NewConfigLxc()
	template, err := ImageStore.GetProviderSpecificImageConfiguration(d.GetYAMLTag(), viper.GetString("defaults.image"))
	if err != nil {
		return "", err
	}
	config.Ostemplate = template
	config.Arch = "amd64"
	config.CMode = "tty"
	config.Console = true
	config.Cores = p.Cores
	config.CPULimit = 0
	config.CPUUnits = 1024
	config.Description = "Managed by the Captain stack"
	config.Hostname = p.FQDN
	config.Memory = p.RAM
	config.Nameserver = viper.GetString("defaults.network.nameservers")
	config.Networks = proxmox.QemuDevices{
		0: {
			"name":     "eth0",
			"bridge":   viper.GetString(d.getConfigItemPath("publicnetwork")),
			"ip":       "dhcp",
			"gw":       viper.GetString("defaults.network.gateway"),
			"firewall": "0",
			"mtu":      viper.GetInt("defaults.network.mtu"),
		},
	}
	config.OnBoot = true
	config.Protection = false
	config.SearchDomain = viper.GetString("defaults.network.searchdomain")
	config.SSHPublicKeys = viper.GetString("defaults.publickey")
	config.Start = true
	config.Storage = viper.GetString(d.getConfigItemPath("diskstorage"))
	config.Swap = p.RAM
	config.Template = false
	config.Tty = 2
	config.Unprivileged = true

	nextID, err := d.client.GetNextID(0)
	if err != nil {
		return "", fmt.Errorf("unable to retreive next available VMID with error: %w", err)
	}
	vmr := proxmox.NewVmRef(nextID)
	vmr.SetNode(viper.GetString(d.getConfigItemPath("defaultnode")))
	err = config.CreateLxc(vmr, d.client)
	if err != nil {
		return "", fmt.Errorf("unable to create LXC container with error: %w", err)
	}
	return fmt.Sprintf("%s:%d", d.GetCUIDPrefix(), vmr.VmId()), nil
}

// DestroyPlane will destroy a plane that is managed by the Proxmox LXC driver.
func (d ProxmoxLxcProviderDriver) DestroyPlane(cuid string, p *GenericPlane) error {
	vmr, err := d.client.GetVmRefByName(p.FQDN)
	if err != nil {
		return fmt.Errorf("unable to obtain reference to underlying LXC container for plane %s: %w", p.FQDN, err)
	}
	_, err = d.client.StopVm(vmr)
	if err != nil {
		return fmt.Errorf("unable to stop LXC container: %w", err)
	}
	err = d.proxmoxOverrideDeleteVMParams(vmr)
	if err != nil {
		return fmt.Errorf("unable to delete LXC container for plane %s: %w", p.FQDN, err)
	}
	return nil
}

// GetCUIDPrefix gets the prefix that should be added to the beginning of the CUID strings for all planes that are
// managed by this driver.
func (d ProxmoxLxcProviderDriver) GetCUIDPrefix() string {
	return "proxmox.lxc"
}

// GetYAMLTag gets the YAML tag that the Proxmox LXC driver uses to identify settings unique to this driver
// in config.yaml.
func (d ProxmoxLxcProviderDriver) GetYAMLTag() string {
	return "proxmoxlxc"
}

// getFullYAMLTag gets the base YAML tag for all settings that are unique to the Proxmox LXC driver.
func (d ProxmoxLxcProviderDriver) getFullYAMLTag() string {
	return fmt.Sprintf("drivers.provisioners.%s", d.GetYAMLTag())
}

// getConfigItemPath returns the full YAML path to a configuration file for the Proxmox LXC driver with the given tag(s) appended
// at the end.
func (d ProxmoxLxcProviderDriver) getConfigItemPath(entryPath string) string {
	return fmt.Sprintf("%s.%s", d.getFullYAMLTag(), entryPath)
}

// proxmoxOverrideDeleteVMParams replaces the method with the same name in the proxmox library because there is a bug where if you pass
// an empty struct to any DELETE endpoint, the Proxmox API returns an error. This method overrides that by passing
// nil to the underlying Session object.
func (d *ProxmoxLxcProviderDriver) proxmoxOverrideDeleteVMParams(vmr *proxmox.VmRef) error {
	err := d.client.CheckVmRef(vmr)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("/nodes/%s/%s/%d", vmr.Node(), vmr.GetVmType(), vmr.VmId())
	return d.customDELETERequest(url)
}

// customDELETERequest performs a DELETE request. This is used as a workaround for broken methods in the
// telmate/proxmox-api-go library.
func (d *ProxmoxLxcProviderDriver) customDELETERequest(url string) error {
	var taskResponse map[string]interface{}
	session, err := d.buildCustomSession()
	if err != nil {
		return err
	}
	resp, err := session.RequestJSON("DELETE", url, nil, nil, nil, &taskResponse)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send DELETE request to the Proxmox API: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the Proxmxox API returned status code %d", resp.StatusCode)
	}
	return nil
}

// buildCustomSession builds a Proxmox API session object without using the telmate/proxmox-api-go library.
func (d *ProxmoxLxcProviderDriver) buildCustomSession() (*proxmox.Session, error) {
	session, err := proxmox.NewSession(viper.GetString(d.getConfigItemPath("url")), nil, &tls.Config{
		InsecureSkipVerify: !viper.GetBool(d.getConfigItemPath("forcessl")),
	})
	if err != nil {
		return &proxmox.Session{}, fmt.Errorf("unable to connect to the Proxmox API: %w", err)
	}
	err = session.Login(os.Getenv("CAPTAIN_PROXMOX_USER"), os.Getenv("CAPTAIN_PROXMOX_PASSWORD"), "")
	if err != nil {
		return &proxmox.Session{}, fmt.Errorf("unable to authenticate with the Proxmox API: %w", err)
	}
	return session, nil
}
