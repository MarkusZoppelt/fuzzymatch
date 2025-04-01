package fuzzymatch

import (
	"slices"
	"testing"
)

func TestLevenshtein(t *testing.T) {
	tests := []struct {
		a, b     string
		expected int
	}{
		{"", "", 0},
		{"", "a", 1},
		{"a", "", 1},
		{"a", "a", 0},
		{"a", "b", 1},
		{"ab", "ab", 0},
		{"ab", "ac", 1},
		{"ab", "a", 1},
		{"ab", "abc", 1},
		{"llm-hint-file", "llm-hint-files", 1},
		{"llm-hints", "llm-hint-files", 5},
		{"Engine", "engine", 1},
	}

	for _, test := range tests {
		t.Run(test.a+"_"+test.b, func(t *testing.T) {
			got := Levenshtein(test.a, test.b)
			if got != test.expected {
				t.Errorf("levenshtein(%q, %q): expected %d, got %d", test.a, test.b, test.expected, got)
			}
		})
	}
}

func TestSuggestClosestMatch(t *testing.T) {
	tests := []struct {
		s               string
		possibleMatches []string
		minDist         int
		expected        string
	}{
		{"foo", []string{"foo", "bar", "baz"}, 9999, ""},
		{"fo", []string{"foo", "bar", "baz"}, 9999, "foo"},
		{"ba", []string{"bar", "baz", "ban"}, 9999, "bar"},
		{"qux", []string{"foo", "bar", "baz"}, 9999, "foo"},
		{"abcd", []string{"x", "y", "z"}, 9999, ""},
		{"", []string{"a", "b", "c"}, 9999, "a"},
		{"a", []string{""}, 9999, ""},
		{"", []string{""}, 9999, ""},
		{"test", []string{}, 9999, ""},
		{"llm-hint-file", []string{"llm-hint-files"}, 9999, "llm-hint-files"},
		{"llm-hintfile", []string{"llm-hint-files"}, 9999, "llm-hint-files"},
	}

	for _, test := range tests {
		t.Run(test.s, func(t *testing.T) {
			got := SuggestClosestMatch(test.s, test.possibleMatches, test.minDist)
			if got != test.expected {
				t.Errorf("SuggestClosestMatch(%q, %v, %d): expected %q, got %q", test.s, test.possibleMatches, test.minDist, test.expected, got)
			}
		})
	}
}

func FuzzSuggestClosestMatch(f *testing.F) {
	possibleMatches := []string{"foo", "bar", "baz", "foobar"}
	f.Add("foo")
	f.Add("fo")
	f.Add("ba")
	f.Add("qux")
	f.Add("abcd")
	f.Add("")
	f.Add("a")

	f.Fuzz(func(t *testing.T, s string) {
		minDist := 9999
		actualMinDist := 9999
		for _, p := range possibleMatches {
			dist := Levenshtein(s, p)
			if dist < actualMinDist {
				actualMinDist = dist
			}
		}

		returned := SuggestClosestMatch(s, possibleMatches, minDist)

		if slices.Contains(possibleMatches, s) {
			if returned != "" {
				t.Errorf("s=%q is in possibleMatches, expected \"\", got %q", s, returned)
			}
		} else {
			if actualMinDist <= 3 {
				if returned == "" {
					t.Errorf("s=%q, min distance %d ≤ 3, expected a suggestion, got \"\"", s, actualMinDist)
				} else if !slices.Contains(possibleMatches, returned) {
					t.Errorf("s=%q, returned %q not in possibleMatches", s, returned)
				} else if dist := Levenshtein(s, returned); dist > 3 {
					t.Errorf("s=%q, returned %q has distance %d > 3", s, returned, dist)
				} else if dist > actualMinDist {
					t.Errorf("s=%q, returned %q has distance %d, but min distance is %d", s, returned, dist, actualMinDist)
				}
			} else {
				if returned != "" {
					t.Errorf("s=%q, min distance %d > 3, expected \"\", got %q", s, actualMinDist, returned)
				}
			}
		}
	})
}
