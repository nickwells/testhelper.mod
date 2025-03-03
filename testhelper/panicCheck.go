package testhelper

import (
	"fmt"
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

// PanicSafe will call the passed func, check if it has panicked and return
// true and the recovered panic value if it has and false and nil
// otherwise. You can call it passing a closure that does the work you are
// testing or any func matching the signature. It is intended to reduce the
// volume of boilerplate code you need.
func PanicSafe(f func()) (panicked bool, panicVal any) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()

	f()

	return panicked, panicVal
}

// CheckExpPanic calls PanicCheckString using the details from the test case to
// supply the parameters
func CheckExpPanic(
	t *testing.T, panicked bool, panicVal any,
	tp TestCaseWithPanic,
) bool {
	t.Helper()

	return PanicCheckString(t, tp.IDStr(),
		panicked, tp.PanicExpected(),
		panicVal, tp.PanicShldCont())
}

// CheckExpPanicWithStack calls PanicCheckStringWithStack using the details from
// the test case to supply the parameters
func CheckExpPanicWithStack(
	t *testing.T, panicked bool, panicVal any,
	tp TestCaseWithPanic, stackTrace []byte,
) bool {
	t.Helper()

	return PanicCheckStringWithStack(t, tp.IDStr(),
		panicked, tp.PanicExpected(),
		panicVal, tp.PanicShldCont(),
		stackTrace)
}

// CheckExpPanicError calls PanicCheckError using the details from the test
// case to supply the parameters
func CheckExpPanicError(t *testing.T, panicked bool, panicVal any,
	tp TestCaseWithPanic,
) bool {
	t.Helper()

	return PanicCheckError(t, tp.IDStr(),
		panicked, tp.PanicExpected(),
		panicVal, tp.PanicShldCont())
}

// CheckExpPanicErrorWithStack calls PanicCheckErrorWithStack using the
// details from the test case to supply the parameters
func CheckExpPanicErrorWithStack(t *testing.T,
	panicked bool, panicVal any,
	tp TestCaseWithPanic, stackTrace []byte,
) bool {
	t.Helper()

	return PanicCheckErrorWithStack(t, tp.IDStr(),
		panicked, tp.PanicExpected(),
		panicVal, tp.PanicShldCont(),
		stackTrace)
}

// PanicCheckString tests the panic value (which should be a string) against
// the passed values. It will report an error if the panic status is
// unexpected.
func PanicCheckString(t *testing.T, testID string,
	panicked, panicExpected bool,
	panicVal any, shouldContain []string,
) bool {
	t.Helper()

	msgs := badPanicString(panicked, panicExpected, panicVal, shouldContain)
	if len(msgs) > 0 {
		t.Log(testID)
		showPanicMsgs(t, panicked, panicVal, msgs)
	}

	return len(msgs) > 0
}

// PanicCheckStringWithStack tests the panic value (which should be a string)
// against the passed values. A stack trace should also be passed which will
// be printed if the panic is not as expected and a panic was seen
func PanicCheckStringWithStack(t *testing.T, testID string,
	panicked, panicExpected bool,
	panicVal any, shouldContain []string, stackTrace []byte,
) bool {
	t.Helper()

	msgs := badPanicString(panicked, panicExpected, panicVal, shouldContain)
	if len(msgs) > 0 {
		t.Log(testID)

		if panicked {
			t.Log(string(stackTrace))
		}

		showPanicMsgs(t, panicked, panicVal, msgs)
	}

	return len(msgs) > 0
}

// PanicCheckError tests the panic value (which should be an error) against
// the passed values. It will report an error if the panic status is
// unexpected.
func PanicCheckError(t *testing.T, testID string,
	panicked, panicExpected bool,
	panicVal any, shouldContain []string,
) bool {
	t.Helper()

	msgs := badPanicError(panicked, panicExpected, panicVal, shouldContain)
	if len(msgs) > 0 {
		t.Log(testID)
		showPanicMsgs(t, panicked, panicVal, msgs)
	}

	return len(msgs) > 0
}

// PanicCheckErrorWithStack tests the panic value (which should be an error)
// against the passed values. A stack trace should also be passed which will
// be printed if the panic is not as expected and a panic was seen
func PanicCheckErrorWithStack(t *testing.T, testID string,
	panicked, panicExpected bool,
	panicVal any, shouldContain []string, stackTrace []byte,
) bool {
	t.Helper()

	msgs := badPanicError(panicked, panicExpected, panicVal, shouldContain)
	if len(msgs) > 0 {
		t.Log(testID)

		if panicked {
			t.Log(string(stackTrace))
		}

		showPanicMsgs(t, panicked, panicVal, msgs)
	}

	return len(msgs) > 0
}

// ReportUnexpectedPanic will check if panicked is true and will report the
// unexpected panic if true. it returns the panicked value
func ReportUnexpectedPanic(t *testing.T, testID string,
	panicked bool, panicVal any, stackTrace []byte,
) bool {
	t.Helper()

	if panicked {
		t.Log(testID)
		t.Logf("\t: panic: %v", panicVal)
		t.Log("\t: At:", string(stackTrace))
		t.Error("\t: An unexpected panic was seen")
	}

	return panicked
}

// badPanicString checks whether the panic which should be a string is
// unexpected in some way and returns true and some explanatory message if
// so, false otherwise
func badPanicString(panicked, panicExpected bool,
	panicVal any, shouldContain []string,
) []string {
	if !(panicked && panicExpected) {
		return badPanic(panicked, panicExpected)
	}

	pvStr, ok := panicVal.(string)
	if !ok {
		return []string{
			fmt.Sprintf("a panic was seen but was not a string: %T", panicVal),
		}
	}

	return badPanicVal(pvStr, shouldContain)
}

// badPanicError checks whether the panic which should be a error is
// unexpected in some way and returns true and some explanatory message if
// so, false otherwise
func badPanicError(panicked, panicExpected bool,
	panicVal any, shouldContain []string,
) []string {
	if !(panicked && panicExpected) {
		return badPanic(panicked, panicExpected)
	}

	pvErr, ok := panicVal.(error)
	if !ok {
		return []string{
			fmt.Sprintf("a panic was seen but was not an error: %T", panicVal),
		}
	}

	pvStr := pvErr.Error()

	return badPanicVal(pvStr, shouldContain)
}

// badPanicVal checks the panic value
func badPanicVal(act string, exp []string) []string {
	missing := missingParts(act, exp)
	if len(missing) > 0 {
		rval := []string{"the panic message should contain:"}
		return append(rval, missing...)
	}

	return nil
}

// badPanic checks the flags
func badPanic(panicked, panicExpected bool) []string {
	if !panicked && !panicExpected {
		return nil
	}

	if panicked && !panicExpected {
		return []string{"there was an unexpected panic"}
	}

	if !panicked && panicExpected {
		return []string{"a panic was expected but not seen"}
	}

	return []string{"badPanic has been called unexpectedly"}
}

// showPanicMsgs reports the problems found with the panic
func showPanicMsgs(t *testing.T, panicked bool, pv any, msgs []string) {
	t.Helper()

	if len(msgs) > 0 {
		if panicked {
			t.Log("\t: Panic value:")
			t.Logf("\t\t%v", pv)
		}

		intro := "\t: "

		for _, msg := range msgs {
			t.Log(intro + msg)
			intro = "\t\t"
		}

		t.Error("\t: Bad Panic")
	}
}
