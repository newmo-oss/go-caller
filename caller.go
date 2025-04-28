package caller // import "github.com/newmo-oss/go-caller"

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// StackTrace is stack of Frames.
type StackTrace []Frame

// New creates callers of [StackTrace].
// The argument skip is the number of stack frames to skip before recording,
// with 0 identifying the frame for [New] itself and 1 identifying the caller of [New].
// It mean that [New] calls [runtime.Callers] with 2+skip as the first argument.
func New(skip int) StackTrace {
	var pc [32]uintptr
	n := runtime.Callers(2+skip, pc[:])
	frames := runtime.CallersFrames(pc[:n])
	st := make(StackTrace, 0, n)
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		st = append(st, Frame{frame: frame})
	}
	return st
}

// Format implements [fmt.Formatter].
func (st StackTrace) Format(s fmt.State, verb rune) {
	fmt.Fprint(s, "[")
	for i, f := range st {
		if i > 0 {
			fmt.Fprint(s, " ")
		}
		f.Format(s, verb)
	}
	fmt.Fprint(s, "]")
}

// Frame is a stack of frame.
type Frame struct {
	frame runtime.Frame
}

// RuntimeFrame returns [runtime.Frame].
func (f Frame) RuntimeFrame() runtime.Frame {
	return f.frame
}

// FuncName returns function name without package path.
// e.g., example.com/sample/a.F.G.func1 -> F.G.func1
func (f Frame) FuncName() string {
	qualified := f.frame.Function
	idx := strings.LastIndex(f.frame.Function, "/")
	if idx >= 0 {
		qualified = qualified[idx:]
	}

	_, funcname, ok := strings.Cut(qualified, ".")
	if !ok {
		return qualified
	}
	return funcname
}

// PkgPath returns the package (import path) of the function.
// e.g., example.com/sample/a.F.G.func1 -> example.com/sample/a
func (f Frame) PkgPath() string {
	var pkgpath string
	qualified := f.frame.Function
	idx := strings.LastIndex(f.frame.Function, "/")
	if idx >= 0 {
		pkgpath = f.frame.Function[:idx]
		qualified = f.frame.Function[idx:]
	}
	pkgname, _, ok := strings.Cut(qualified, ".")
	if !ok {
		return pkgpath
	}
	return pkgpath + pkgname
}

// PkgName returns the package name of the function.
// e.g., example.com/sample/a.F.G.func1 -> a
func (f Frame) PkgName() string {
	qualified := f.frame.Function
	idx := strings.LastIndex(f.frame.Function, "/")
	if idx >= 0 {
		qualified = qualified[idx+1:]
	}
	pkgname, _, ok := strings.Cut(qualified, ".")
	if !ok {
		return ""
	}
	return pkgname
}

// File returns the file path.
func (f Frame) File() string {
	return f.frame.File
}

// Line returns the line number.
func (f Frame) Line() int {
	return f.frame.Line
}

// Format implements [fmt.Formatter].
//
// The verb means as follows.
//
//	verb	description	example
//	s	file name	a.go
//	s+	file path	example.com/sample/a/a.go
//	d	line number	10
//	n	function name	F.G.func1
//	P	package name	a
//	P+	import path	example.com/sample/a
//	v	file:line	a.go:10
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			fmt.Fprint(s, f.File())
		default:
			fmt.Fprint(s, filepath.Base(f.File()))
		}
	case 'd':
		fmt.Fprintf(s, "%d", f.Line())
	case 'n':
		fmt.Fprint(s, f.FuncName())
	case 'P': // 'p' (pointer) and 'T' (type) is reserved in fmt package
		switch {
		case s.Flag('+'):
			fmt.Fprint(s, f.PkgPath())
		default:
			fmt.Fprint(s, f.PkgName())
		}
	case 'v':
		f.Format(s, 's')
		fmt.Fprint(s, ":")
		f.Format(s, 'd')
	}
}
