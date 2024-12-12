package formatter

import (
	"github.com/bykof/go-plantuml/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatField__non_nullabe(t *testing.T) {
	field := domain.Field{Name: "FieldName", Type: "string", Nullable: false}
	expected := "+FieldName: string"
	result := FormatField(field)
	assert.Equal(t, expected, result)
}

func TestFormatField__nullabe(t *testing.T) {
	field := domain.Field{Name: "FieldName", Type: "string", Nullable: true}
	expected := "+FieldName: string"
	result := FormatField(field)
	assert.Equal(t, expected, result)
}

func TestFormatField__pointer(t *testing.T) {
	field := domain.Field{Name: "FieldName", Type: "*string"}
	expected := "+FieldName: *string"
	result := FormatField(field)
	assert.Equal(t, expected, result)
}

func TestFormatField__private(t *testing.T) {
	field := domain.Field{Name: "fieldName", Type: "string"}
	expected := "-fieldName: string"
	result := FormatField(field)
	assert.Equal(t, expected, result)
}

func TestFormatFields(t *testing.T) {
	fields := domain.Fields{
		{Name: "Field1", Type: "int"},
		{Name: "Field2", Type: "string"},
	}
	expected := "+Field1: int\n+Field2: string"
	result := FormatFields(fields)
	assert.Equal(t, expected, result)
}

func TestFormatInterface(t *testing.T) {
	domainInterface := domain.Interface{
		Name:      "MyInterface",
		Functions: domain.Functions{},
	}
	expected := "interface MyInterface{\n\n}"
	result := FormatInterface(domainInterface)
	assert.Equal(t, expected, result)
}

func TestFormatInterfaces(t *testing.T) {
	interfaces := domain.Interfaces{
		{Name: "Interface1"},
		{Name: "Interface2"},
	}
	expected := "interface Interface1{\n\n}\ninterface Interface2{\n\n}"
	result := FormatInterfaces(interfaces)
	assert.Equal(t, expected, result)
}

func TestFormatParameter(t *testing.T) {
	parameter := domain.Field{Name: "param", Type: "int"}
	expected := "param int"
	result := FormatParameter(parameter)
	assert.Equal(t, expected, result)
}

func TestFormatParameters(t *testing.T) {
	parameters := domain.Fields{
		{Name: "param1", Type: "int"},
		{Name: "param2", Type: "string"},
	}
	expected := "param1 int, param2 string"
	result := FormatParameters(parameters)
	assert.Equal(t, expected, result)
}

func TestFormatReturnField(t *testing.T) {
	returnField := domain.Field{Type: "string"}
	expected := "string"
	result := FormatReturnField(returnField)
	assert.Equal(t, expected, result)
}

func TestFormatReturnFields(t *testing.T) {
	returnFields := domain.Fields{
		{Type: "string"},
		{Type: "int"},
	}
	expected := "string, int"
	result := FormatReturnFields(returnFields)
	assert.Equal(t, expected, result)
}

func TestFormatPlantUMLWrapper(t *testing.T) {
	content := []string{"Content1", "Content2"}
	expected := "@startuml\nContent1\nContent2\n@enduml"
	result := FormatPlantUMLWrapper(content...)
	assert.Equal(t, expected, result)
}

func TestFormatPackages(t *testing.T) {
	packages := domain.Packages{
		{Name: "Package1"},
		{Name: "Package2"},
	}
	expected := "package Package1{\n\n}\npackage Package2{\n\n}"
	result := FormatPackages(packages)
	assert.Equal(t, expected, result)
}

func TestFormatPackageAnnotations(t *testing.T) {
	pkg := domain.Package{Name: "TestPackage", Constants: domain.Fields{
		domain.Field{Name: "Field1", Type: "int"},
		domain.Field{Name: "Field2", Type: "string"},
	}}
	expected := "annotation TestPackage {\n+Field1: int\n+Field2: string\n\n\n}"
	result := FormatPackageAnnotations(pkg)
	assert.Equal(t, expected, result)
}

func TestFormatPackage(t *testing.T) {
	pkg := domain.Package{Name: "TestPackage", Constants: domain.Fields{
		domain.Field{Name: "Field1", Type: "int"},
		domain.Field{Name: "Field2", Type: "string"},
	}}
	expected := "package TestPackage{\nannotation TestPackage {\n+Field1: int\n+Field2: string\n\n\n}\n}"
	result := FormatPackage(pkg)
	assert.Equal(t, expected, result)
}

func TestFormatFunction(t *testing.T) {
	function := domain.Function{
		Name:       "TestFunction",
		Parameters: domain.Fields{},
	}
	expected := "+TestFunction()"
	result := FormatFunction(function)
	assert.Equal(t, expected, result)
}

func TestFormatFunctions(t *testing.T) {
	functions := domain.Functions{
		{Name: "Function1"},
		{Name: "Function2"},
	}
	expected := "+Function1()\n+Function2()"
	result := FormatFunctions(functions)
	assert.Equal(t, expected, result)
}

func TestFormatClass(t *testing.T) {
	class := domain.Class{
		Name: "TestClass",
	}
	expected := "class TestClass {\n\n\n}"
	result := FormatClass(class)
	assert.Equal(t, expected, result)
}

func TestFormatClasses(t *testing.T) {
	classes := domain.Classes{
		{Name: "Class1"},
		{Name: "Class2"},
	}
	expected := "class Class1 {\n\n\n}\nclass Class2 {\n\n\n}"
	result := FormatClasses(classes)
	assert.Equal(t, expected, result)
}

func TestFormatPlantUML(t *testing.T) {
	packages := domain.Packages{}
	expected := "@startuml\n\n\n\n@enduml"
	result := FormatPlantUML(packages)
	assert.Equal(t, expected, result)
}

func TestFormatRelation(t *testing.T) {
	class1 := domain.Class{Name: "Class1"}
	class2 := domain.Class{Name: "Class2"}
	expected := `"Class1" --> "Class2"`
	result := FormatRelation(class1, class2)
	assert.Equal(t, expected, result)
}

func TestFormatRelations(t *testing.T) {
	classes := domain.Classes{
		{Name: "Class1"},
		{Name: "Class2"},
	}
	expected := ""
	result := FormatRelations(classes)
	assert.Equal(t, expected, result)
}

func TestFormatImplementationRelation(t *testing.T) {
	class := domain.Class{Name: "Class1"}
	iface := domain.Interface{Name: "Interface1"}
	expected := `"Class1" --|> "Interface1"`
	result := FormatImplementationRelation(class, iface)
	assert.Equal(t, expected, result)
}

func TestFormatImplementationRelations(t *testing.T) {
	classes := domain.Classes{}
	interfaces := domain.Interfaces{}
	expected := ""
	result := FormatImplementationRelations(classes, interfaces)
	assert.Equal(t, expected, result)
}
