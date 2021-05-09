package providers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDummyDriverInit(t *testing.T) {
	d := DummyProviderDriver{}
	assert.NotNil(t, d)
}

func TestDummyDriverConnect(t *testing.T) {
	d := DummyProviderDriver{}
	assert.NoError(t, d.Connect())
}

func TestDummyDriverBuildPlane(t *testing.T) {
	d := DummyProviderDriver{}
	p := GenericPlane{
		FQDN:  "plane1.example.com",
		CUID:  "",
		Cores: 1,
		RAM:   512,
		Disk:  8,
	}
	cuid, err := d.BuildPlane(&p)
	assert.NoError(t, err)
	assert.Equal(t, "plane1.example.com", cuid)
}

func TestDummyDriverDestroyPlane(t *testing.T) {
	d := DummyProviderDriver{}
	p := GenericPlane{
		FQDN:  "plane1.example.com",
		CUID:  "",
		Cores: 1,
		RAM:   512,
		Disk:  8,
	}
	assert.NoError(t, d.DestroyPlane("", &p))
}

func TestDummyDriverGetCUIDPrefix(t *testing.T) {
	d := DummyProviderDriver{}
	assert.Equal(t, "dummy", d.GetCUIDPrefix())
}

func TestDummyDriverGetYAMLTag(t *testing.T) {
	d := DummyProviderDriver{}
	assert.Equal(t, "dummy", d.GetYAMLTag())
}
