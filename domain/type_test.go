package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType_ToString(t *testing.T) {
	assert.Equal(t, "Test", Type("Test").ToString())
}
