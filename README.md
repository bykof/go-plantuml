# go-plantuml 
[![Build Status](https://github.com/bykof/go-plantuml/actions/workflows/test.yml/badge.svg)](https://github.com/bykof/go-plantuml/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/bykof/go-plantuml)](https://goreportcard.com/report/github.com/bykof/go-plantuml)

go-plantuml generates plantuml diagrams from go source files or directories.

## Installation

```bash
go install github.com/bykof/go-plantuml@latest
```

Please consider that `$GOPATH/bin` should be on our `$PATH`.


## Usage

```bash
Usage:
  go-plantuml [command]

Available Commands:
  generate    Generate a plantuml diagram from given paths
  help        Help about any command

Flags:
  -h, --help     help for go-plantuml
  -t, --toggle   Help message for toggle

Use "go-plantuml [command] --help" for more information about a command.

```

```bash
Usage:
  go-plantuml generate [flags]

Flags:
  -d, --directories strings   the go source directories (default [.])
  -f, --files strings         the go source files
  -h, --help                  help for generate
  -o, --out string            the graphfile (default "graph.puml")
  -r, --recursive             traverse the given directories recursively
```

## Example

For example we have to files in the directory `test`.

```go
// address.go
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
	return fmt.Sprintf(
		"%s %s %d", 
		PackageVariable, 
		AnotherPackageVariable, 
		StartingStreetNumber, 
	)
}

```

```go
// user.go
package models

import "github.com/bykof/go-plantuml/test/address/models"

type (
	User struct {
		FirstName      string
		LastName       string
		Age            uint8
		Address        *models.Address
		privateAddress models.Address
	}
)

func (user *User) SetFirstName(firstName string) {
	user.FirstName = firstName
}

func PackageFunction() string {
	return "Hello World"
}

```

Then we run `go-plantuml generate` or `go-plantuml generate -d . -o graph.puml`.

This will create a `graph.puml` file and check for .go files inside your current directory.

Which looks like this:
<p align="center">
  <img src="https://raw.githubusercontent.com/bykof/go-plantuml/master/docs/assets/graph.svg">
</p>

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)