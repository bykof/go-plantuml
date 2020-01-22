package formatter

import (
	"fmt"
	"github.com/bykof/go-plantuml/domain"
	"strings"
)

const PlantUMLClassFormat = `class %s {
%s
%s
}`

const PlantUMLWrapper = `@startuml
%s
@enduml`

func FormatField(field domain.Field) string {
	visibilityCharacter := "+"
	if field.IsPrivate() {
		visibilityCharacter = "-"
	}
	return fmt.Sprintf("%s%s: %s", visibilityCharacter, field.Name, field.Type)
}

func FormatFields(fields domain.Fields) string {
	var formattedFields []string
	for _, field := range fields {
		formattedFields = append(formattedFields, FormatField(field))
	}
	return strings.Join(formattedFields, "\n")
}

func FormatParameter(parameter domain.Field) string {
	return fmt.Sprintf("%s %s", parameter.Name, parameter.Type)
}

func FormatParameters(parameters domain.Fields) string {
	var formattedParameters []string
	for _, parameter := range parameters {
		formattedParameters = append(formattedParameters, FormatParameter(parameter))
	}
	return strings.Join(formattedParameters, ", ")
}

func FormatReturnField(returnField domain.Field) string {
	return string(returnField.Type)
}

func FormatReturnFields(returnFields domain.Fields) string {
	var formattedReturnFields []string
	for _, returnField := range returnFields {
		formattedReturnFields = append(formattedReturnFields, FormatReturnField(returnField))
	}
	return strings.Join(formattedReturnFields, ", ")
}

func FormatPlantUMLWrapper(content ...string) string {
	return fmt.Sprintf(PlantUMLWrapper, strings.Join(content, "\n"))
}

func FormatFunction(function domain.Function) string {
	visibilityCharacter := "+"
	if function.IsPrivate() {
		visibilityCharacter = "-"
	}
	formattedReturnFields := FormatReturnFields(function.ReturnFields)
	formattedParameters := FormatParameters(function.Parameters)
	return fmt.Sprintf(
		"%s%s(%s): %s",
		visibilityCharacter,
		function.Name,
		formattedParameters,
		formattedReturnFields,
	)
}

func FormatFunctions(functions domain.Functions) string {
	var formattedFunctions []string
	for _, function := range functions {
		formattedFunctions = append(formattedFunctions, FormatFunction(function))
	}
	return strings.Join(formattedFunctions, "\n")
}

func FormatClass(class domain.Class) string {
	formattedFields := FormatFields(class.Fields)
	formattedFunctions := FormatFunctions(class.Functions)
	return fmt.Sprintf(PlantUMLClassFormat, class.Name, formattedFields, formattedFunctions)
}

func FormatClasses(classes domain.Classes) string {
	var formattedClasses []string
	for _, class := range classes {
		formattedClasses = append(formattedClasses, FormatClass(class))
	}

	return strings.Join(formattedClasses, "\n")
}

func FormatPlantUML(classes domain.Classes) string {
	formattedClasses := FormatClasses(classes)
	formattedRelations := FormatRelations(classes)
	return FormatPlantUMLWrapper(formattedClasses, formattedRelations)
}

func FormatRelation(class domain.Class, class2 domain.Class) string {
	return fmt.Sprintf("%s -- %s", class.Name, class2.Name)
}

func FormatRelations(classes domain.Classes) string {
	var formattedRelations []string
	if len(classes) < 2 {
		return ""
	}
	for i := 0; i <= len(classes)-2; i++ {
		for j := i + 1; j <= len(classes)-1; j++ {
			if classes[i].HasRelation(classes[j]) {
				formattedRelations = append(formattedRelations, FormatRelation(classes[i], classes[j]))
			}
			if classes[j].HasRelation(classes[i]) {
				formattedRelations = append(formattedRelations, FormatRelation(classes[j], classes[i]))
			}
		}
	}
	return strings.Join(formattedRelations, "\n")
}
