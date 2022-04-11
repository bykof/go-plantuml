package domain

import "strings"

type (
	Type string

	Types []Type
)

func (t Type) ToString() string {
	return string(t)
}

func (t Type) ToClassString() string {
	return strings.ReplaceAll(string(t), "*", "")
}
