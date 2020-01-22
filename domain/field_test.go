package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestField_IsPrivate(t *testing.T) {
	fieldIsPublic := Field{Name: "Test"}
	fieldIsPrivate := Field{Name: "test"}
	fieldNameEmpty := Field{Name: ""}
	assert.True(t, fieldIsPrivate.IsPrivate())
	assert.False(t, fieldIsPublic.IsPrivate())
	assert.False(t, fieldNameEmpty.IsPrivate())
}
