package metadata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetCaptainVersion(t *testing.T) {
	result := GetCaptainVersion()
	assert.NotNil(t, result)
}