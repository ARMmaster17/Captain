package providers

type DummyProviderDriver struct {

}

func (d DummyProviderDriver) Connect() error {
	return nil
}

func (d DummyProviderDriver) BuildPlane(p *GenericPlane) (string, error) {
	return p.FQDN, nil
}

func (d DummyProviderDriver) DestroyPlane(cuid string, p *GenericPlane) error {
	return nil
}

func (d DummyProviderDriver) GetCUIDPrefix() string {
	return "dummy"
}

func (d DummyProviderDriver) GetYAMLTag() string {
	return "dummy"
}
