package providers

// DummyProviderDriver is a provider driver used for testing both in unit tests, and for trying out Captain without
// connecting the stack to a live hypervisor cluster.
type DummyProviderDriver struct {
}

// Connect simulates connecting to a provider driver. Always returns no error.
func (d DummyProviderDriver) Connect() error {
	return nil
}

// BuildPlane simulates building a plane. The GenericPlane's FQDN is returned as the CUID and no action is performed.
func (d DummyProviderDriver) BuildPlane(p *GenericPlane) (string, error) {
	return p.FQDN, nil
}

// DestroyPlane simulates destroying a plane. No action is performed as the dummy driver does not use real resources.
func (d DummyProviderDriver) DestroyPlane(cuid string, p *GenericPlane) error {
	return nil
}

// GetCUIDPrefix returns the prefix used in the CUID of planes that are managed by this driver.
func (d DummyProviderDriver) GetCUIDPrefix() string {
	return "dummy"
}

// GetYAMLTag returns the tag used for configuration files that are specific to the dummy driver in config.yaml.
func (d DummyProviderDriver) GetYAMLTag() string {
	return "dummy"
}
