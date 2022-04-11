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

func (function Function) EqualImplementation(otherFunction Function) bool {
	return function.Name == otherFunction.Name && function.Parameters.EqualImplementations(
		otherFunction.Parameters,
	) && function.ReturnFields.EqualImplementations(
		otherFunction.ReturnFields,
	)
}
