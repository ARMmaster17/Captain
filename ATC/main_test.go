package main

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGenerateConfigFile(t *testing.T) {
	_ = os.Remove("/etc/captain/atc/config.yaml")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/captain/atc")
	err := generateConfigFile()
	assert.NoError(t, err)
	info, err := os.Stat("/etc/captain/atc/config.yaml")
	assert.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}
