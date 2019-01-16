package testhelper

import (
	"strings"
	"testing"
)

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
