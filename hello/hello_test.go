package main

import "testing"

func TestSubstract(t *testing.T) {

	if Substract(35, 12) != 23 {
		t.Error("Expected result to equal 46")
	}
}

func TestTableSubstract(t *testing.T) {
	var tests = []struct {
		x        int
		y        int
		expected int
	}{
		{4, 6, -2},
		{10, -2, 12},
		{-20, 1000, -1020},
		{1000, 210, 790},
	}

	for _, test := range tests {
		if output := Substract(test.x, test.y); output != test.expected {
			t.Error("Test Failed: {} inputted, {} expected, recieved: {}", test.x, test.y, test.expected, output)
		}
	}
}
