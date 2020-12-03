package testhelper

import (
	"strings"
	"testing"
)

// ShouldContain checks that the string 'act' contains all of the strings in
// the 'exp' argument and reports an error if it does not. The desc
// parameter is used to describe the string being checked. It returns true if
// a problem was found, false otherwise.
func ShouldContain(t *testing.T, testID, desc, act string, exp []string) bool {
	t.Helper()

	missing := missingParts(act, exp)
	if len(missing) > 0 {
		t.Log(testID)
		t.Logf("\t: an unexpected %s value was seen:", desc)
		t.Log("\t\t" + act)
		t.Log("\t: it should contain:")

		for _, part := range missing {
			t.Log("\t\t" + part)
		}
		t.Error("\t: Parts of the string were missing")
		return true
	}
	return false
}

// missingParts returns the entries in exp which are not in act
func missingParts(act string, exp []string) []string {
	missing := []string{}
	for _, part := range exp {
		if !strings.Contains(act, part) {
			missing = append(missing, part)
		}
	}
	return missing
}
