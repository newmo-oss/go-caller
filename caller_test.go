package caller_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/newmo-oss/go-caller"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// see: stacktraces_for_test.go
	cases := map[string]struct {
		stacktrace caller.StackTrace
		want       string
	}{
		"skip 0": {stacktrace: callNew(0), want: "callNew.go:2 github.com/newmo-oss/go-caller_test.callNew"},
		"skip 1": {stacktrace: wrapCallNew(1), want: "wrapCallNew.go:2 github.com/newmo-oss/go-caller_test.wrapCallNew"},
		"skip 0 - wrapped": {stacktrace: wrapCallNew(0), want: "callNew.go:2 github.com/newmo-oss/go-caller_test.callNew"},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if len(tt.stacktrace) == 0 {
				t.Fatal("stacktrace must not be nil or empty")
			}

			frame := tt.stacktrace[0]
			got := fmt.Sprintf("%s:%d %s.%s", frame.File(), frame.Line(), frame.PkgPath(), frame.FuncName())
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Error("(-got, +want) =", diff)
			}
		})
	}
}

func TestStackTrance_Format(t *testing.T) {
	t.Parallel()

	st := caller.StackTrace([]caller.Frame{
		caller.NewFrameForTest(runtime.Frame{
			Function: "example.com/sample/a.F.G.func1",
			File:     "example.com/sample/a/a.go",
			Line:     10,
		}),
		caller.NewFrameForTest(runtime.Frame{
			Function: "example.com/sample/b.F.G.func2",
			File:     "example.com/sample/b/b.go",
			Line:     11,
		}),
	})

	cases := []struct {
		format string
		want   string
	}{
		{"%s", "[a.go b.go]"},
		{"%+s", "[example.com/sample/a/a.go example.com/sample/b/b.go]"},
		{"%d", "[10 11]"},
		{"%n", "[F.G.func1 F.G.func2]"},
		{"%P", "[a b]"},
		{"%+P", "[example.com/sample/a example.com/sample/b]"},
		{"%v", "[a.go:10 b.go:11]"},
	}

	for _, tt := range cases {
		name := fmt.Sprintf("%q->%q", tt.format, tt.want)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := fmt.Sprintf(tt.format, st); got != tt.want {
				t.Errorf("Format does not match: (got, want) = (%q, %q)", got, tt.want)
			}
		})
	}
}

func TestStackTrance_RuntimeFrame(t *testing.T) {
	t.Parallel()

	want := runtime.Frame{
		Function: "example.com/sample/a.F.G.func1",
		File:     "example.com/sample/a/a.go",
		Line:     10,
	}

	if got := caller.NewFrameForTest(want).RuntimeFrame(); got != want {
		t.Errorf("RuntimeFrame does not match: (got, want) = (%v, %v)", got, want)
	}
}

func TestFrame_FuncName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		frameFunction string
		want          string
	}{
		{"F", "F"},
		{"a.F", "F"},
		{"example.com/sample/a.F", "F"},
		{"example.com/sample/a.F.G.func1", "F.G.func1"},
	}

	for _, tt := range cases {
		name := fmt.Sprintf("%q->%q", tt.frameFunction, tt.want)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			f := caller.NewFrameForTest(runtime.Frame{
				Function: tt.frameFunction,
			})

			if got := f.FuncName(); got != tt.want {
				t.Errorf("FuncName does not match: (got, want) = (%q, %q)", got, tt.want)
			}
		})
	}
}

func TestFrame_PkgPath(t *testing.T) {
	t.Parallel()

	cases := []struct {
		frameFunction string
		want          string
	}{
		{"F", ""},
		{"a.F", "a"},
		{"example.com/sample/a.F", "example.com/sample/a"},
		{"example.com/sample/a.F.G.func1", "example.com/sample/a"},
	}

	for _, tt := range cases {
		name := fmt.Sprintf("%q->%q", tt.frameFunction, tt.want)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			f := caller.NewFrameForTest(runtime.Frame{
				Function: tt.frameFunction,
			})

			if got := f.PkgPath(); got != tt.want {
				t.Errorf("PkgPath does not match: (got, want) = (%q, %q)", got, tt.want)
			}
		})
	}
}

func TestFrame_PkgName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		frameFunction string
		want          string
	}{
		{"F", ""},
		{"a.F", "a"},
		{"example.com/sample/a.F", "a"},
		{"example.com/sample/a.F.G.func1", "a"},
	}

	for _, tt := range cases {
		name := fmt.Sprintf("%q->%q", tt.frameFunction, tt.want)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			f := caller.NewFrameForTest(runtime.Frame{
				Function: tt.frameFunction,
			})

			if got := f.PkgName(); got != tt.want {
				t.Errorf("PkgName does not match: (got, want) = (%q, %q)", got, tt.want)
			}
		})
	}
}

func TestFrame_Format(t *testing.T) {
	t.Parallel()

	frame := caller.NewFrameForTest(runtime.Frame{
		Function: "example.com/sample/a.F.G.func1",
		File:     "example.com/sample/a/a.go",
		Line:     10,
	})

	cases := []struct {
		format string
		want   string
	}{
		{"%s", "a.go"},
		{"%+s", "example.com/sample/a/a.go"},
		{"%d", "10"},
		{"%n", "F.G.func1"},
		{"%P", "a"},
		{"%+P", "example.com/sample/a"},
		{"%v", "a.go:10"},
	}

	for _, tt := range cases {
		name := fmt.Sprintf("%q->%q", tt.format, tt.want)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := fmt.Sprintf(tt.format, frame); got != tt.want {
				t.Errorf("Format does not match: (got, want) = (%q, %q)", got, tt.want)
			}
		})
	}
}
