package astParser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/bykof/go-plantuml/domain"
)

func ParseDirectory(directoryPath string, recursive bool) domain.Packages {
	var packages domain.Packages
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		log.Fatal(err)
	}
	currentPackage := domain.Package{
		FilePath:   directoryPath,
		Name:       directoryPath,
		Variables:  domain.Fields{},
		Constants:  domain.Fields{},
		Interfaces: domain.Interfaces{},
		Classes:    domain.Classes{},
		Functions:  domain.Functions{},
	}
	for _, file := range files {
		fullPath := filepath.Join(directoryPath, file.Name())
		if !file.IsDir() {
			if filepath.Ext(file.Name()) != ".go" || strings.Contains(file.Name(), "_test") {
				continue
			}
			parsedPackage := ParseFile(fullPath)
			currentPackage = currentPackage.Add(parsedPackage)
		} else {
			if recursive {
				packages = append(packages, ParseDirectory(fullPath, recursive)...)
			}
		}
	}

	if !currentPackage.IsEmpty() {
		packages = append(packages, currentPackage)
	}

	return packages
}

func ParseFile(filePath string) domain.Package {
	var domainPackage domain.Package

	node, err := parser.ParseFile(
		token.NewFileSet(),
		filePath,
		nil,
		parser.ParseComments,
	)
	if err != nil {
		log.Fatal(err)
	}

	domainPackage = domain.Package{
		FilePath:   filePath,
		Name:       filePath,
		Interfaces: domain.Interfaces{},
		Classes:    domain.Classes{},
		Functions:  domain.Functions{},
		Constants:  domain.Fields{},
		Variables:  domain.Fields{},
	}

	if node.Scope != nil {
		for name, object := range node.Scope.Objects {
			// If object is not a type
			switch object.Kind {
			case ast.Var:
				field, err := valueSpecToField(object.Name, object.Decl.(*ast.ValueSpec))
				if err != nil {
					log.Fatal(err)
				}
				field.Name = fmt.Sprintf("var %s", field.Name)
				domainPackage.Variables = append(domainPackage.Variables, *field)
			case ast.Con:
				field, err := valueSpecToField(object.Name, object.Decl.(*ast.ValueSpec))
				if err != nil {
					log.Fatal(err)
				}
				field.Name = fmt.Sprintf("const %s", field.Name)
				domainPackage.Constants = append(domainPackage.Constants, *field)
			case ast.Typ:
				typeSpec := object.Decl.(*ast.TypeSpec)

				switch typeSpec.Type.(type) {
				case *ast.StructType:
					structType := typeSpec.Type.(*ast.StructType)
					class := domain.Class{
						Name:   name,
						Fields: ParseFields(structType.Fields.List),
					}

					domainPackage.Classes = append(domainPackage.Classes, class)
				case *ast.InterfaceType:
					var functions domain.Functions
					interfaceType := typeSpec.Type.(*ast.InterfaceType)

					for _, field := range interfaceType.Methods.List {
						if funcType, ok := field.Type.(*ast.FuncType); ok {
							parsedFields, err := ParseField(field)
							if err != nil {
								log.Fatal(err)
							}
							for _, parsedField := range parsedFields {
								functions = append(functions, funcTypeToFunction(parsedField.Name, funcType))
							}

						}
					}

					domainInterface := domain.Interface{
						Name:      name,
						Functions: functions,
					}

					domainPackage.Interfaces = append(domainPackage.Interfaces, domainInterface)
				}
			}
		}
	}

	for _, decl := range node.Decls {
		if functionDecl, ok := decl.(*ast.FuncDecl); ok {
			var className string

			// Function is not bound to a struct
			if functionDecl.Recv == nil {
				function := createFunction(functionDecl.Name.Name, functionDecl)
				domainPackage.Functions = append(domainPackage.Functions, function)
				continue
			}

			for _, receiverClass := range functionDecl.Recv.List {
				classField, err := exprToField("", receiverClass.Type)
				if err != nil {
					log.Fatal(err)
				}

				className = classField.Type.ToClassString()
				classIndex := domainPackage.Classes.ClassIndexByName(className)

				// Handle the case where className could not be found in classes
				if classIndex < 0 {
					log.Printf("Could not find class: %s for function %s", className, functionDecl.Name.Name)
					continue
				}

				function := createFunction(functionDecl.Name.Name, functionDecl)
				domainPackage.Classes[classIndex].Functions = append(
					domainPackage.Classes[classIndex].Functions,
					function,
				)
			}
		}
	}
	return domainPackage
}

func createFunction(name string, functionDecl *ast.FuncDecl) domain.Function {
	function := domain.Function{
		Name: name,
	}
	if functionDecl.Type.Params != nil {
		function.Parameters = ParseFields(functionDecl.Type.Params.List)
	}
	if functionDecl.Type.Results != nil {
		function.ReturnFields = ParseFields(functionDecl.Type.Results.List)
	}
	return function
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

func ParseField(field *ast.Field) (domain.Fields, error) {
	var fields domain.Fields

	if field.Names != nil && len(field.Names) > 0 {
		for _, fieldName := range field.Names {
			parsedField, err := exprToField(fieldName.Name, field.Type)
			if err != nil {
				return fields, err
			}
			fields = append(fields, *parsedField)
		}
	} else {
		parsedField, err := exprToField("", field.Type)
		if err != nil {
			return fields, err
		}
		fields = append(fields, *parsedField)
	}
	return fields, nil

}

func ParseFields(fieldList []*ast.Field) domain.Fields {
	fields := domain.Fields{}
	for _, field := range fieldList {
		parsedFields, err := ParseField(field)
		if err != nil {
			log.Fatal(err)
		}

		fields = append(fields, parsedFields...)
	}
	return fields
}
