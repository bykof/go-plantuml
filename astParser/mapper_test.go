package astParser

import (
	"go/ast"
	"testing"

	"github.com/magiconair/properties/assert"

	"github.com/bykof/go-plantuml/domain"
)

func Test_starExprToField(t *testing.T) {
	assert.Equal(
		t,
		starExprToField("test", &ast.StarExpr{
			X: &ast.Ident{
				Name: "string",
			},
		}),
		domain.Field{
			Name:     "test",
			Type:     "*string",
			Nullable: true,
		},
	)
}

func Test_identToField(t *testing.T) {
	assert.Equal(
		t,
		identToField("test", &ast.Ident{Name: "string"}),
		domain.Field{
			Name: "test",
			Type: "string",
		},
	)
}

func Test_selectorExprToField(t *testing.T) {
	assert.Equal(
		t,
		selectorExprToField(
			"test",
			&ast.SelectorExpr{
				Sel: &ast.Ident{Name: "testSelect"},
				X:   &ast.Ident{Name: "testPackage"},
			},
		),
		domain.Field{
			Name: "test",
			Type: "testSelect",
		},
	)
}

func Test_arrayTypeToField(t *testing.T) {
	assert.Equal(
		t,
		arrayTypeToField(
			"test",
			&ast.ArrayType{
				Elt: &ast.Ident{Name: "string"},
			},
		),
		domain.Field{
			Name: "test",
			Type: "[]string",
		},
	)
}
func Test_ellipsisToField(t *testing.T) {
	assert.Equal(
		t,
		ellipsisToField(
			"test",
			&ast.Ellipsis{
				Elt: &ast.Ident{Name: "string"},
			},
		),
		domain.Field{
			Name: "test",
			Type: "... string",
		},
	)
}
func Test_interfaceTypeToField(t *testing.T) {
	assert.Equal(
		t,
		interfaceTypeToField(
			"test",
			&ast.InterfaceType{},
		),
		domain.Field{
			Name: "test",
			Type: "interface",
		},
	)
}
func Test_mapTypeToField(t *testing.T) {
	assert.Equal(
		t,
		mapTypeToField(
			"test",
			&ast.MapType{
				Key: &ast.Ident{
					Name: "string",
				},
				Value: &ast.Ident{
					Name: "int64",
				},
			},
		),
		domain.Field{
			Name: "test",
			Type: "map[string]int64",
		},
	)
}
func Test_funcTypeToField(t *testing.T) {
	params := []*ast.Field{
		{
			Names: []*ast.Ident{{Name: "bla"}},
			Type:  &ast.Ident{Name: "string"},
		},
		{
			Names: []*ast.Ident{{Name: "anotherBla"}},
			Type:  &ast.Ident{Name: "int"},
		},
	}

	returns := []*ast.Field{
		{
			Names: []*ast.Ident{{Name: "bla"}},
			Type:  &ast.Ident{Name: "string"},
		},
	}
	assert.Equal(
		t,
		funcTypeToField(
			"test",
			&ast.FuncType{
				Params:  &ast.FieldList{List: params},
				Results: &ast.FieldList{List: returns},
			},
		),
		domain.Field{
			Name: "test",
			Type: "func(bla string, anotherBla int) string",
		},
	)
}
func Test_structTypeToField(t *testing.T) {
	assert.Equal(
		t,
		structTypeToField(
			"test",
			&ast.StructType{},
		),
		domain.Field{
			Name: "test",
			Type: "interface{}",
		},
	)
}
func Test_chanTypeToField(t *testing.T) {
	assert.Equal(
		t,
		chanTypeToField(
			"test",
			&ast.ChanType{
				Value: &ast.Ident{
					Name: "string",
				},
			},
		),
		domain.Field{
			Name: "test",
			Type: "chan string",
		},
	)
}
