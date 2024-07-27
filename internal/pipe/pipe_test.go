package pipe

import (
	"testing"
)

func TestValue(t *testing.T) {
	v := Value(10).Value()
	if v != 10 {
		t.Errorf("Value() = %d; want 10", v)
	}
}

func TestValue_1x1(t *testing.T) {
	v1 := Value(10)
	v2 := Value(20)
	r2 := v2.Value()
	r1 := v1.Value()
	if r1 != 10 {
		t.Errorf("v1.Value() = %d; want 10", r1)
	}
	if r2 != 20 {
		t.Errorf("v2.Value() = %d; want 20", r2)
	}
}

func TestPipe(t *testing.T) {
	plus3 := func(n int) int { return n + 3 }
	times10 := func(n int) int { return n * 10 }
	n := Value(10)(plus3)(times10).Value()
	if n != 130 {
		t.Errorf("Value() = %d; want 130", n)
	}
}
