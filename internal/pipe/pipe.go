// Package pipe provides utilities like pipe operator in the other language.
//
// # Warning
//
// The [Value] function must be inlined because it uses a pointer to the function to keep a state each pipeline.
// See [Compiler Optimizations].
//
// [Compiler Optimizations]: https://go.dev/wiki/CompilerOptimizations
package pipe

import (
	"unsafe"
)

type state struct {
	v1, v2 any
}

var states = make(map[uintptr]*state)

// Pipe represents the term of the pipeline.
type Pipe[T any] func(f func(T) T) Pipe[T]

// Value returns a term.
func Value[T any](v T) Pipe[T] {
	s := &state{v1: v}
	var f Pipe[T]
	f = func(g func(T) T) Pipe[T] {
		v1 := s.v1.(T)
		s.v1 = g(v1)
		return f
	}
	addr := **(**uintptr)(unsafe.Pointer(&f))
	states[addr] = s
	// TODO: runtime.SetFinalizer
	return f
}

// Value returns the result of the pipeline.
func (p Pipe[T]) Value() T {
	addr := **(**uintptr)(unsafe.Pointer(&p))
	s := states[addr]
	if s == nil {
		panic("no state")
	}
	delete(states, addr)
	return s.v1.(T)
}

func From[T1, T2 any](p Pipe[T1], f func(T1) T2) Pipe[T2] {
	v := p.Value()
	return Value(f(v))
}

/*
If Go supports the type parameter on method, we will add Pipe.To method.

  pipe.Value(auth()).To(tokenFrom).To(fetch)
*/
