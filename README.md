# go-caller [![Go Reference](https://pkg.go.dev/badge/github.com/newmo-oss/go-caller.svg)](https://pkg.go.dev/github.com/newmo-oss/go-caller)[![Go Report Card](https://goreportcard.com/badge/github.com/newmo-oss/go-caller)](https://goreportcard.com/report/github.com/newmo-oss/go-caller)

go-caller is a library of stack trace.

## Usage

```go
package main

import (
	"fmt"

	"github.com/newmo-oss/go-caller"
)

func main() {
	stacktrace := caller.New(1)
	frame := stacktrace[0]
	// main.go:10 main.main
	fmt.Printf("%s:%d %s.%s\n", frame.File(), frame.Line(), frame.PkgName(), frame.FuncName())
}
```

## License
MIT
