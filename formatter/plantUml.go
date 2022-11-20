package formatter

import (
	"fmt"
	"strings"

	"github.com/bykof/go-plantuml/domain"
)

const PlantUMLInterfaceFormat = `interface %s{
%s
}`

const PlantUMLPackageFormat = `package %s{
%s
}`

const PlantUMLAnnotationFormat = `annotation %s {
%s
%s
}`

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

func FormatInterface(domainInterface domain.Interface) string {
	return fmt.Sprintf(
		PlantUMLInterfaceFormat,
		domainInterface.Name,
		FormatFunctions(domainInterface.Functions),
	)
}

func FormatInterfaces(interfaces domain.Interfaces) string {
	var formattedInteraces []string
	for _, domainInterface := range interfaces {
		formattedInteraces = append(formattedInteraces, FormatInterface(domainInterface))
	}
	return strings.Join(formattedInteraces, "\n")
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

func FormatPackages(domainPackages domain.Packages) string {
	var formattedPackages []string
	for _, domainPackage := range domainPackages {
		formattedPackages = append(formattedPackages, FormatPackage(domainPackage))
	}
	return strings.Join(formattedPackages, "\n")
}

func FormatPackageAnnotations(domainPackage domain.Package) string {
	var formattedPackageFunctions string

	if len(domainPackage.Functions) == 0 && len(domainPackage.Constants) == 0 && len(domainPackage.Variables) == 0 {
		return formattedPackageFunctions
	}

	if len(domainPackage.Functions) > 0 {
		formattedPackageFunctions = FormatFunctions(domainPackage.Functions)
	}

	var formattedFields []string
	formattedFields = append(formattedFields, FormatFields(domainPackage.Constants))
	formattedFields = append(formattedFields, FormatFields(domainPackage.Variables))

	return fmt.Sprintf(
		PlantUMLAnnotationFormat,
		domainPackage.Name,
		strings.Join(formattedFields, "\n"),
		formattedPackageFunctions,
	)
}

func FormatPackage(domainPackage domain.Package) string {
	var packageContent []string
	formattedPackageAnnotations := FormatPackageAnnotations(domainPackage)
	if formattedPackageAnnotations != "" {
		packageContent = append(packageContent, formattedPackageAnnotations)
	}

	if len(domainPackage.Interfaces) > 0 {
		packageContent = append(packageContent, FormatInterfaces(domainPackage.Interfaces))

	}

	if len(domainPackage.Classes) > 0 {
		packageContent = append(packageContent, FormatClasses(domainPackage.Classes))
	}

	return fmt.Sprintf(
		PlantUMLPackageFormat,
		domainPackage.Name,
		strings.Join(packageContent, "\n"),
	)
}

func FormatFunction(function domain.Function) string {
	var formattedReturnFields string
	visibilityCharacter := "+"
	if function.IsPrivate() {
		visibilityCharacter = "-"
	}
	formattedParameters := FormatParameters(function.Parameters)

	if len(function.ReturnFields) > 0 {
		formattedReturnFields = fmt.Sprintf(": %s", FormatReturnFields(function.ReturnFields))
	}

	return fmt.Sprintf(
		"%s%s(%s)%s",
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
	formattedClass := fmt.Sprintf(PlantUMLClassFormat, class.Name, formattedFields, formattedFunctions)
	return formattedClass
}

func FormatClasses(classes domain.Classes) string {
	var formattedClasses []string
	for _, class := range classes {
		formattedClasses = append(formattedClasses, FormatClass(class))
	}

	return strings.Join(formattedClasses, "\n")
}

func FormatPlantUML(packages domain.Packages) string {
	formattedPackages := FormatPackages(packages)
	formattedRelations := FormatRelations(packages.AllClasses())
	formattedImplementationRelations := FormatImplementationRelations(
		packages.AllClasses(),
		packages.AllInterfaces(),
	)
	return FormatPlantUMLWrapper(formattedPackages, formattedRelations, formattedImplementationRelations)
}

func FormatRelation(class domain.Class, class2 domain.Class) string {
	return fmt.Sprintf(`"%s" --> "%s"`, class.Name, class2.Name)
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

func FormatImplementationRelation(class domain.Class, domainInterface domain.Interface) string {
	return fmt.Sprintf(`"%s" --|> "%s"`, class.Name, domainInterface.Name)
}

func FormatImplementationRelations(classes domain.Classes, domainInterfaces domain.Interfaces) string {
	var formattedImplementationRelations []string
	for _, class := range classes {
		for _, domainInterface := range domainInterfaces {
			if domainInterface.IsImplementedByClass(class) {
				formattedImplementationRelations = append(
					formattedImplementationRelations,
					FormatImplementationRelation(class, domainInterface),
				)
			}
		}
	}
	return strings.Join(formattedImplementationRelations, "\n")
}
