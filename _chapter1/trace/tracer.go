package trace

import (
	"io"
	"fmt"
)

// Tracer is the interface that describes an object capable of
// tracing events throuout the code.
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// Off creates a Tracer that will ignore call to Trace.
func Off() Tracer{
	return &nilTracer{}
}