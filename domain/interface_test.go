package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterface_IsImplementedBy(t *testing.T) {
	assert.True(t, Interface{
		Name: "NiceInterface",
		Functions: Functions{
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
			},
			Function{
				Name: "OtherRun",
				Parameters: Fields{
					Field{
						Name:     "Same",
						Nullable: false,
						Type:     "Same",
					},
				},
				ReturnFields: Fields{
					Field{
						Name:     "Return",
						Nullable: false,
						Type:     "Same",
					},
				},
			},
		},
	}.IsImplementedByClass(
		Class{
			Name: "NiceInterface",
			Functions: Functions{
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
				},
				Function{
					Name: "OtherRun",
					Parameters: Fields{
						Field{
							Name:     "Same",
							Nullable: false,
							Type:     "Same",
						},
					},
					ReturnFields: Fields{
						Field{
							Name:     "Return",
							Nullable: false,
							Type:     "Same",
						},
					},
				},
			},
		},
	))

	assert.False(t, Interface{
		Name: "NiceInterface",
		Functions: Functions{
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
			},
			Function{
				Name: "OtherRun",
				Parameters: Fields{
					Field{
						Name:     "Same",
						Nullable: false,
						Type:     "Same",
					},
				},
				ReturnFields: Fields{
					Field{
						Name:     "Return",
						Nullable: false,
						Type:     "Same",
					},
				},
			},
		},
	}.IsImplementedByClass(
		Class{
			Name: "NiceInterface",
			Functions: Functions{
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
				},
				Function{
					Name: "OtherRun",
					Parameters: Fields{
						Field{
							Name:     "Same",
							Nullable: false,
							Type:     "Different",
						},
					},
					ReturnFields: Fields{
						Field{
							Name:     "Return",
							Nullable: false,
							Type:     "Same",
						},
					},
				},
			},
		},
	))
}
