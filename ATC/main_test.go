package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGenerateConfigFile(t *testing.T) {
	_ = os.Remove("/etc/captain/atc/config.yaml")
	err := generateConfigFile()
	assert.NoError(t, err)
	info, err := os.Stat("/etc/captain/atc/config.yaml")
	assert.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}
