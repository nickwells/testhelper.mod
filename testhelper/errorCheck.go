package testhelper

import (
	"testing"
)

// ErrExp records common details about error expectations for a test case
type ErrExp struct {
	Expected      bool
	ShouldContain []string
}

// MkErrExp is a constructor for the ErrExp struct. The Expected flag is
// always set to true and the slice of strings that the error should contain
// is set to the slice of strings passed. For an ErrExp where the error is
// not expected just leave it in its default state.
func MkErrExp(s ...string) ErrExp {
	return ErrExp{
		Expected:      true,
		ShouldContain: s,
	}
}

// Exp returns true or false according to the value of the Expected field
func (e ErrExp) Exp() bool {
	return e.Expected
}

// ShldCont returns the value of the ShouldContain field
func (e ErrExp) ShldCont() []string {
	return e.ShouldContain
}

// TestErr is an interface wrapping the error expectation methods
type TestErr interface {
	Exp() bool
	ShldCont() []string
}

// TestCaseWithErr combines the TestCase and TestErr interfaces
type TestCaseWithErr interface {
	TestCase
	TestErr
}

// ErrCheck calls CheckError using the details from the test case to supply
// the parameters
func ErrCheck(t *testing.T, err error, tce TestCaseWithErr) bool {
	t.Helper()
	return CheckError(t, tce.IDStr(), err, tce.Exp(), tce.ShldCont())
}

// ErrCheckWithID calls CheckError using the details from the TestErr to
// supply the parameters. The testID is supplied separately
func ErrCheckWithID(t *testing.T, testID string, err error, te TestErr) bool {
	t.Helper()
	return CheckError(t, testID, err, te.Exp(), te.ShldCont())
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
