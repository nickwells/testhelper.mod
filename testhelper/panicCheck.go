package testhelper

import (
	"strings"
	"testing"
)

// PanicExp records common details about panic expectations for a test case
type PanicExp struct {
	Expected      bool
	ShouldContain []string
}

// MkPanicExp is a constructor for the PanicExp struct. The Expected flag is
// always set to true and the slice of strings that the panic should contain
// is set to the slice of strings passed. For a PanicExp where a panic is
// not expected just leave it in its default state.
func MkPanicExp(s ...string) PanicExp {
	return PanicExp{
		Expected:      true,
		ShouldContain: s,
	}
}

// Exp returns true or false according to the value of the Expected field
func (p PanicExp) Exp() bool {
	return p.Expected
}

// ShldCont returns the value of the ShouldContain field
func (p PanicExp) ShldCont() []string {
	return p.ShouldContain
}

// TestPanic is an interface wrapping the panic expectation methods
type TestPanic interface {
	Exp() bool
	ShldCont() []string
}

// TestCaseWithPanic combines the TestCase and TestPanic interfaces
type TestCaseWithPanic interface {
	TestCase
	TestPanic
}

// CheckPanic calls PanicCheckString using the details from the test case to
// supply the parameters
func CheckPanic(t *testing.T, panicked bool, panicVal interface{},
	tp TestCaseWithPanic) bool {
	t.Helper()
	return PanicCheckString(t, tp.IDStr(),
		panicked, tp.Exp(),
		panicVal, tp.ShldCont())
}

// CheckPanicWithStack calls PanicCheckStringWithStack using the details from
// the test case to supply the parameters
func CheckPanicWithStack(t *testing.T, panicked bool, panicVal interface{},
	tp TestCaseWithPanic, stackTrace []byte) bool {
	t.Helper()
	return PanicCheckStringWithStack(t, tp.IDStr(),
		panicked, tp.Exp(),
		panicVal, tp.ShldCont(),
		stackTrace)
}

// PanicCheckString tests the panic value (which should be a string) against
// the passed values
func PanicCheckString(t *testing.T, testID string,
	panicked, panicExpected bool,
	panicVal interface{}, shouldContain []string) bool {
	t.Helper()

	panicIsBad, msg :=
		badPanicString(panicked, panicExpected, panicVal, shouldContain)
	if panicIsBad {
		t.Log(testID)
		if panicked {
			t.Logf("\t: %v\n", panicVal)
		}
		t.Errorf("\t: %s", msg)
	}
	return panicIsBad
}

// PanicCheckStringWithStack tests the panic value (which should be a string)
// against the passed values. A stack trace should also be passed which will
// be printed if the panic is not as expected and a panic was seen
func PanicCheckStringWithStack(t *testing.T, testID string,
	panicked, panicExpected bool,
	panicVal interface{}, shouldContain []string, stackTrace []byte) bool {
	t.Helper()

	panicIsBad, msg :=
		badPanicString(panicked, panicExpected, panicVal, shouldContain)
	if panicIsBad {
		t.Log(testID)
		if panicked {
			t.Logf("\t: %v\n", panicVal)
			t.Log(string(stackTrace))
		}
		t.Errorf("\t: %s", msg)
	}
	return panicIsBad
}

// ReportUnexpectedPanic will check if panicked is true and will report the
// unexpected panic if true. it returns the panicked value
func ReportUnexpectedPanic(t *testing.T, testID string,
	panicked bool, panicVal interface{}, stackTrace []byte) bool {
	if panicked {
		t.Log(testID)
		t.Logf("\t: panic: %v\n", panicVal)
		t.Log("\t: At:", string(stackTrace))
		t.Error("\t: An unexpected panic was seen")
	}

	return panicked
}

// badPanicString checks whether the panic is unexpected in some way and
// returns true and some explanatory message if so, false otherwise
func badPanicString(panicked, panicExpected bool,
	panicVal interface{}, shouldContain []string) (bool, string) {
	if !panicked && !panicExpected {
		return false, ""
	}

	if panicked && !panicExpected {
		return true, "there was an unexpected panic"
	}

	if !panicked && panicExpected {
		return true, "a panic was expected but not seen"
	}

	pvStr, ok := panicVal.(string)
	if !ok {
		return true, "a panic was seen but the value was not a string"
	}

	missing := missingParts(pvStr, shouldContain)
	if len(missing) > 0 {
		return true, "the panic message should contain: " +
			strings.Join(missing, "\n\t\t: and: ")
	}
	return false, ""
}
