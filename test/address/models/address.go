package models

import "fmt"

var (
	EmptyVariable, AnotherEmptyVariable string
	A, B                                = "1", 2
	PackageVariable                     = "Teststreet"
	AnotherPackageVariable              = "Anotherteststreet"
)

const (
	StartingStreetNumber = 1
)

type (
	AddressLike interface {
		FullAddress(withPostalCode bool) string
	}

	Address struct {
		A, B          string
		Street        string
		City          string
		PostalCode    string
		Country       string
		CustomChannel chan string
		AnInterface   *interface{}
	}
)

func (address Address) FullAddress(withPostalCode bool) string {
	return fmt.Sprintf("%s %s %d", PackageVariable, AnotherPackageVariable, StartingStreetNumber)
}
