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
