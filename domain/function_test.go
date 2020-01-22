package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunction_IsPrivate(t *testing.T) {
	functionIsPublic := Function{Name: "Test"}
	functionIsPrivate := Function{Name: "test"}
	functionNameEmpty := Function{Name: ""}
	assert.True(t, functionIsPrivate.IsPrivate())
	assert.False(t, functionIsPublic.IsPrivate())
	assert.False(t, functionNameEmpty.IsPrivate())
}
