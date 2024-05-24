# stdlibx-go

Extensions to [golang](https://go.dev/) stdlib. Zero dependencies.

## Usage

```go
package app

import "github.com/ahawker/stdlibx-go/pkg"


func main() {
	stdlibx.Must[*string](nil) // panics
}
```

## Local Development

```shell
$ make help
modules                        Tidy and vendor Go modules for local development.
test                           Run tests.
test-benchmark                 Run benchmark tests.
```

## License

[Apache 2.0](LICENSE)
