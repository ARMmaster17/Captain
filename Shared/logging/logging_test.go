package logging

import (
	"github.com/ARMmaster17/Captain/Shared/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_InitLogging(t *testing.T) {
	require.NoError(t, config.InitConfiguration("shared"))
	config.SetAppString("LOG_PATH", "/etc/captain/shared/test.log")
	assert.NoError(t, InitLogging("shared"))
}
