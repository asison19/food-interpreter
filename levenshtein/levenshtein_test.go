package levenshtein

import (
	"testing"
)

type ld struct {
	a string
	b string
	n int
}

func TestLevenshteinDistance(t *testing.T) {
	var l = []ld{
		{"kitten", "sitting", 3},
		{"sitting", "kitten", 3},
		{"kitten", "sittin", 2},
		{"apple", "banana", 5},
		{"abc", "", 3},
		{"", "", 0},
		{"12345", "67890", 5},
		{",./;'[]\\", "abcdefg", 8},
		{"<>?:\"{}|", "abcdefg", 8},
		{"`~!@#$%^&*()_+-=", "abcdefg", 16},
		{"ABC", "abc", 3}, // TODO should we care about casing?
		{"ABc", "abc", 2},
		{"Who's on first?", "What's on second?", 8},
	}

	for _, e := range l {
		r := LevenshteinDistance(e.a, e.b)
		if r != e.n {
			t.Fatalf("Expected Levenshtein Distance of %d for \"%s\" and \"%s\". Got %d instead.", e.n, e.a, e.b, r)
		}
	}
}
