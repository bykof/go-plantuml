package astParser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_dotNotatedModulePath(t *testing.T) {
	assert.Equal(
		t,
		"models",
		dotNotatedModulePath("address.go", "models"),
	)

	assert.Equal(
		t,
		"test.address.models",
		dotNotatedModulePath("test/address/models/address.go", "models"),
	)

	assert.Equal(
		t,
		"test.address.models",
		dotNotatedModulePath("test/address/models/models.go", "models"),
	)

	assert.Equal(
		t,
		"test.address.models.layer.models",
		dotNotatedModulePath("test/address/models/layer/models/address_layer.go", "models"),
	)
}
