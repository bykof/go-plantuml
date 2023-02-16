package astParser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/bykof/go-plantuml/domain"
)

var earlySeenFunctions map[string]domain.Functions

func ParseDirectory(directoryPath string, opts ...ParserOptionFunc) domain.Packages {
	options := &parserOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var packages domain.Packages
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		log.Fatal(err)
	}

	earlySeenFunctions = make(map[string]domain.Functions)
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
			if isExcluded(file.Name(), options.excludedFilesRegex) {
				continue
			}

			parsedPackage := ParseFile(fullPath)
			currentPackage = currentPackage.Add(parsedPackage)
		} else {
			if options.recursive {
				packages = append(packages, ParseDirectory(fullPath, opts...)...)
			}
		}
	}

	if !currentPackage.IsEmpty() {
		packages = append(packages, currentPackage)
	}

	// Handle the case where className could not be found in classes
	for className, earlyFuncs := range earlySeenFunctions {
		for _, function := range earlyFuncs {
			log.Printf("Could not find class: %s for function %s", className, function.Name)
		}
	}

	return packages
}

func isExcluded(fileName string, regex *regexp.Regexp) bool {
	if filepath.Ext(fileName) != ".go" {
		return true
	}

	if strings.HasSuffix(fileName, "_test.go") {
		return true
	}

	if regex == nil {
		return false
	}

	if regex.Match([]byte(fileName)) {
		return true
	}

	return false
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

				if typeSpec.TypeParams != nil {
					var params []string
					for _, v := range typeSpec.TypeParams.List {
						for _, v := range v.Names {
							params = append(params, v.String())
						}
					}
					if len(params) > 0 {
						name = name + "[" + strings.Join(params, ",") + "]"
					}
				}
				switch typeSpec.Type.(type) {
				case *ast.StructType:
					structType := typeSpec.Type.(*ast.StructType)
					class := domain.Class{
						Name:   name,
						Fields: ParseFields(structType.Fields.List),
					}

					addEarlyFunctions(&class)
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
				default:
					class := domain.Class{
						Name: name,
					}

					addEarlyFunctions(&class)
					domainPackage.Classes = append(domainPackage.Classes, class)
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

				function := createFunction(functionDecl.Name.Name, functionDecl) // Handle the case where className could not be found in classes

				if classIndex < 0 {
					earlySeenFunctions[className] = append(earlySeenFunctions[className], function)
				} else {
					domainPackage.Classes[classIndex].Functions = append(
						domainPackage.Classes[classIndex].Functions,
						function,
					)
				}
			}
		}
	}
	return domainPackage
}

func addEarlyFunctions(class *domain.Class) {
	if funcs, ok := earlySeenFunctions[class.Name]; ok {
		class.Functions = append(class.Functions, funcs...)
		delete(earlySeenFunctions, class.Name)
	}
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
	case *ast.IndexExpr: // generic goes here
		field, err := indexExprToField(fieldName, fieldType)
		return &field, err
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
