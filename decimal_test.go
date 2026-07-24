package tinydecimal_test

import (
	"math"
	"testing"

	tinydecimal "github.com/krish-bedi/simple-go-decimal"
)

func TestFloat(t *testing.T) {
	a := 0.1
	b := 0.2

	got := a + b
	want := 0.3
	
	// Expected to fail
	if got == want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestDecimalAndString(t *testing.T) {
	cases := []struct {
		value int64
		exponent int8
		want string
	} {
		{123, 0, "123"},
		{123, -2, "1.23"},
		{5, -1, "0.5"},
		{0, 0, "0"},
		{-7, -4, "-0.0007"},
		{-123, -2, "-1.23"},
		{-123, 0, "-123"},
		{-7000, -4, "-0.7"},
		{7, -4, "0.0007"},
		{7000, -4, "0.7"},
		{-120, 0, "-120"},
		{120, 4, "1200000"},
		{70, -5, "0.0007"},
	}

	for _, c := range cases {
		decimal, err := tinydecimal.New(c.value, c.exponent)
		if err != nil {
			t.Fatalf("received unexpected error: %v", err)
		}

		got := decimal.String()

		assertEqual(t, got, c.want)
	}
}

func TestOverflowAndPrecisionLoss(t *testing.T) {
	t.Run("test positive overflow", func(t *testing.T) {
		// value too large to store as scaled integer
		_, err := tinydecimal.New(math.MaxInt64, -2)
		if err != tinydecimal.ErrOverflow {
			t.Fatalf("expected ErrOverflow, got %v", err)
		}
	})
	t.Run("test negative overflow", func(t *testing.T) {
		// value too small to store as scaled integer
		// -2 exp multiplies MinInt64 by 10 which is < MinInt64
		_, err := tinydecimal.New(math.MinInt64, -2)
		if err != tinydecimal.ErrOverflow {
			t.Fatalf("expected ErrOverflow, got %v", err)
		}
	})
	t.Run("test precision loss", func(t *testing.T) {
		// can not store 0.00007 with a max precision of 0.0001
		_, err := tinydecimal.New(7, -5)
		if err != tinydecimal.ErrPrecisionLoss {
			t.Fatalf("expected ErrPrecisionLoss, got %v", err)
		}
	})
} 

func TestAddAndSub(t *testing.T) {
	t.Run("addition", func(t *testing.T) {
		// 0.1 + 0.2 = 0.3
		a, _ := tinydecimal.New(1, -1) // 0.1
		b, _ := tinydecimal.New(2, -1) // 0.2

		got := a.Add(b).String()
		want := "0.3"

		assertEqual(t, got, want)
	})
	
	t.Run("subtraction", func(t *testing.T) {
		// 5.75 - 1.25 = 4.5
		a, _ := tinydecimal.New(575, -2)
		b, _ := tinydecimal.New(125, -2)

		got := a.Sub(b).String()
		want := "4.5"

		assertEqual(t, got, want)
	})
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}