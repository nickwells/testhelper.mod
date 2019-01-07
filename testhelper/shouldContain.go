package testhelper

import (
	"strings"
	"testing"
)

// ShouldContain checks that the string s contains all of the strings in
// expect and reports an error if it does not. The desc parameter is used to
// describe the string being checked. It returns true if a problem was
// found, false otherwise.
func ShouldContain(t *testing.T, testID, desc, s string, expect []string) bool {
	t.Helper()

	missing := missingParts(s, expect)
	if len(missing) > 0 {
		t.Log(testID)
		t.Errorf("\t: an unexpected %s value was seen: %s\n", desc, s)
		t.Log("\t: it should contain:\n")

		for _, part := range missing {
			t.Log("\t\t", part, "\n")
		}
		return true
	}
	return false
}

// missingParts returns the entries in expect which are not in s
func missingParts(s string, expect []string) []string {
	missing := []string{}
	for _, part := range expect {
		if !strings.Contains(s, part) {
			missing = append(missing, part)
		}
	}
	return missing
}
