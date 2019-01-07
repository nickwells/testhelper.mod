package testhelper

import (
	"fmt"
	"testing"
)

func TestMissingParts(t *testing.T) {
	testCases := []struct {
		name    string
		s       string
		expect  []string
		missing []string
	}{
		{
			name: "no expectations",
			s:    "any string",
		},
		{
			name:   "all present",
			s:      "any string",
			expect: []string{"any", "string"},
		},
		{
			name:    "some missing",
			s:       "any string",
			expect:  []string{"hello", "world"},
			missing: []string{"hello", "world"},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s :\n", i, tc.name)
		missing := missingParts(tc.s, tc.expect)
		if StringSliceDiff(missing, tc.missing) {
			t.Log(tcID)
			t.Logf("\t: expected: %v", tc.missing)
			t.Logf("\t:      got: %v", missing)
			t.Errorf("\t: missingParts did not return the expected results\n")
		}
	}

}
