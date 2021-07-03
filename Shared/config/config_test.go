package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_InitConfiguration(t *testing.T) {
	InitConfiguration("shared")
}

func Test_GetGlobalString(t *testing.T) {
	InitConfiguration("shared")
	require.NoError(t, os.Setenv("CAPTAIN_TEST", "TEST"))
	assert.Equal(t, "TEST", GetGlobalString("TEST"))
}

func Test_GetAppString(t *testing.T) {
	InitConfiguration("shared")
	require.NoError(t, os.Setenv("CAPTAIN_SHARED_TEST", "TEST"))
	assert.Equal(t, "TEST", GetAppString("TEST"))
}

func Test_ConvertsAppNameToUppercase(t *testing.T) {
	InitConfiguration("shared")
	require.NoError(t, os.Setenv("CAPTAIN_SHARED_TEST", "TEST"))
	assert.Equal(t, "TEST", GetAppString("TEST"))
}

func Test_SetAppString(t *testing.T) {
	InitConfiguration("shared")
	SetAppString("TEST", "testvalue")
	assert.Equal(t, "testvalue", GetAppString("TEST"))
}

func Test_GetAppInt(t *testing.T) {
	InitConfiguration("shared")
	require.NoError(t, os.Setenv("CAPTAIN_SHARED_TEST", "1"))
	assert.Equal(t, 1, GetAppInt("TEST"))
}
