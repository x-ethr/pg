# `pg`

## Documentation

Official `godoc` documentation (with examples) can be found at the [Package Registry](https://pkg.go.dev/github.com/x-ethr/pg).

## Usage

###### Add Package Dependency

```bash
go get -u github.com/x-ethr/pg
```

###### Import & Implement

`main.go`

```go
package main

import (
    "fmt"

    "github.com/x-ethr/pg"
)

func main() {
    ctx, level := context.Background(), slog.LevelInfo

    uri := pg.DSN()
    connection, e := database.Connection(ctx, uri)
    if e != nil {
        panic(e)
    }
}
```

- Please refer to the [code examples](./example_test.go) for additional usage and implementation details.
- See https://pkg.go.dev/github.com/x-ethr/environment for additional documentation.

## Contributions

See the [**Contributing Guide**](./CONTRIBUTING.md) for additional details on getting started.
