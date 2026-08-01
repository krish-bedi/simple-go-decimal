package tinydecimal_test

import (
	"testing"
	"testing/quick"

	tinydecimal "github.com/krish-bedi/simple-go-decimal"
)

// Property: a + b == b + a
func TestAddIsCommutative(t *testing.T) {
	property := func(x, y int64, e, f int8) bool {
		a, errA := tinydecimal.New(x, modExp(e))
		if errA != nil {
			return true // error returns 0 as value
		}

		b, errB := tinydecimal.New(y, modExp(f))
		if errB != nil {
			return true // error returns 0 as value
		}

		ab, errAB := a.Add(b)
		ba, errBA := b.Add(a)

		return ab == ba && errAB == errBA
	}

	if err := quick.Check(property, nil); err != nil {
		t.Error(err)
	}
}

// Property: (a + b) - b == a
func TestAddThenSub(t *testing.T) {
	property := func(x, y int64, e, f int8) bool {
		a, errA := tinydecimal.New(x, modExp(e))
		if errA != nil {
			return true // error returns 0 as value
		}

		b, errB := tinydecimal.New(y, modExp(f))
		if errB != nil {
			return true // error returns 0 as value
		}

		ab, errAB := a.Add(b)
		// skip cases that overflow on a + b
		if errAB != nil {
			return true
		}
		subB, errSubB := ab.Sub(b)
		// (a + b) - b should not overflow
		if errSubB != nil {
			return false
		}

		return subB == a
	}

	if err := quick.Check(property, nil); err != nil {
		t.Error(err)
	}
}

// Property: (a * b) == (b * a)
func TestMultiplyIsCommutative(t *testing.T) {
	property := func(x, y int64, e, f int8) bool {
		a, errA := tinydecimal.New(x, modExp(e))
		if errA != nil {
			return true
		}

		b, errB := tinydecimal.New(y, modExp(f))
		if errB != nil {
			return true
		}

		ab, errAB := a.Multiply(b)
		if errAB != nil {
			return true // skip cases where a * b overflows
		}

		ba, errBA := b.Multiply(a)
		if errBA != nil {
			return false // b * a should not overflow when a * b didn't
		}

		return ab == ba
	}

	if err := quick.Check(property, nil); err != nil {
		t.Error(err)
	}
}

// Property: a/a == 1
func TestDivideByItself(t *testing.T) {
	property := func(x int64, e int8) bool {
		a, errA := tinydecimal.New(x, modExp(e))
		if errA != nil {
			return true // failed to initialize a
		}

		one, _ := tinydecimal.New(1, 0)

		div, errDiv := a.Divide(a)
		if errDiv != nil {
			return false // division by self should not overflow
		}

		return div == one
	}

	if err := quick.Check(property, nil); err != nil {
		t.Error(err)
	}
}

// get exponents in range of -4 to +4
func modExp(e int8) int8 {
	return e % 5
}