package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunction_IsPrivate(t *testing.T) {
	functionIsPublic := Function{Name: "Test"}
	functionIsPrivate := Function{Name: "test"}
	functionNameEmpty := Function{Name: ""}
	assert.True(t, functionIsPrivate.IsPrivate())
	assert.False(t, functionIsPublic.IsPrivate())
	assert.False(t, functionNameEmpty.IsPrivate())
}

func TestFunction_EqualImplementation(t *testing.T) {
	assert.False(t, Function{
		Name: "Run",
		Parameters: Fields{
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
		},
	}.EqualImplementation(
		Function{
			Name: "RunA",
			Parameters: Fields{
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
			},
		},
	))
	assert.True(t, Function{
		Name: "Run",
		Parameters: Fields{
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
		},
		ReturnFields: Fields{
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
		},
	}.EqualImplementation(
		Function{
			Name: "Run",
			Parameters: Fields{
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
			},
			ReturnFields: Fields{
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
			},
		},
	))
}
