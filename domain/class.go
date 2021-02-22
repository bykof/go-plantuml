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

func (classes Classes) ClassByName(name string) *Class {
	if classes == nil {
		return nil
	}

	for _, class := range classes {
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

func (classes Classes) ClassIndexByPointerName(pointerName string) int {
	if pointerName[0] == '*' {
		return classes.ClassIndexByName(pointerName[1:])
	}
	return -1
}
