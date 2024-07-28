// Package pipe provides utilities like pipe operator in the other language.
//
// # CAUTION
//
// The [Value] and [ValueErr] must be inlined because the [Pipe] uses the pointer to the function to store a intermediate state.
// See [Compiler Optimizations].
//
// [Compiler Optimizations]: https://go.dev/wiki/CompilerOptimizations
package pipe

import "unsafe"

type state struct {
	v   any
	err error
}

var states = make(map[uintptr]*state)

// Pipe represents the term of the pipeline.
type Pipe[T any] func(f func(T) T) Pipe[T]

// PipeErr represents the term of the pipeline with error.
type PipeErr[T any] func(f func(T) (T, error)) PipeErr[T]

// Value returns a term.
func Value[T any](v T) Pipe[T] {
	s := &state{v, nil}
	var f Pipe[T]
	f = func(g func(T) T) Pipe[T] {
		s.v = g(s.v.(T))
		return f
	}
	addr := **(**uintptr)(unsafe.Pointer(&f))
	states[addr] = s
	// TODO: runtime.SetFinalizer
	return f
}

// ValueErr returns a term.
func ValueErr[T any](v T, err error) PipeErr[T] {
	s := &state{v, err}
	var f PipeErr[T]
	f = func(g func(T) (T, error)) PipeErr[T] {
		if s.err == nil {
			s.v, s.err = g(s.v.(T))
		}
		return f
	}
	addr := **(**uintptr)(unsafe.Pointer(&f))
	states[addr] = s
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
	return s.v.(T)
}

// Value returns the result of the pipeline.
func (p PipeErr[T]) Value() (T, error) {
	addr := **(**uintptr)(unsafe.Pointer(&p))
	s := states[addr]
	if s == nil {
		panic("no state")
	}
	delete(states, addr)
	return s.v.(T), s.err
}

func From[T1, T2 any](p Pipe[T1], f func(T1) T2) Pipe[T2] {
	v := p.Value()
	return Value(f(v))
}

func FromWithErr[T1, T2 any](p PipeErr[T1], f func(T1) (T2, error)) PipeErr[T2] {
	v, err := p.Value()
	if err != nil {
		var zero T2
		return ValueErr(zero, err)
	}
	return ValueErr(f(v))
}

/*
If Go supports the type parameter on method, we will add Pipe.To method.

  pipe.Value(auth()).To(tokenFrom).To(fetch)
*/
