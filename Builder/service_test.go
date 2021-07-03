package Builder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewBuilder(t *testing.T) {
	builder, err := NewBuilder()
	assert.NoError(t, err)
	assert.NotNil(t, builder)
}
