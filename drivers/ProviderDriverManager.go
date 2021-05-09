package drivers

import (
	"fmt"
	"github.com/ARMmaster17/Captain/drivers/providers"
)

func BuildPlaneOnAnyProvider(p *providers.GenericPlane) (string, error) {
	// We're going to assume that the Proxmox LXC driver is
	// initialized since it's the only driver available at the
	// moment.
	driver := providers.ProxmoxLxcProviderDriver{}
	err := driver.Connect()
	if err != nil {
		return "", fmt.Errorf("unable to initialize %s driver: %w", driver.GetYAMLTag(), err)
	}
	cuid, err := driver.BuildPlane(p)
	if err != nil {
		return "", fmt.Errorf("unable to build plane with driver %s: %w", driver.GetYAMLTag(), err)
	}
	return fmt.Sprintf("%s:%s", driver.GetCUIDPrefix(), cuid), nil
}

func DestroyPlane(p *providers.GenericPlane) error {
	// Again, we can assume that any given GenericPlane exists
	// on a hypervisor managed by the Proxmox LXC driver.
	driver := providers.ProxmoxLxcProviderDriver{}
	err := driver.Connect()
	if err != nil {
		return fmt.Errorf("unable to initialize %s driver: %w", driver.GetYAMLTag(), err)
	}
	err = driver.DestroyPlane(p)
	if err != nil {
		return fmt.Errorf("unable to destroy plane with driver %s: %w", driver.GetYAMLTag(), err)
	}
	return nil
}
