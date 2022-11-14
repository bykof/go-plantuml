package calculator

import (
	"github.com/bykof/go-plantuml/test/user/models"
)


type Number interface {
	~int | ~float32
}

type Calculator[T Number] struct {
	owner models.User
}

func (c *Calculator[T]) Add(a, b T) T {
	return a + b
}

func (c *Calculator[T]) Multiply(a, b T) T {
	return a * b
}
