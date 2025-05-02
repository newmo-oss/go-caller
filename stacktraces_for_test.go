package caller_test

import "github.com/newmo-oss/go-caller"

// see: https://pkg.go.dev/cmd/compile#hdr-Compiler_Directives

//line callNew.go:1
func callNew(skip int) caller.StackTrace {
	return caller.New(skip)
}

//line wrapCallNew.go:1
func wrapCallNew(skip int) caller.StackTrace {
	return callNew(skip)
}
