package ptime

import (
	"testing"
)

func TestLastMonthString(t *testing.T) {
	lm, _ := LastMonthString("2022-02")
	expected := "2022-01"
	if lm != expected {
		t.Errorf("got %q, wanted %q", lm, expected)
	}
}

func TestLastMonthStringJan(t *testing.T) {
	lm, _ := LastMonthString("2022-01")
	expected := "2021-12"
	if lm != expected {
		t.Errorf("got %q, wanted %q", lm, expected)
	}
}
