package testhelper_test

import (
	"github.com/nickwells/testhelper.mod/testhelper"
	"testing"
)

func TestStringSliceDiff(t *testing.T) {
	testCases := []struct {
		name         string
		s1, s2       []string
		diffExpected bool
	}{
		{
			name: "both nil",
		},
		{
			name:         "first nil, second not nil",
			s2:           []string{"a"},
			diffExpected: true,
		},
		{
			name:         "first not nil, second nil",
			s1:           []string{"a"},
			diffExpected: true,
		},
		{
			name:         "neither nil, first shorter than second",
			s1:           []string{"a"},
			s2:           []string{"a", "b"},
			diffExpected: true,
		},
		{
			name:         "neither nil, second shorter than first",
			s1:           []string{"a", "b"},
			s2:           []string{"a"},
			diffExpected: true,
		},
		{
			name:         "same length, differing",
			s1:           []string{"a", "b"},
			s2:           []string{"a", "c"},
			diffExpected: true,
		},
		{
			name:         "same length, same strings but reversed",
			s1:           []string{"a", "b"},
			s2:           []string{"b", "a"},
			diffExpected: true,
		},
		{
			name: "same length, same strings",
			s1:   []string{"a", "b"},
			s2:   []string{"a", "b"},
		},
	}

	for i, tc := range testCases {
		differs := testhelper.StringSliceDiff(tc.s1, tc.s2)
		if differs && !tc.diffExpected {
			t.Logf("test %d: %s :\n", i, tc.name)
			t.Logf("\t: %v\n", tc.s1)
			t.Logf("\t: %v\n", tc.s2)
			t.Errorf("\t: should not be reported as differing\n")
		} else if !differs && tc.diffExpected {
			t.Logf("test %d: %s :\n", i, tc.name)
			t.Logf("\t: %v\n", tc.s1)
			t.Logf("\t: %v\n", tc.s2)
			t.Errorf("\t: should be reported as differing\n")
		}
	}

}
