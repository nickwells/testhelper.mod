package testhelper

import "testing"

func TestBadPanicString(t *testing.T) {
	testCases := []struct {
		name                    string
		panicked, panicExpected bool
		panicVal                interface{}
		expVal                  []string
		expBad                  bool
		badMsgExpected          string
	}{
		{
			name: "no panic",
		},
		{
			name:           "no panic when expected",
			panicExpected:  true,
			expVal:         []string{"xxx"},
			expBad:         true,
			badMsgExpected: "a panic was expected but not seen",
		},
		{
			name:           "panic when not expected",
			panicked:       true,
			expVal:         []string{"xxx"},
			expBad:         true,
			badMsgExpected: "there was an unexpected panic",
		},
		{
			name:          "panic value not a string",
			panicked:      true,
			panicExpected: true,
			panicVal: struct {
				a int
				b string
			}{a: 1, b: "Hello, World!"},
			expVal:         []string{"xxx"},
			expBad:         true,
			badMsgExpected: "a panic was seen but the value was not a string",
		},
		{
			name:           "panic value does not contain the 1 expected value",
			panicked:       true,
			panicExpected:  true,
			panicVal:       "Hello, World!",
			expVal:         []string{"xxx"},
			expBad:         true,
			badMsgExpected: "the panic message should contain: xxx",
		},
		{
			name:           "panic value does not any of the expected values",
			panicked:       true,
			panicExpected:  true,
			panicVal:       "Hello, World!",
			expVal:         []string{"X", "Y"},
			expBad:         true,
			badMsgExpected: "the panic message should contain: X\n\t\t: and: Y",
		},
	}

	for i, tc := range testCases {
		panicIsBad, msg := badPanicString(tc.panicked, tc.panicExpected,
			tc.panicVal, tc.expVal)
		if panicIsBad && !tc.expBad {
			t.Logf("test %d: %s :\n", i, tc.name)
			t.Logf("\t: badPanic message: %s\n", msg)
			t.Errorf("\t: the panic was unexpectedly reported as bad\n")
		} else if !panicIsBad && tc.expBad {
			t.Logf("test %d: %s :\n", i, tc.name)
			t.Errorf(
				"\t: the panic was expected to be reported as bad but wasn't\n")
		} else if panicIsBad && tc.expBad {
			if tc.badMsgExpected != msg {
				t.Logf("test %d: %s :\n", i, tc.name)
				t.Logf("\t: badPanic message: %s\n", msg)
				t.Errorf("\t: the message was expected to be: %s\n",
					tc.badMsgExpected)
			}
		}
	}

}
