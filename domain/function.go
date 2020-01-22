package domain

import "unicode"

type (
	Function struct {
		Name         string
		Parameters   Fields
		ReturnFields Fields
	}

	Functions []Function
)

func (function Function) IsPrivate() bool {
	if len(function.Name) > 0 {
		return unicode.IsLower(rune(function.Name[0]))
	}
	return false
}
