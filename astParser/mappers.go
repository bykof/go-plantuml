package astParser

import (
	"go/ast"

	"github.com/bykof/go-plantuml/domain"
)

func starExprToField(fieldName string, starExpr *ast.StarExpr) domain.Field {
	var fieldType string
	xField, err := exprToField("", starExpr.X)
	if err == nil && xField != nil {
		fieldType = xField.Type.ToString()
	}
	return domain.Field{
		Name:     fieldName,
		Type:     domain.Type(formatPointer(fieldType)),
		Nullable: true,
	}
}

func identToField(fieldName string, ident *ast.Ident) domain.Field {
	return domain.Field{
		Name: fieldName,
		Type: domain.Type(ident.Name),
	}
}

func selectorExprToField(fieldName string, selectorExpr *ast.SelectorExpr) domain.Field {
	return domain.Field{
		Name: fieldName,
		Type: domain.Type(selectorExpr.Sel.Name),
	}
}

func arrayTypeToField(fieldName string, arrayType *ast.ArrayType) domain.Field {
	var typeName string
	eltField, err := exprToField("", arrayType.Elt)

	if err == nil && eltField != nil {
		typeName = eltField.Type.ToString()
	}
	return domain.Field{
		Name: fieldName,
		Type: domain.Type(formatArrayType(typeName)),
	}
}

func ellipsisToField(fieldName string, ellipsis *ast.Ellipsis) domain.Field {
	var typeName string
	eltField, err := exprToField("", ellipsis.Elt)

	if err == nil && eltField != nil {
		typeName = eltField.Type.ToString()
	}
	return domain.Field{
		Name: fieldName,
		Type: domain.Type(formatEllipsis(typeName)),
	}
}

func interfaceTypeToField(fieldName string, interfaceType *ast.InterfaceType) domain.Field {
	return domain.Field{
		Name: fieldName,
		Type: domain.Type("interface"),
	}
}

func mapTypeToField(fieldName string, mapType *ast.MapType) domain.Field {
	var err error
	var keyType string
	var valueType string
	keyField, err := exprToField("", mapType.Key)
	if err == nil && keyField != nil {
		keyType = keyField.Type.ToString()
	}

	valueField, err := exprToField("", mapType.Value)
	if err == nil && valueField != nil {
		valueType = valueField.Type.ToString()
	}
	return domain.Field{
		Name: fieldName,
		Type: domain.Type(formatMapType(keyType, valueType)),
	}
}

func funcTypeToField(fieldName string, funcType *ast.FuncType) domain.Field {
	function := funcTypeToFunction(fieldName, funcType)
	return domain.Field{
		Name: function.Name,
		Type: domain.Type(formatFuncType(function)),
	}
}

func structTypeToField(fieldName string, structType *ast.StructType) domain.Field {
	return domain.Field{
		Name: fieldName,
		Type: "interface{}",
	}
}

func chanTypeToField(fieldName string, chanType *ast.ChanType) domain.Field {
	var valueFieldType string
	valueField, err := exprToField("", chanType.Value)

	if err == nil && valueField != nil {
		valueFieldType = valueField.Type.ToString()
	}
	return domain.Field{
		Name: fieldName,
		Type: domain.Type(formatChanType(valueFieldType)),
	}
}

func valueSpecToField(fieldName string, valueSpec *ast.ValueSpec) (*domain.Field, error) {
	var fieldType string

	if valueSpec.Type != nil {
		return exprToField(fieldName, valueSpec.Type)
	}

	if valueSpec.Values != nil && len(valueSpec.Values) > 0 {
		switch valueSpec.Values[0].(type) {
		case *ast.BasicLit:
			fieldType = valueSpec.Values[0].(*ast.BasicLit).Kind.String()
		}
	}

	return &domain.Field{
		Name: fieldName,
		Type: domain.Type(fieldType),
	}, nil
}

func funcTypeToFunction(fieldName string, funcType *ast.FuncType) domain.Function {
	function := domain.Function{
		Name: fieldName,
	}
	if funcType.Params != nil {
		function.Parameters = ParseFields(funcType.Params.List)
	}
	if funcType.Results != nil {
		function.ReturnFields = ParseFields(funcType.Results.List)
	}
	return function
}
