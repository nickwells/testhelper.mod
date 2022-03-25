package testhelper

import (
	"errors"
	"testing"
)

func TestBadPanicString(t *testing.T) {
	testCases := []struct {
		ID
		ExpPanic
		panicked bool
		panicVal any
		expBad   bool
		expMsgs  []string
	}{
		{
			ID: MkID("no panic"),
		},
		{
			ID:       MkID("no panic when expected"),
			ExpPanic: MkExpPanic("xxx"),
			expBad:   true,
			expMsgs:  []string{"a panic was expected but not seen"},
		},
		{
			ID:       MkID("panic when not expected"),
			panicked: true,
			expBad:   true,
			expMsgs:  []string{"there was an unexpected panic"},
		},
		{
			ID:       MkID("panic value not a string"),
			panicked: true,
			panicVal: 1,
			ExpPanic: MkExpPanic("xxx"),
			expBad:   true,
			expMsgs:  []string{"a panic was seen but was not a string: int"},
		},
		{
			ID: MkID(
				"panic value does not contain the expected value"),
			panicked: true,
			panicVal: "Hello, World!",
			ExpPanic: MkExpPanic("xxx"),
			expBad:   true,
			expMsgs:  []string{"the panic message should contain:", "xxx"},
		},
		{
			ID: MkID(
				"panic value does not contain any of the expected values"),
			panicked: true,
			panicVal: "Hello, World!",
			ExpPanic: MkExpPanic("X", "Y"),
			expBad:   true,
			expMsgs:  []string{"the panic message should contain:", "X", "Y"},
		},
		{
			ID:       MkID("panic value contains the expected values"),
			panicked: true,
			panicVal: "Hello, World!",
			ExpPanic: MkExpPanic("Hello", "World"),
		},
	}

	for _, tc := range testCases {
		var panicIsBad bool
		msgs := badPanicString(tc.panicked, tc.ExpPanic.Expected,
			tc.panicVal, tc.ExpPanic.ShouldContain)
		if len(msgs) > 0 {
			panicIsBad = true
		}
		if panicIsBad && !tc.expBad {
			t.Logf(tc.IDStr())
			t.Logf("\t: badPanic message:\n")
			for _, msg := range msgs {
				t.Log("\t\t", msg, "\n")
			}
			t.Errorf("\t: the panic was unexpectedly reported as bad\n")
		} else if !panicIsBad && tc.expBad {
			t.Logf(tc.IDStr())
			t.Errorf(
				"\t: the panic was expected to be reported as bad but wasn't\n")
		} else if panicIsBad && tc.expBad {
			DiffStringSlice(t, tc.IDStr(), "panic value", msgs, tc.expMsgs)
		}
	}
}

func TestBadPanicError(t *testing.T) {
	testCases := []struct {
		ID
		ExpPanic
		panicked bool
		panicVal any
		expBad   bool
		expMsgs  []string
	}{
		{
			ID: MkID("no panic"),
		},
		{
			ID:       MkID("no panic when expected"),
			ExpPanic: MkExpPanic("xxx"),
			expBad:   true,
			expMsgs:  []string{"a panic was expected but not seen"},
		},
		{
			ID:       MkID("panic when not expected"),
			panicked: true,
			expBad:   true,
			expMsgs:  []string{"there was an unexpected panic"},
		},
		{
			ID:       MkID("panic value not an error"),
			panicked: true,
			panicVal: 1,
			ExpPanic: MkExpPanic("xxx"),
			expBad:   true,
			expMsgs:  []string{"a panic was seen but was not an error: int"},
		},
		{
			ID: MkID(
				"panic value does not contain the expected value"),
			panicked: true,
			panicVal: errors.New("hello, world"),
			ExpPanic: MkExpPanic("xxx"),
			expBad:   true,
			expMsgs:  []string{"the panic message should contain:", "xxx"},
		},
		{
			ID: MkID(
				"panic value does not contain any of the expected values"),
			panicked: true,
			panicVal: errors.New("hello, world"),
			ExpPanic: MkExpPanic("X", "Y"),
			expBad:   true,
			expMsgs:  []string{"the panic message should contain:", "X", "Y"},
		},
		{
			ID:       MkID("panic value contains the expected values"),
			panicked: true,
			panicVal: errors.New("hello, world"),
			ExpPanic: MkExpPanic("hello", "world"),
		},
	}

	for _, tc := range testCases {
		var panicIsBad bool
		msgs := badPanicError(tc.panicked, tc.ExpPanic.Expected,
			tc.panicVal, tc.ExpPanic.ShouldContain)
		if len(msgs) > 0 {
			panicIsBad = true
		}
		if panicIsBad && !tc.expBad {
			t.Logf(tc.IDStr())
			t.Logf("\t: badPanic message:\n")
			for _, msg := range msgs {
				t.Log("\t\t", msg, "\n")
			}
			t.Errorf("\t: the panic was unexpectedly reported as bad\n")
		} else if !panicIsBad && tc.expBad {
			t.Logf(tc.IDStr())
			t.Errorf(
				"\t: the panic was expected to be reported as bad but wasn't\n")
		} else if panicIsBad && tc.expBad {
			DiffStringSlice(t, tc.IDStr(), "panic value", msgs, tc.expMsgs)
		}
	}
}
