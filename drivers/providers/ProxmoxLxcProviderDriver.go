package providers

import (
	"crypto/tls"
	"fmt"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

type ProxmoxLxcProviderDriver struct {
	client *proxmox.Client
}

func (d *ProxmoxLxcProviderDriver) Connect() error {
	tlsConf := &tls.Config{InsecureSkipVerify: true}
	d.client, _ = proxmox.NewClient(os.Getenv("CAPTAIN_PROXMOX_URL"), nil, tlsConf, 300)
	return d.client.Login(os.Getenv("CAPTAIN_PROXMOX_USER"), os.Getenv("CAPTAIN_PROXMOX_PASSWORD"), "")
}

func (d *ProxmoxLxcProviderDriver) BuildPlane(p *GenericPlane) (string, error) {
	config := proxmox.NewConfigLxc()
	config.Ostemplate = viper.GetString(d.getConfigItemPath("image"))
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
		0 : {
			"name": "eth0",
			"bridge": viper.GetString(d.getConfigItemPath("publicnetwork")),
			"ip": "dhcp",
			"gw": viper.GetString("defaults.network.gateway"),
			"firewall": "0",
			"mtu": viper.GetInt("defaults.network.mtu"),
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

func (d *ProxmoxLxcProviderDriver) DestroyPlane(p *GenericPlane) error {
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

func (d *ProxmoxLxcProviderDriver) GetCUIDPrefix() string {
	return "proxmox.lxc"
}

func (d *ProxmoxLxcProviderDriver) GetYAMLTag() string {
	return "proxmoxlxc"
}

func (d *ProxmoxLxcProviderDriver) getFullYAMLTag() string {
	return fmt.Sprintf("drivers.provisioners.%s", d.GetYAMLTag())
}

func (d *ProxmoxLxcProviderDriver) getConfigItemPath(entryPath string) string {
	return fmt.Sprintf("%s.%s", d.getFullYAMLTag(), entryPath)
}

// This method replaces the method with the same name in the proxmox library because there is a bug where if you pass
// an empty struct to any DELETE endpoint, the Proxmox API returns an error. This method overrides that by passing
// nil to the underlying Session object.
func (d *ProxmoxLxcProviderDriver) proxmoxOverrideDeleteVMParams(vmr *proxmox.VmRef) error {
	err := d.client.CheckVmRef(vmr)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("/nodes/%s/%s/%d", vmr.Node(), vmr.GetVmType(), vmr.VmId())
	var taskResponse map[string]interface{}
	session, err := proxmox.NewSession(os.Getenv("CAPTAIN_PROXMOX_URL"), nil, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return fmt.Errorf("unable to connect to the Proxmox API: %w", err)
	}
	err = session.Login(os.Getenv("CAPTAIN_PROXMOX_USER"), os.Getenv("CAPTAIN_PROXMOX_PASSWORD"), "")
	if err != nil {
		return fmt.Errorf("unable to authenticate with the Proxmox API: %w", err)
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