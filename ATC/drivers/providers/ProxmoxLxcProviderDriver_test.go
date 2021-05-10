package providers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProxmoxLxcInit(t *testing.T) {
	d := ProxmoxLxcProviderDriver{}
	assert.NotNil(t, d)
}

func TestProxmoxLxcGetCUIDPrefix(t *testing.T) {
	d := ProxmoxLxcProviderDriver{}
	assert.Equal(t, "proxmox.lxc", d.GetCUIDPrefix())
}

func TestProxmoxLxcGetYAMLTag(t *testing.T) {
	d := ProxmoxLxcProviderDriver{}
	assert.Equal(t, "proxmoxlxc", d.GetYAMLTag())
}
