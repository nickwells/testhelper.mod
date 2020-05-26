package testhelper

import (
	"strings"
	"testing"
)

// ExpPanic records common details about panic expectations for a test
// case. It is intended that this should be embedded in a test case
// structure, which will also have an ID structure embedded. The resulting
// test case can then be passed to the CheckExpPanic func. It is similar to the
// ExpErr structure in form and intended use.
//
// Note that this expects any panic value to be a string and it will report
// an error if that is not the case.
type ExpPanic struct {
	Expected      bool
	ShouldContain []string
}

// MkExpPanic is a constructor for the ExpPanic struct. The Expected
// flag is always set to true and the slice of strings that the panic should
// contain is set to the slice of strings passed. For an ExpPanic where a
// panic is not expected just leave it in its default state.
func MkExpPanic(s ...string) ExpPanic {
	return ExpPanic{
		Expected:      true,
		ShouldContain: s,
	}
}

// PanicExpected returns true or false according to the value of the
// PanicExpected field
func (p ExpPanic) PanicExpected() bool {
	return p.Expected
}

// PanicShldCont returns the value of the ShouldContain field
func (p ExpPanic) PanicShldCont() []string {
	return p.ShouldContain
}

// TestPanic is an interface wrapping the panic expectation methods
type TestPanic interface {
	PanicExpected() bool
	PanicShldCont() []string
}

// TestCaseWithPanic combines the TestCase and TestPanic interfaces
type TestCaseWithPanic interface {
	TestCase
	TestPanic
}

// CheckExpPanic calls PanicCheckString using the details from the test case to
// supply the parameters
func CheckExpPanic(t *testing.T, panicked bool, panicVal interface{},
	tp TestCaseWithPanic) bool {
	t.Helper()

	return PanicCheckString(t, tp.IDStr(),
		panicked, tp.PanicExpected(),
		panicVal, tp.PanicShldCont())
}

// CheckExpPanicWithStack calls PanicCheckStringWithStack using the details from
// the test case to supply the parameters
func CheckExpPanicWithStack(t *testing.T, panicked bool, panicVal interface{},
	tp TestCaseWithPanic, stackTrace []byte) bool {
	t.Helper()

	return PanicCheckStringWithStack(t, tp.IDStr(),
		panicked, tp.PanicExpected(),
		panicVal, tp.PanicShldCont(),
		stackTrace)
}

// PanicCheckString tests the panic value (which should be a string) against
// the passed values. It will report an error if the panic status is
// unexpected.
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
	t.Helper()

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
