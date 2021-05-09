package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestE2EBootstrapDryRun(t *testing.T) {
	helperDeleteDBIfExists()
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	assert.NoError(t, BootstrapCluster())
}
