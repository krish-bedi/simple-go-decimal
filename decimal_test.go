package tinydecimal_test

import (
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
	}

	for _, c := range cases {
		decimal, err := tinydecimal.New(c.value, c.exponent)
		if err != nil {
			t.Fatalf("received unexpected error: %v", err)
		}

		got := decimal.String()

		if got != c.want {
			t.Errorf("for case: %+v, got %v", c, got)
		}
	}
}