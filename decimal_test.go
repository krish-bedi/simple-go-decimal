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

func TestNewAndString(t *testing.T) {
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
		{math.MinInt64, -4, "-922337203685477.5808"},
		{math.MaxInt64, -4, "922337203685477.5807"},
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
		cases := []struct{
			x 	 tinydecimal.Decimal
			y 	 tinydecimal.Decimal
			want string
			err	 error
		} {
			// 0.1 + 0.2 = 0.3
			{newDecimal(t, 1,-1), newDecimal(t, 2,-1), "0.3", nil},
			// + overflow: max + max = 2max
			{newDecimal(t, math.MaxInt64, -4), newDecimal(t, math.MaxInt64, -4), "", tinydecimal.ErrOverflow},
			// - overflow: (-max) + (-max) = -2max (2min)
			{newDecimal(t, math.MinInt64, -4), newDecimal(t, math.MinInt64, -4), "", tinydecimal.ErrOverflow},
			{newParse(t, "30.99"), newParse(t, "78.1234"), "109.1134", nil},
		}

		for _, c := range cases {
			addition, err := c.x.Add(c.y)

			if err != c.err {
				t.Fatalf("expected error: %v, got %v", c.err, err)
			}

			if err == nil {
				got := addition.String()
				assertEqual(t, got, c.want)
			}
		}
	})
	
	t.Run("subtraction", func(t *testing.T) {
		cases := []struct{
			x 	 tinydecimal.Decimal
			y 	 tinydecimal.Decimal
			want string
			err  error
		} {
			// 5.75 - 1.25 = 4.5
			{newDecimal(t, 575,-2), newDecimal(t, 125,-2), "4.5", nil},
			// + overflow: max - (-max) = 2max
			{newDecimal(t, math.MaxInt64, -4), newDecimal(t, math.MinInt64, -4), "", tinydecimal.ErrOverflow},
			// - overflow: -max - max = -2max (2min)
			{newDecimal(t, math.MinInt64, -4), newDecimal(t, math.MaxInt64, -4), "", tinydecimal.ErrOverflow},
			// MinInt64 = -9,223,...,808 whereas MaxInt64 = 9,223,...,807
			// Therefore, 0 - MinInt becomes +9,223,...,808 which doesn't fit in MaxInt64
			// And overflows by 1 back to MinInt64 (+ overflow)
			{newDecimal(t, 0, 0), newDecimal(t, math.MinInt64, -4), "", tinydecimal.ErrOverflow},
			{newParse(t, "25.1234"), newParse(t, "-100.55"), "125.6734", nil},
		}

		for _, c := range cases {
			subtraction, err := c.x.Sub(c.y)

			if err != c.err {
				t.Fatalf("expected error: %v, got %v", c.err, err)
			}

			if err == nil {
				got := subtraction.String()
				assertEqual(t, got, c.want)
			}
		}
	})
}

func TestMultiplyAndDivide(t *testing.T) {
	t.Run("test multiply", func(t *testing.T) {
		cases := []struct{
			x 	 tinydecimal.Decimal
			y 	 tinydecimal.Decimal
			want string
			err	 error
		} {
			// 55.1 * 100.5 = 5,537.55
			{newDecimal(t, 551,-1), newDecimal(t, 1005,-1), "5537.55", nil},
			// + overflow: max * 2 = 2max
			{newDecimal(t, math.MaxInt64, -4), newDecimal(t, 2, 0), "", tinydecimal.ErrOverflow},
			// - overflow: max * (-2) = -2max (2min)
			{newDecimal(t, math.MaxInt64, -4), newDecimal(t, -2, 0), "", tinydecimal.ErrOverflow},
			{newParse(t, "100.24"), newParse(t, "20.1203"), "2016.8588", nil},

		}

		for _, c := range cases {
			multiplication, err := c.x.Multiply(c.y)

			if err != c.err {
				t.Fatalf("expected error: %v, got %v", c.err, err)
			}

			if err == nil {
				got := multiplication.String()
				assertEqual(t, got, c.want)
			}
		}
	})

	t.Run("test divide", func(t *testing.T) {
		cases := []struct{
			x 	 tinydecimal.Decimal
			y 	 tinydecimal.Decimal
			want string
			err	 error
		} {
			// 10.5 / 0.12 = 87.5
			{newDecimal(t, 105,-1), newDecimal(t, 12,-2), "87.5", nil},
			// + overflow: max * 0.5 = 2max
			{newDecimal(t, math.MaxInt64, -4), newDecimal(t, 5,-1), "", tinydecimal.ErrOverflow},
			// - overflow: min * 0.5 = -2max (2min)
			{newDecimal(t, math.MinInt64, -4), newDecimal(t, 5,-1), "", tinydecimal.ErrOverflow},
			// division by 0
			{newDecimal(t, 123, -1), newDecimal(t, 0, 0), "", tinydecimal.ErrDivisionByZero},
			{newParse(t, "999.999"), newParse(t, "0.5678"), "1761.1817", nil},
		}

		for _, c := range cases {
			division, err := c.x.Divide(c.y)

			if err != c.err {
				t.Fatalf("expected error: %v, got %v", c.err, err)
			}

			if err == nil {
				got := division.String()
				assertEqual(t, got, c.want)
			}	
		}
	})
}

func TestParse(t *testing.T) {
	cases := []struct {
		input string
		want tinydecimal.Decimal
		err error
	} {
		{"123", newDecimal(t, 123, 0), nil},
		{"123.456", newDecimal(t, 123456, -3), nil},
		{"", newDecimal(t, 0, 0), nil},
		{"abc", newDecimal(t, 0, 0), tinydecimal.ErrInvalidFormat},
	}

	for _, c := range cases {
		got, err := tinydecimal.Parse(c.input)

		if err != c.err {
			t.Fatalf("expected error %v, but got %v", c.err, err)
		}

		if err == nil {
			assertEqual(t, got, c.want)
		}
	}
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func newDecimal(t *testing.T, value int64, exp int8) tinydecimal.Decimal {
	t.Helper()

	dec, err := tinydecimal.New(value, exp)
	if err != nil {
		t.Fatalf("received error during Decimal creation, %v", err)
	}

	return dec
}

func newParse(t *testing.T, v string) tinydecimal.Decimal {
	t.Helper()
	
	value, _ := tinydecimal.Parse(v)
	return value
}