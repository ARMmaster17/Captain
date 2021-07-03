package framework

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_NewFramework(t *testing.T) {
	framework, err := NewFramework("")
	assert.NoError(t, err)
	assert.NotNil(t, framework)
}

func Test_NewFrameworkWithName(t *testing.T) {
	framework, err := NewFramework("shared")
	assert.NoError(t, err)
	require.NotEqual(t, Framework{}, framework)
	assert.Equal(t, "shared", framework.AppName)
}

func Test_FrameworkCreatesRouterOnInit(t *testing.T) {
	framework, err := NewFramework("shared")
	assert.NoError(t, err)
	require.NotEqual(t, Framework{}, framework)
	assert.NotNil(t, framework.Router)
}

// End goal ----------------------------------
//func Test_E2ECustomEndpoint(t *testing.T) {
//	framework := NewFramework()
//	framework.AddEndpoint()
//	framework.StartAsync()
//	// Connect through HTTP/S
//	// Verify with assert(t, ...)
//	framework.Stop()
//}
