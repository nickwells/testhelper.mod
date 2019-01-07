package testhelper

import (
	"strings"
	"testing"
)

// PanicCheckString tests the panic value (which should be a string) against
// the passed values
func PanicCheckString(t *testing.T, name string,
	panicked, panicExpected bool,
	panicVal interface{}, expVal []string) bool {
	t.Helper()

	panicIsBad, msg := badPanicString(panicked, panicExpected, panicVal, expVal)
	if panicIsBad {
		t.Logf("test %s :\n", name)
		if panicked {
			t.Logf("\t: %v\n", panicVal)
		}
		t.Errorf("\t: %s", msg)
	}
	return panicIsBad
}

// badPanicString checks whether the panic is unexpected in some way and
// returns true and some explanatory message if so, false otherwise
func badPanicString(panicked, panicExpected bool,
	panicVal interface{}, expVal []string) (bool, string) {
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

	missing := missingParts(pvStr, expVal)
	if len(missing) > 0 {
		return true, "the panic message should contain: " +
			strings.Join(missing, "\n\t\t: and: ")
	}
	return false, ""
}
