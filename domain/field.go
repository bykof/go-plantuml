package domain

import "unicode"

type (
	Field struct {
		Name     string
		Nullable bool
		Type     Type
	}

	Fields []Field
)

func (field Field) IsPrivate() bool {
	if len(field.Name) > 0 {
		return unicode.IsLower(rune(field.Name[0]))
	}
	return false
}

func (field Field) EqualImplementation(anotherField Field) bool {
	return field.Type == anotherField.Type && field.Nullable == anotherField.Nullable
}

func (fields Fields) EqualImplementations(otherFields Fields) bool {
	if len(fields) != len(otherFields) {
		return false
	}

	for i := 0; i < len(fields); i++ {
		if !fields[i].EqualImplementation(otherFields[i]) {
			return false
		}
	}

	return true
}
