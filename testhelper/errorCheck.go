package testhelper

import (
	"testing"
)

// ErrInfo records common details about error expectations for a test case
type ErrInfo struct {
	Expected      bool
	ShouldContain []string
}

// Exp returns true or false according to the value of the Expected field
func (e ErrInfo) Exp() bool {
	return e.Expected
}

// ShldCont returns the value of the ShouldContain field
func (e ErrInfo) ShldCont() []string {
	return e.ShouldContain
}

// Err is an interface wrapping the error expectation methods
type Err interface {
	Exp() bool
	ShldCont() []string
}

// TestCaseWithErr combines the TestCase and Err interfaces
type TestCaseWithErr interface {
	TestCase
	Err
}

// ErrCheck will call CheckError using the details from the test case to
// supply the parameters
func ErrCheck(t *testing.T, i int, err error, tc TestCaseWithErr) bool {
	t.Helper()
	return CheckError(t, tc.MakeID(i), err, tc.Exp(), tc.ShldCont())
}

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
		}

		return !ShouldContain(t, testID, "error", err.Error(), shouldContain)
	}

	if expected {
		t.Log(testID)
		t.Errorf("\t: an error was expected but none was returned\n")
		return false
	}

	return true
}
