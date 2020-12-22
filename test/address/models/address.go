package models

type (
	Address struct {
		Street     string
		City       string
		PostalCode string
		Country    string
		Bla        chan string
		A          *interface{}
	}
)

func (address Address) FullAddress(withPostalCode bool) string {
	return ""
}
