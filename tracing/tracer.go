package trace

import (
	"fmt"
	"io"
)

// Tracer is an interface that describes an object capable of
// tracing events throughout the code

type Tracer interface {
	Trace(...interface{})
}

type trace struct {
	out io.Writer
}

type nilTracer struct {}

func (t *trace) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

func (t *nilTracer) Trace(a ...interface{}) {}

func New(w io.Writer) Tracer {
	return &trace{out: w}
}

// Off creates a Tracer that will ignore calls to Trace.
func Off() Tracer {
	return &nilTracer{}
}
