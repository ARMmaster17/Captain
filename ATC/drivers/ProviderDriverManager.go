package drivers

import (
	"fmt"
	"github.com/ARMmaster17/Captain/ATC/drivers/providers"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
)

// BuildPlaneOnAnyProvider chooses a suitable driver based on the available drivers defined in config.yaml, and
// triggers a build on that platform through the selected driver. Returns the fully-qualified CUID identifier for that
// plane.
func BuildPlaneOnAnyProvider(p *providers.GenericPlane) (string, error) {
	driver, err := getActiveBuildDriver()
	if err != nil {
		return "", fmt.Errorf("unable to get a provisioning driver:\n%w", err)
	}
	err = driver.Connect()
	if err != nil {
		return "", fmt.Errorf("unable to initialize %s driver:\n%w", driver.GetYAMLTag(), err)
	}

	cuid, err := driver.BuildPlane(p)
	if err != nil {
		return "", fmt.Errorf("unable to build plane with driver %s:\n%w", driver.GetYAMLTag(), err)
	}
	return fmt.Sprintf("%s:%s", driver.GetCUIDPrefix(), cuid), nil
}

// DestroyPlane will destroy a plane and clean up any leftover configuration. A driver is chosen based on the plane's
// CUID string, which is of the format driver:UID.
func DestroyPlane(p *providers.GenericPlane) error {
	driver, err := getDestroyDriver(p.CUID)
	if err != nil {
		return fmt.Errorf("unable to get a destruction driver:\n%w", err)
	}
	err = driver.Connect()
	if err != nil {
		return fmt.Errorf("unable to initialize %s driver:\n%w", driver.GetYAMLTag(), err)
	}
	cuidsplit := strings.Split(p.CUID, ":")[1]
	err = driver.DestroyPlane(cuidsplit, p)
	if err != nil {
		return fmt.Errorf("unable to destroy plane with driver %s:\n%w", driver.GetYAMLTag(), err)
	}
	return nil
}

// getActiveBuildDriver looks at the current configuration and picks an appropriate driver to use for provisioning
// a new plane. Drivers are tried in-order as they are listed in config.yaml.
func getActiveBuildDriver() (providers.ProviderDriver, error) {
	activedrivers := viper.GetStringSlice("config.drivers.provisioners")
	if len(activedrivers) == 0 {
		return nil, fmt.Errorf("expected at least 1 driver in config.yaml, found 0")
	}
	for _, driverentry := range activedrivers {
		driver, err := driverLookupByYAMLTag(driverentry)
		if err != nil {
			log.Warn().Msg(fmt.Sprintf("unknown driver %s in config.yaml... skipping", driverentry))
		} else {
			return driver, nil
		}
	}
	return nil, fmt.Errorf("no usable drivers found in config.yaml")
}

// getDestroyDriver returns a driver instance from a fully-qualified CUID string from a plane instance.
func getDestroyDriver(cuid string) (providers.ProviderDriver, error) {
	driverentry := strings.Split(cuid, ":")[0]
	return driverLookupByCUIDPrefix(driverentry)
}

// driverLookupByYAMLTag returns a driver instance by comparing the given string to every registered driver's YAML tag.
func driverLookupByYAMLTag(driverentry string) (providers.ProviderDriver, error) {
	if driverentry == (&providers.ProxmoxLxcProviderDriver{}).GetYAMLTag() {
		// In the future there could be a check here to see if a driver is ready to receive new planes.
		return &providers.ProxmoxLxcProviderDriver{}, nil
	}
	if driverentry == (&providers.DummyProviderDriver{}).GetYAMLTag() {
		return providers.DummyProviderDriver{}, nil
	}
	return nil, fmt.Errorf("%s is not a valid driver type", driverentry)
}

// driverLookupByCUIDPrefix returns a driver instance by comparing the given string to every registered driver's YAML tag.
func driverLookupByCUIDPrefix(cuidprefix string) (providers.ProviderDriver, error) {
	if cuidprefix == (&providers.ProxmoxLxcProviderDriver{}).GetCUIDPrefix() {
		// In the future there could be a check here to see if a driver is ready to receive new planes.
		return &providers.ProxmoxLxcProviderDriver{}, nil
	}
	if cuidprefix == (&providers.DummyProviderDriver{}).GetCUIDPrefix() {
		return providers.DummyProviderDriver{}, nil
	}
	return nil, fmt.Errorf("%s is not a valid CUID prefix", cuidprefix)
}
