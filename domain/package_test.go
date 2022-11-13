package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackages_AllClasses(t *testing.T) {
	classes1 := Classes{Class{Name: "One"}}
	classes2 := Classes{Class{Name: "Second"}}
	packages := Packages{Package{Classes: classes1}, Package{Classes: classes2}}
	assert.Equal(t, append(classes1, classes2...), packages.AllClasses())
}

func TestPackages_AllInterfaces(t *testing.T) {
	interfaces1 := Interfaces{Interface{Name: "One"}}
	interfaces2 := Interfaces{Interface{Name: "Second"}}
	packages := Packages{Package{Interfaces: interfaces1}, Package{Interfaces: interfaces2}}
	assert.Equal(t, append(interfaces1, interfaces2...), packages.AllInterfaces())
}
