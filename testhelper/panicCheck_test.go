package testhelper

import "testing"

func TestBadPanicString(t *testing.T) {
	testCases := []struct {
		ID
		ExpPanic
		panicked       bool
		panicVal       interface{}
		expBad         bool
		badMsgExpected string
	}{
		{
			ID: MkID("no panic"),
		},
		{
			ID:             MkID("no panic when expected"),
			ExpPanic:       MkExpPanic("xxx"),
			expBad:         true,
			badMsgExpected: "a panic was expected but not seen",
		},
		{
			ID:             MkID("panic when not expected"),
			panicked:       true,
			expBad:         true,
			badMsgExpected: "there was an unexpected panic",
		},
		{
			ID:             MkID("panic value not a string"),
			panicked:       true,
			panicVal:       1,
			ExpPanic:       MkExpPanic("xxx"),
			expBad:         true,
			badMsgExpected: "a panic was seen but the value was not a string",
		},
		{
			ID: MkID(
				"panic value does not contain the expected value"),
			panicked:       true,
			panicVal:       "Hello, World!",
			ExpPanic:       MkExpPanic("xxx"),
			expBad:         true,
			badMsgExpected: "the panic message should contain: xxx",
		},
		{
			ID: MkID(
				"panic value does not contain any of the expected values"),
			panicked:       true,
			panicVal:       "Hello, World!",
			ExpPanic:       MkExpPanic("X", "Y"),
			expBad:         true,
			badMsgExpected: "the panic message should contain: X\n\t\t: and: Y",
		},
	}

	for _, tc := range testCases {
		panicIsBad, msg := badPanicString(tc.panicked, tc.ExpPanic.Expected,
			tc.panicVal, tc.ExpPanic.ShouldContain)
		if panicIsBad && !tc.expBad {
			t.Logf(tc.IDStr())
			t.Logf("\t: badPanic message: %s\n", msg)
			t.Errorf("\t: the panic was unexpectedly reported as bad\n")
		} else if !panicIsBad && tc.expBad {
			t.Logf(tc.IDStr())
			t.Errorf(
				"\t: the panic was expected to be reported as bad but wasn't\n")
		} else if panicIsBad && tc.expBad {
			if tc.badMsgExpected != msg {
				t.Logf(tc.IDStr())
				t.Logf("\t: badPanic message: %s\n", msg)
				t.Errorf("\t: the message was expected to be: %s\n",
					tc.badMsgExpected)
			}
		}
	}

}
