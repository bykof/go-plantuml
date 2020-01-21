package domain

import "unicode"

type (
	Field struct {
		Name     string
		Nullable bool
		Package  Package
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
