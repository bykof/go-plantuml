package domain

type (
	Class struct {
		Name      string
		Package   Package
		Fields    Fields
		Functions Functions
	}

	Classes []Class
)

func (class Class) HasRelation(toClass Class) bool {
	for _, field := range class.Fields {
		if string(field.Type) == toClass.Name {
			return true
		}
	}
	return false
}

func (classes *Classes) ClassByName(name string) *Class {
	if classes == nil {
		return nil
	}

	for _, class := range *classes {
		if class.Name == name {
			return &class
		}
	}
	return nil
}

func (classes Classes) ClassIndexByName(name string) int {
	for index, class := range classes {
		if class.Name == name {
			return index
		}
	}
	return -1
}

func (classes Classes) AllPackages() Packages {
	var packages Packages
	checked := make(map[Package]bool)
	for _, class := range classes {
		if _, found := checked[class.Package]; !found {
			if class.Package != "" {
				checked[class.Package] = true
				packages = append(packages, class.Package)
			}
		}
	}
	return packages
}
