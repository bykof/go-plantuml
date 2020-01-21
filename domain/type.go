package domain

type (
	Type string

	Types []Type
)

func (t Type) ToString() string {
	return string(t)
}
