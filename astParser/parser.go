package astParser

import (
	"errors"
	"fmt"
	"github.com/bykof/go-plantuml/domain"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"strings"
)

func ParseDirectory(directoryPath string, recursive bool) domain.Classes {
	var classes domain.Classes
	files, err := ioutil.ReadDir(directoryPath)
	for _, file := range files {
		fullPath := filepath.Join(directoryPath, file.Name())
		if !file.IsDir() {
			if filepath.Ext(file.Name()) != ".go" || strings.Contains(file.Name(), "_test") {
				continue
			}
			classes = append(classes, ParseFile(fullPath)...)
		} else {
			if recursive {
				classes = append(classes, ParseDirectory(fullPath, recursive)...)
			}
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	return classes
}

func ParseFile(filePath string) domain.Classes {
	var classes domain.Classes
	node, err := parser.ParseFile(
		token.NewFileSet(),
		filePath,
		nil,
		parser.ParseComments,
	)

	if err != nil {
		log.Fatal(err)
	}

	if node.Scope != nil {
		for name, object := range node.Scope.Objects {
			// If object is not a type
			if object.Kind != ast.Typ {
				continue
			}
			typeSpec := object.Decl.(*ast.TypeSpec)
			structType, ok := typeSpec.Type.(*ast.StructType)
			// We are probably dealing with interface types
			if !ok {
				continue
			}
			class := domain.Class{Name: name, Package: domain.Package(node.Name.Name)}
			class.Fields = ParseFields(structType.Fields.List)
			classes = append(classes, class)
		}
	}

	for _, decl := range node.Decls {
		if functionDecl, ok := decl.(*ast.FuncDecl); ok {
			var className string

			// Function is not bound to a struct
			if functionDecl.Recv == nil {
				continue
			}

			classField, err := exprToField("", functionDecl.Recv.List[0].Type)
			if err != nil {
				log.Fatal(err)
			}

			className = classField.Type.ToString()

			if len(classes) == 0 {
				return classes
			}

			classIndex := classes.ClassIndexByName(className)

			if classIndex != -1 {
				function := domain.Function{
					Name: functionDecl.Name.Name,
				}
				if functionDecl.Type.Params != nil {
					function.Parameters = ParseFields(functionDecl.Type.Params.List)
				}
				if functionDecl.Type.Results != nil {
					function.ReturnFields = ParseFields(functionDecl.Type.Results.List)
				}
				classes[classIndex].Functions = append(classes[classIndex].Functions, function)
			}
		}
	}
	return classes
}

func exprToField(fieldName string, expr ast.Expr) (*domain.Field, error) {
	switch fieldType := expr.(type) {
	case *ast.Ident:
		field := identToField(fieldName, fieldType)
		return &field, nil
	case *ast.SelectorExpr:
		field := selectorExprToField(fieldName, fieldType)
		return &field, nil
	case *ast.StarExpr:
		field := starExprToField(fieldName, fieldType)
		return &field, nil
	case *ast.ArrayType:
		field := arrayTypeToField(fieldName, fieldType)
		return &field, nil
	case *ast.Ellipsis:
		field := ellipsisToField(fieldName, fieldType)
		return &field, nil
	case *ast.InterfaceType:
		field := interfaceTypeToField(fieldName, fieldType)
		return &field, nil
	case *ast.MapType:
		field := mapTypeToField(fieldName, fieldType)
		return &field, nil
	case *ast.FuncType:
		field := funcTypeToField(fieldName, fieldType)
		return &field, nil
	case *ast.StructType:
		field := structTypeToField(fieldName, fieldType)
		return &field, nil
	case *ast.ChanType:
		field := chanTypeToField(fieldName, fieldType)
		return &field, nil
	default:
		return nil, fmt.Errorf("unknown Field Type %s", reflect.TypeOf(expr).String())
	}
}

func ParseField(field *ast.Field) (*domain.Field, error) {
	var fieldName string

	if field.Names != nil && len(field.Names) > 0 {
		fieldName = field.Names[0].Name
	}
	return exprToField(fieldName, field.Type)

}

func ParseFields(fieldList []*ast.Field) domain.Fields {
	fields := domain.Fields{}
	for _, field := range fieldList {
		parsedField, err := ParseField(field)
		if err != nil {
			log.Fatal(err)
		}

		if parsedField == nil {
			log.Fatal(errors.New("unexpected error: parsedField is nil"))
		}
		fields = append(fields, *parsedField)
	}
	return fields
}
