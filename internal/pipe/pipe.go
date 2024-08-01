// Package pipe provides utilities like pipe operator in the other language.
//
// # CAUTION
//
// The [Value] and [From] must be inlined because the [Pipe] uses the pointer to the function to store a intermediate state.
// See [Compiler Optimizations].
//
// [Compiler Optimizations]: https://go.dev/wiki/CompilerOptimizations
package pipe

import "unsafe"

type state struct {
	next evaluator
}

type result struct {
	v   any
	err error
}

type evaluator interface {
	eval() *result
}

type scalar[T any] struct {
	v T
}

func (s *scalar[T]) eval() *result {
	return &result{s.v, nil}
}

type selection[In, Out any] struct {
	parent evaluator
	fn     func(v In) Out
}

func (s *selection[In, Out]) eval() *result {
	r := s.parent.eval()
	if r.err != nil {
		var zero In
		return &result{zero, r.err}
	}
	v := r.v.(In)
	return &result{s.fn(v), nil}
}

type selection2[In, Out any] struct {
	parent evaluator
	fn     func(in In) (Out, error)
}

func (s *selection2[In, Out]) eval() *result {
	r1 := s.parent.eval()
	if r1.err != nil {
		var zero Out
		return &result{zero, r1.err}
	}
	v := r1.v.(In)
	r2, err := s.fn(v)
	return &result{r2, err}
}

var states = make(map[uintptr]*state)

// Pipe represents the term of the pipeline.
type Pipe[T any] func(f func(T) T) Pipe[T]

// Value returns a term.
func Value[T any](v T) Pipe[T] {
	s := &state{&scalar[T]{v}}
	var f Pipe[T]
	f = func(g func(T) T) Pipe[T] {
		n := s.next
		s.next = &selection[T, T]{n, g}
		return f
	}
	addr := **(**uintptr)(unsafe.Pointer(&f))
	states[addr] = s
	// TODO: runtime.SetFinalizer
	return f
}

func (p Pipe[T]) Catch(f func(v T) (T, error)) Pipe[T] {
	addr := **(**uintptr)(unsafe.Pointer(&p))
	s := states[addr]
	n := s.next
	s.next = &selection2[T, T]{n, f}
	return p
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

	r := s.next.eval()
	if r.err != nil {
		var zero T
		return zero, r.err
	}
	return r.v.(T), nil
}

func From[In, Out any](p Pipe[In], f func(In) (Out, error)) Pipe[Out] {
	addr := **(**uintptr)(unsafe.Pointer(&p))
	s := states[addr]
	delete(states, addr)
	n := s.next
	s.next = &selection2[In, Out]{n, f}

	var p2 Pipe[Out]
	p2 = func(g func(Out) Out) Pipe[Out] {
		n := s.next
		s.next = &selection[Out, Out]{n, g}
		return p2
	}
	addr = **(**uintptr)(unsafe.Pointer(&p2))
	states[addr] = s
	return p2
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
