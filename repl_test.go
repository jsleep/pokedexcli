package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		// add more cases here
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(c.expected) {
			t.Errorf("Expected slice of len %d, got %d", len(c.expected), len(actual))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Expected %s, got %s", expectedWord, word)
			}
		}
	}
}
