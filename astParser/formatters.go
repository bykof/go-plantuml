package astParser

import (
	"fmt"
	"github.com/bykof/go-plantuml/domain"
	"strings"
)

func formatArrayType(typeName string) string {
	return fmt.Sprintf("[]%s", typeName)
}

func formatEllipsis(typeName string) string {
	return fmt.Sprintf("... %s", typeName)
}

func formatMapType(keyType string, valueType string) string {
	return fmt.Sprintf("map[%s]%s", keyType, valueType)
}

func formatFunctionParameterField(field domain.Field) string {
	return fmt.Sprintf("%s %s", field.Name, field.Type)
}

func formatFunctionReturnField(field domain.Field) string {
	return field.Type.ToString()
}

func formatPointer(name string) string {
	return fmt.Sprintf("*%s", name)
}

func formatFuncType(function domain.Function) string {
	var parameters []string
	var returns []string
	for _, parameterField := range function.Parameters {
		parameters = append(parameters, formatFunctionParameterField(parameterField))
	}

	for _, returnField := range function.ReturnFields {
		returns = append(returns, formatFunctionReturnField(returnField))
	}
	return fmt.Sprintf("func(%s) %s", strings.Join(parameters, ", "), strings.Join(returns, ", "))
}

func formatChanType(valueType string) string {
	return fmt.Sprintf("chan %s", valueType)
}
