package testhelper_test

import (
	"testing"

	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestStringSliceDiff(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		s1, s2  []string
		expDiff bool
	}{
		{
			ID: testhelper.MkID("both nil"),
		},
		{
			ID:      testhelper.MkID("first nil, second not nil"),
			s2:      []string{"a"},
			expDiff: true,
		},
		{
			ID:      testhelper.MkID("first not nil, second nil"),
			s1:      []string{"a"},
			expDiff: true,
		},
		{
			ID:      testhelper.MkID("neither nil, first shorter than second"),
			s1:      []string{"a"},
			s2:      []string{"a", "b"},
			expDiff: true,
		},
		{
			ID:      testhelper.MkID("neither nil, second shorter than first"),
			s1:      []string{"a", "b"},
			s2:      []string{"a"},
			expDiff: true,
		},
		{
			ID:      testhelper.MkID("same length, differing"),
			s1:      []string{"a", "b"},
			s2:      []string{"a", "c"},
			expDiff: true,
		},
		{
			ID:      testhelper.MkID("same length, same strings but reversed"),
			s1:      []string{"a", "b"},
			s2:      []string{"b", "a"},
			expDiff: true,
		},
		{
			ID: testhelper.MkID("same length, same strings"),
			s1: []string{"a", "b"},
			s2: []string{"a", "b"},
		},
	}

	for _, tc := range testCases {
		differs := testhelper.StringSliceDiff(tc.s1, tc.s2)
		if differs && !tc.expDiff {
			t.Logf(tc.IDStr())
			t.Logf("\t: %v\n", tc.s1)
			t.Logf("\t: %v\n", tc.s2)
			t.Errorf("\t: should not be reported as differing\n")
		} else if !differs && tc.expDiff {
			t.Logf(tc.IDStr())
			t.Logf("\t: %v\n", tc.s1)
			t.Logf("\t: %v\n", tc.s2)
			t.Errorf("\t: should be reported as differing\n")
		}
	}
}
