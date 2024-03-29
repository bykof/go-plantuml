package astParser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bykof/go-plantuml/domain"
)

func Test_formatArrayType(t *testing.T) {
	assert.Equal(t, "[]test", formatArrayType("test"))
	assert.Equal(t, "[]", formatArrayType(""))
}

func Test_formatEllipsis(t *testing.T) {
	assert.Equal(t, "... test", formatEllipsis("test"))
	assert.Equal(t, "... ", formatEllipsis(""))
}

func Test_formatMapType(t *testing.T) {
	assert.Equal(t, "map[test]anotherTest", formatMapType("test", "anotherTest"))
	assert.Equal(t, "map[]", formatMapType("", ""))
}

func Test_formatFunctionField(t *testing.T) {
	assert.Equal(t, "test anotherTest", formatFunctionParameterField(domain.Field{Name: "test", Type: "anotherTest"}))
	assert.Equal(t, " ", formatFunctionParameterField(domain.Field{Name: "", Type: ""}))
}

func Test_formatFuncType(t *testing.T) {
	assert.Equal(
		t,
		"func(arg1 string, arg2 []bla) string, error",
		formatFuncType(domain.Function{
			Parameters: domain.Fields{
				domain.Field{Name: "arg1", Type: "string"},
				domain.Field{Name: "arg2", Type: "[]bla"},
			},
			ReturnFields: domain.Fields{
				domain.Field{Type: "string"},
				domain.Field{Type: "error"},
			},
		}),
	)
	assert.Equal(
		t,
		"func() string, error",
		formatFuncType(domain.Function{
			ReturnFields: domain.Fields{
				domain.Field{Type: "string"},
				domain.Field{Type: "error"},
			},
		}),
	)
	assert.Equal(
		t,
		"func() ",
		formatFuncType(domain.Function{}),
	)
}

func Test_formatChanType(t *testing.T) {
	assert.Equal(t, "chan test", formatChanType("test"))
	assert.Equal(t, "chan ", formatChanType(""))
}

func Test_formatPointerType(t *testing.T) {
	assert.Equal(t, "*test", formatPointer("test"))
	assert.Equal(t, "*", formatPointer(""))
}
