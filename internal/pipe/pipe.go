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

// Value returns a term.
func Value[T any](v T) Pipe[T] {
	s := &state{v, nil}
	var f Pipe[T]
	f = func(g func(T) T) Pipe[T] {
		if s.err == nil {
			s.v = g(s.v.(T))
		}
		return f
	}
	addr := **(**uintptr)(unsafe.Pointer(&f))
	states[addr] = s
	// TODO: runtime.SetFinalizer
	return f
}

func (p Pipe[T]) withErr(err error) Pipe[T] {
	addr := **(**uintptr)(unsafe.Pointer(&p))
	s := states[addr]
	s.err = err
	return p
}

func (p Pipe[T]) Catch(f func(v T) (T, error)) Pipe[T] {
	var zero T

	v1, err := p.ValueErr()
	if err != nil {
		return Value(zero).withErr(err)
	}
	v2, err := f(v1)
	if err != nil {
		return Value(zero).withErr(err)
	}
	return Value(v2)
}

func (p Pipe[T]) Then(f func(T) T) Pipe[T] {
	return p(f)
}

// Value returns the result of the pipeline.
func (p Pipe[T]) Value() T {
	v, err := p.ValueErr()
	if err != nil {
		panic(err)
	}
	return v
}

// ValueErr returns the result of the pipeline.
func (p Pipe[T]) ValueErr() (T, error) {
	addr := **(**uintptr)(unsafe.Pointer(&p))
	s := states[addr]
	if s == nil {
		panic("no state")
	}
	delete(states, addr)
	return s.v.(T), s.err
}

func From[T1, T2 any](p Pipe[T1], f func(T1) (T2, error)) Pipe[T2] {
	var zero T2

	v1, err := p.ValueErr()
	if err != nil {
		return Value(zero).withErr(err)
	}
	v2, err := f(v1)
	if err != nil {
		return Value(zero).withErr(err)
	}
	return Value(v2)
}

func Errorable[T1, T2 any](f func(T1) T2) func(T1) (T2, error) {
	return func(v T1) (T2, error) {
		return f(v), nil
	}
}

/*
If Go supports the type parameter on method, we will add Pipe.To method.

  pipe.Value(auth()).To(tokenFrom).To(fetch)
*/
