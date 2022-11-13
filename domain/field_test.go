package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestField_IsPrivate(t *testing.T) {
	fieldIsPublic := Field{Name: "Test"}
	fieldIsPrivate := Field{Name: "test"}
	fieldNameEmpty := Field{Name: ""}
	assert.True(t, fieldIsPrivate.IsPrivate())
	assert.False(t, fieldIsPublic.IsPrivate())
	assert.False(t, fieldNameEmpty.IsPrivate())
}

func TestField_EqualImplementation(t *testing.T) {
	field1 := Field{
		Name:     "Same",
		Nullable: false,
		Type:     "Same",
	}
	field2 := Field{
		Name:     "Same",
		Nullable: false,
		Type:     "Same",
	}
	field3 := Field{
		Name:     "Nullable",
		Nullable: true,
		Type:     "Same",
	}
	field4 := Field{
		Name:     "Type",
		Nullable: false,
		Type:     "Different",
	}
	field5 := Field{
		Name:     "Different",
		Nullable: true,
		Type:     "Different",
	}

	assert.True(t, field1.EqualImplementation(field2))
	assert.False(t, field1.EqualImplementation(field3))
	assert.False(t, field1.EqualImplementation(field4))
	assert.False(t, field1.EqualImplementation(field5))
}

func TestFields_EqualImplementations(t *testing.T) {
	assert.False(t, Fields{Field{}, Field{}}.EqualImplementations(Fields{Field{}}))
	assert.True(t, Fields{
		Field{
			Name:     "Same",
			Nullable: false,
			Type:     "Same",
		},
		Field{
			Name:     "Same2",
			Nullable: false,
			Type:     "Same2",
		},
	}.EqualImplementations(Fields{
		Field{
			Name:     "Same",
			Nullable: false,
			Type:     "Same",
		},
		Field{
			Name:     "Same2",
			Nullable: false,
			Type:     "Same2",
		},
	}))

	assert.False(t, Fields{
		Field{
			Name:     "Same",
			Nullable: false,
			Type:     "Same",
		},
		Field{
			Name:     "Same2",
			Nullable: false,
			Type:     "Same2",
		},
	}.EqualImplementations(Fields{
		Field{
			Name:     "Same2",
			Nullable: false,
			Type:     "Same2",
		},
		Field{
			Name:     "Same",
			Nullable: false,
			Type:     "Same",
		},
	}))

	assert.False(t, Fields{
		Field{
			Name:     "Same",
			Nullable: false,
			Type:     "Same",
		},
		Field{
			Name:     "Same2",
			Nullable: false,
			Type:     "Same2",
		},
	}.EqualImplementations(Fields{
		Field{
			Name:     "Same",
			Nullable: false,
			Type:     "Same",
		},
	}))

}
