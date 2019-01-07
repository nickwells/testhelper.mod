package testhelper

import (
	"strings"
	"testing"
)

// CheckError checks that the error is nil if it is not expected, that it is
// non-nil if it is expected and that it contains the expected content if it
// is expected and non-nil. It will return false if there is any problem with
// the error, true otherwise
func CheckError(t *testing.T, testID string, err error, expected bool, mustContain []string) bool {
	t.Helper()

	if err != nil {
		if !expected {
			t.Logf("%s\n", testID)
			t.Errorf("\t: there was an unexpected err: %s\n", err)
			return false
		} else {
			emsg := err.Error()
			errorFound := false
			for _, s := range mustContain {
				if !strings.Contains(emsg, s) {
					if !errorFound {
						t.Logf("%s\n", testID)
						t.Logf("\t: Error: '%s'\n", emsg)
						errorFound = true
					}
					t.Errorf("\t: the error did not contain: '%s'\n", s)
				}
			}
			if ShouldContain(t, testID, "error", err.Error(), mustContain) {
				return false
			}
		}
	} else if err == nil && expected {
		t.Logf("%s\n", testID)
		t.Errorf("\t: an error was expected but none was returned\n")
		return false
	}

	return true
}
