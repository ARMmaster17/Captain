package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_InitConfiguration(t *testing.T) {
	assert.NoError(t, InitConfiguration("SHARED"))
}

func Test_GetGlobalString(t *testing.T) {
	require.NoError(t, InitConfiguration("SHARED"))
	require.NoError(t, os.Setenv("CAPTAIN_TEST", "TEST"))
	assert.Equal(t, "TEST", GetGlobalString("TEST"))
}

func Test_GetAppString(t *testing.T) {
	require.NoError(t, InitConfiguration("SHARED"))
	require.NoError(t, os.Setenv("CAPTAIN_SHARED_TEST", "TEST"))
	assert.Equal(t, "TEST", GetAppString("TEST"))
}

func Test_ConvertsAppNameToUppercase(t *testing.T) {
	require.NoError(t, InitConfiguration("shared"))
	require.NoError(t, os.Setenv("CAPTAIN_SHARED_TEST", "TEST"))
	assert.Equal(t, "TEST", GetAppString("TEST"))
}

func Test_SetAppString(t *testing.T) {
	require.NoError(t, InitConfiguration("shared"))
	SetAppString("TEST", "testvalue")
	assert.Equal(t, "testvalue", GetAppString("TEST"))
}
