Netstring
=========

* [Netstring README][netstring-github]
* [pkg.go.dev documentation][docs]

Introduction
------------

_Netstring_ is a library for packing and parsing [netstrings][netstring],
self-delimiting encoding of strings. The library is extremely simple and well-tested.

Netstrings may be used as a basic building block for reliable network protocols.
Most high-level protocols, in effect, transmit a sequence of strings; those
strings may be encoded as netstrings and then concatenated into a sequence of
characters, which in turn may be transmitted over a reliable stream protocol
such as TCP.

Installation
------------

### Go modules

Gobrake can be installed like any other Go package that supports [Go
modules][go-mod].

#### Installing in an existing project

Just `go get` the library:

```sh
go get github.com/kyrylo/netstring
```

Example
-------

### Parsing a byte string

```go
import (
	"log"

	"github.com/kyrylo/netstring"
)

func main() {
	s := "8:sunshine,"
	b := bufio.NewReader(bytes.NewReader([]byte(s)))

	parsed, err := netstring.Parse(b)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(parsed) // "sunshine"
}
```

### Packing bytes into a netstring

```go
import (
    "log"

    "github.com/kyrylo/netstring"
)

func main() {
	s := "sunshine"
    packed := netstring.Pack([]byte(s))
    fmt.Println(packed) // "sunshine"
}
```

Supported Go versions
---------------------

The library supports Go v1.11+. The CI file would be the best source of truth
because it contains all Go versions that are tested against.

Contact
-------

In case you have a problem, question or a bug report, feel free to:

* [file an issue][issues]
* [tweet at me][twitter]

License
-------

The project uses the MIT License. See LICENSE.md for details.

[netstring-github]: https://github.com/kyrylo/netstring
[netstring]: https://cr.yp.to/proto/netstrings.txt
[semver2]: http://semver.org/spec/v2.0.0.html
[go-mod]: https://github.com/golang/go/wiki/Modules
[issues]: https://github.com/kyrylo/netstring/issues
[twitter]: https://twitter.com/kyrylosilin
[docs]: https://pkg.go.dev/github.com/kyrylo/netstring
