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
			// Function is not bound to a struct
			if functionDecl.Recv == nil {
				continue
			}

			className := ""
			if ident, ok := functionDecl.Recv.List[0].Type.(*ast.Ident); ok {
				className = ident.Name
			}

			if starExpr, ok := functionDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
				className = starExpr.X.(*ast.Ident).Name
			}

			classIndex := classes.ClassIndexByName(className)
			if classIndex != -1 {
				function := domain.Function{}
				function.Name = functionDecl.Name.Name
				if functionDecl.Type == nil {
					print(functionDecl)
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

func formatArrayType(ident *ast.Ident) string {
	return fmt.Sprintf("[]%s", ident.Name)
}

func formatEllipsis(typeName string) string {
	return fmt.Sprintf("... %s", typeName)
}

func startExprToField(name string, starExpr *ast.StarExpr) domain.Field {
	if expr, ok := starExpr.X.(*ast.SelectorExpr); ok {
		return domain.Field{
			Name:     name,
			Type:     domain.Type(expr.Sel.Name),
			Nullable: true,
		}
	}

	if ident, ok := starExpr.X.(*ast.Ident); ok {
		return domain.Field{
			Name:     name,
			Type:     domain.Type(ident.Name),
			Nullable: true,
		}
	}
	return domain.Field{}
}

func identToField(fieldName string, fieldType *ast.Ident) domain.Field {
	return domain.Field{
		Name: fieldName,
		Type: domain.Type(fieldType.Name),
	}
}

func selectorExprToField(fieldName string, fieldType *ast.SelectorExpr) domain.Field {
	var packageName string
	selectorField, err := exprToField("", fieldType.X)

	if err == nil && selectorField != nil {
		packageName = selectorField.Type.ToString()
	}

	return domain.Field{
		Name:    fieldName,
		Type:    domain.Type(fieldType.Sel.Name),
		Package: domain.Package(packageName),
	}
}

func arrayTypeToField(fieldName string, fieldType *ast.ArrayType) domain.Field {
	return domain.Field{
		Name: fieldName,
		Type: domain.Type(formatArrayType(fieldType.Elt.(*ast.Ident))),
	}
}

func ellipsisToField(fieldName string, fieldType *ast.Ellipsis) domain.Field {
	var typeName string
	eltField, err := exprToField("", fieldType.Elt)

	if err == nil && eltField != nil {
		typeName = eltField.Type.ToString()
	}
	return domain.Field{
		Name: fieldName,
		Type: domain.Type(formatEllipsis(typeName)),
	}
}

func interfaceTypeToField(fieldName string, fieldType *ast.InterfaceType) domain.Field {
	return domain.Field{
		Name: fieldName,
		Type: domain.Type("interface"),
	}
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
		field := startExprToField(fieldName, fieldType)
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
	default:
		return nil, fmt.Errorf("unknown Field Type %s", reflect.TypeOf(expr).String())
	}
}

func getField(field *ast.Field) (*domain.Field, error) {
	var fieldName string

	if field.Names != nil && len(field.Names) > 0 {
		fieldName = field.Names[0].Name
	}
	return exprToField(fieldName, field.Type)

}

func ParseFields(fieldList []*ast.Field) domain.Fields {
	fields := domain.Fields{}
	for _, field := range fieldList {
		parsedField, err := getField(field)
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
