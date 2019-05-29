package testhelper

import (
	"testing"
)

func TestMissingParts(t *testing.T) {
	testCases := []struct {
		ID
		s       string
		expect  []string
		missing []string
	}{
		{
			ID: MkID("no expectations"),
			s:  "any string",
		},
		{
			ID:     MkID("all present"),
			s:      "any string",
			expect: []string{"any", "string"},
		},
		{
			ID:      MkID("some missing"),
			s:       "any string",
			expect:  []string{"hello", "world"},
			missing: []string{"hello", "world"},
		},
	}

	for _, tc := range testCases {
		missing := missingParts(tc.s, tc.expect)
		if StringSliceDiff(missing, tc.missing) {
			t.Log(tc.IDStr())
			t.Logf("\t: expected: %v", tc.missing)
			t.Logf("\t:      got: %v", missing)
			t.Errorf("\t: missingParts did not return the expected results\n")
		}
	}

}
