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

func formatArrayType(typeName string) string {
	return fmt.Sprintf("[]%s", typeName)
}

func formatEllipsis(typeName string) string {
	return fmt.Sprintf("... %s", typeName)
}

func formatMapType(keyType string, valueType string) string {
	return fmt.Sprintf("map[%s]%s", keyType, valueType)
}

func formatFunctionField(field domain.Field) string {
	return fmt.Sprintf("%s %s", field.Name, field.Type)
}

func formatFuncType(function domain.Function) string {
	var parameters []string
	var returns []string
	for _, parameterField := range function.Parameters {
		parameters = append(parameters, formatFunctionField(parameterField))
	}

	for _, returnField := range function.ReturnFields {
		returns = append(parameters, formatFunctionField(returnField))
	}
	return fmt.Sprintf("func(%s) %s", strings.Join(parameters, ", "), strings.Join(returns, ", "))
}

func formatChanType(valueType string) string {
	return fmt.Sprintf("chan %s", valueType)
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
	var typeName string
	eltField, err := exprToField("", fieldType.Elt)

	if err == nil && eltField != nil {
		typeName = eltField.Type.ToString()
	}
	return domain.Field{
		Name: fieldName,
		Type: domain.Type(formatArrayType(typeName)),
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

func mapTypeToField(fieldName string, fieldType *ast.MapType) domain.Field {
	var err error
	var keyType string
	var valueType string
	keyField, err := exprToField("", fieldType.Key)
	if err == nil && keyField != nil {
		keyType = keyField.Type.ToString()
	}

	valueField, err := exprToField("", fieldType.Value)
	if err == nil && valueField != nil {
		valueType = valueField.Type.ToString()
	}
	return domain.Field{
		Name: fieldName,
		Type: domain.Type(formatMapType(keyType, valueType)),
	}
}

func funcTypeToField(fieldName string, fieldType *ast.FuncType) domain.Field {
	function := domain.Function{}
	if fieldType.Params != nil {
		function.Parameters = ParseFields(fieldType.Params.List)
	}
	if fieldType.Results != nil {
		function.ReturnFields = ParseFields(fieldType.Results.List)
	}
	return domain.Field{
		Name: fieldName,
		Type: domain.Type(formatFuncType(function)),
	}
}

func structTypeToField(fieldName string, fieldType *ast.StructType) domain.Field {
	return domain.Field{
		Name: fieldName,
		Type: "interface{}",
	}
}

func chanTypeToField(fieldName string, fieldType *ast.ChanType) domain.Field {
	var valueFieldType string
	valueField, err := exprToField("", fieldType.Value)

	if err == nil && valueField != nil {
		valueFieldType = valueField.Type.ToString()
	}
	return domain.Field{
		Name: fieldName,
		Type: domain.Type(formatChanType(valueFieldType)),
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
