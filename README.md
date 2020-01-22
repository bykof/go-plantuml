# go-plantuml 
[![Build Status](https://travis-ci.com/bykof/go-plantuml.svg?branch=master)](https://travis-ci.com/bykof/go-plantuml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bykof/go-plantuml)](https://goreportcard.com/report/github.com/bykof/go-plantuml)

go-plantuml generates plantuml diagrams from go source files or directories.

## Installation

```bash
go get -u github.com/bykof/go-plantuml
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

## Example

For example we have to files in the directory `testGraph`.

```go
// addres.go
package testGraph

type (
	Address struct {
		Street     string
		City       string
		PostalCode string
		Country    string
	}
)

func (address Address) FullAddress(withPostalCode bool) string {
    return ""
}
```

```go
// user.go
package testGraph

type (
	User struct {
		FirstName      string
		LastName       string
		age            uint8
		Address        Address
		privateAddress Address
	}
)

func (user *User) SetFirstName(firstName string) {}
```

Then we run `go-plantuml generate` or `go-plantuml generate -d . -o graph.puml`.

This will create a `graph.puml` file and check for .go files inside your current directory.

Which looks like this:
<p align="center">
  <img src="https://raw.githubusercontent.com/bykof/go-plantuml/master/docs/assets/graph.png">
</p>

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)