package testhelper

import (
	"testing"
)

// CheckError checks that the error is nil if it is not expected, that it is
// non-nil if it is expected and that it contains the expected content if it
// is expected and non-nil. It will return false if there is any problem with
// the error, true otherwise
func CheckError(t *testing.T, testID string, err error, expected bool, shouldContain []string) bool {
	t.Helper()

	if err != nil {
		if !expected {
			t.Log(testID)
			t.Errorf("\t: there was an unexpected err: %s\n", err)
			return false
		} else {
			return !ShouldContain(t, testID, "error", err.Error(), shouldContain)
		}
	} else if err == nil && expected {
		t.Log(testID)
		t.Errorf("\t: an error was expected but none was returned\n")
		return false
	}

	return true
}
