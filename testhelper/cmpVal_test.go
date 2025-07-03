package testhelper

import (
	"testing"
)

func TestStringFirstDiff(t *testing.T) {
	testCases := []struct {
		ID
		s1           string
		s2           string
		expFirstDiff int
	}{
		{
			ID:           MkID("no diff: len(s2)"),
			s1:           "Hello",
			s2:           "Hello",
			expFirstDiff: 5,
		},
		{
			ID:           MkID("len s2 > len s1: len(s1)"),
			s1:           "Hello",
			s2:           "Hello, World",
			expFirstDiff: 5,
		},
		{
			ID:           MkID("len s1 > len s2: len(s2)"),
			s1:           "Hello, World",
			s2:           "Hello",
			expFirstDiff: 5,
		},
		{
			ID:           MkID("differ at first rune"),
			s1:           "Hello",
			s2:           "Bye",
			expFirstDiff: 0,
		},
		{
			ID:           MkID("differ at last rune"),
			s1:           "Hello",
			s2:           "Hella",
			expFirstDiff: 4,
		},
	}

	for _, tc := range testCases {
		fd := stringFirstDiff(tc.s1, tc.s2)
		DiffInt(t, tc.IDStr(), "firstDiff", fd, tc.expFirstDiff)
	}
}

func TestAlmostEqual(t *testing.T) {
	testCases := []struct {
		ID
		a, b, epsilon float64
		expected      bool
	}{
		{
			ID:       MkID("a == b, epsilon: 0"),
			a:        1.234,
			b:        1.234,
			epsilon:  0,
			expected: true,
		},
		{
			ID:       MkID("a == b, epsilon: 1"),
			a:        1.234,
			b:        1.234,
			epsilon:  1,
			expected: true,
		},
		{
			ID:       MkID("a != b, epsilon: 0"),
			a:        1.234,
			b:        2.345,
			epsilon:  0,
			expected: false,
		},
		{
			ID:       MkID("a != b, epsilon < diff"),
			a:        1.234,
			b:        1.235,
			epsilon:  0.0001,
			expected: false,
		},
		{
			ID:       MkID("a != b, epsilon > diff"),
			a:        1.234,
			b:        1.235,
			epsilon:  0.01,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := almostEqual(tc.a, tc.b, tc.epsilon)
			DiffBool(t, tc.IDStr(), "almostEqual", result, tc.expected)
		})
	}
}
