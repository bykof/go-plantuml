package domain

import (
	"strings"
)

type (
	Class struct {
		Name      string
		Fields    Fields
		Functions Functions
	}

	Classes []Class
)

func (class Class) HasRelation(toClass Class) bool {
	for _, field := range class.Fields {
		fieldTypeName := string(field.Type)

		for _, prefix := range []string{"*", "[]"} {
			fieldTypeName = strings.TrimPrefix(fieldTypeName, prefix)
		}

		if fieldTypeName == toClass.Name {
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
