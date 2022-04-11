package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestType_ToString(t *testing.T) {
	assert.Equal(t, "Test", Type("Test").ToString())
}
