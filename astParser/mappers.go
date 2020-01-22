package astParser

import (
	"github.com/bykof/go-plantuml/domain"
	"go/ast"
)

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
