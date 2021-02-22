package models

type (
	AddressLike interface {
		FullAddress(withPostalCode bool) string
	}

	Address struct {
		Street        string
		City          string
		PostalCode    string
		Country       string
		CustomChannel chan string
		AnInterface   *interface{}
	}
)

func (address Address) FullAddress(withPostalCode bool) string {
	return ""
}
