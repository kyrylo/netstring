# Netstring

[![.github/workflows/test.yml](https://github.com/kyrylo/netstring/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/kyrylo/netstring/actions/workflows/test.yml)

- [Netstring README][netstring-github]
- [pkg.go.dev documentation][docs]

## Introduction

_Netstring_ is a library for packing and parsing [netstrings][netstring],
self-delimiting encoding of strings. The library is extremely simple and well-tested.

Netstrings may be used as a basic building block for reliable network protocols.
Most high-level protocols, in effect, transmit a sequence of strings; those
strings may be encoded as netstrings and then concatenated into a sequence of
characters, which in turn may be transmitted over a reliable stream protocol
such as TCP.

## Installation

### Go modules

Netstring can be installed like any other Go package that supports [Go
modules][go-mod].

#### Installing in an existing project

Just `go get` the library:

```sh
go get github.com/kyrylo/netstring
```

## Example

### Parsing a netstring into a byte string

```go
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"

	"github.com/kyrylo/netstring"
)

func main() {
	// The netstring is "8:sunshine,"
	netstr := []byte{
		0x38, 0x3a, 0x73, 0x75, 0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65, 0x2c,
	}

	// Create a reader.
	buf := bufio.NewReader(bytes.NewReader(netstr))

	// Parse the "8:sunshine," into "sunshine".
	str, err := netstring.Parse(buf)
	if err != nil {
		log.Fatal(err)
	}

	// Output: "sunshine"
	fmt.Printf("Input netstring: %s\n", netstr)
	fmt.Printf("  Parsed string: %s\n", str)
}
```

### Packing a byte string into a netstring

```go
package main

import (
	"fmt"
	"log"

	"github.com/kyrylo/netstring"
)

func main() {
	s := []byte("sunshine")
	netstr, err := netstring.Pack(s)
	if err != nil {
		log.Fatal(err)
	}

	// netstr is "8:sunshine,"
	// bytes: [0x38, 0x3a, 0x73, 0x75, 0x6e, 0x73, 0x68, 0x69, 0x6e, 0x65, 0x2c]
	fmt.Printf("    Input string: %s\n", s)
	fmt.Printf("Output netstring: %s\n", netstr)
}
```

## Supported Go versions

The library supports Go v1.17+. The CI file would be the best source of truth
because it contains all Go versions that are tested against.

## Contact

In case you have a problem, question or a bug report, feel free to:

- [file an issue][issues]
- [tweet at me][twitter]

## License

The project uses the MIT License. See LICENSE.md for details.

[netstring-github]: https://github.com/kyrylo/netstring
[netstring]: https://cr.yp.to/proto/netstrings.txt
[semver2]: http://semver.org/spec/v2.0.0.html
[go-mod]: https://github.com/golang/go/wiki/Modules
[issues]: https://github.com/kyrylo/netstring/issues
[twitter]: https://twitter.com/kyrylosilin
[docs]: https://pkg.go.dev/github.com/kyrylo/netstring
