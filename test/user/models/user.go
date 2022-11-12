package models

import "github.com/bykof/go-plantuml/test/address/models"

type (
	User struct {
		FirstName      string
		LastName       string
		Age            uint8
		Address        *models.Address
	}
)

func (user *User) SetFirstName(firstName string) {
	user.FirstName = firstName
}

func PackageFunction() string {
	return "Hello World"
}
