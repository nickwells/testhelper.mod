package testhelper_test

import (
	"fmt"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func ExampleID_IDStr() {
	testCases := []struct {
		testhelper.ID
		v1, v2 int
	}{
		{
			ID: testhelper.MkID("my first example"),
			v1: 1,
			v2: 2,
		},
	}

	for _, tc := range testCases {
		if tc.v1 != tc.v2 {
			fmt.Println(tc.IDStr()) // in a real test this will be t.Error(...)
		}
	}
	// Output:
	// test: examples_test.go:15: my first example
}

func ExampleExpErr() {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		v int
	}{
		{
			ID: testhelper.MkID("No error expected"),
			v:  1,
		},
		{
			ID:     testhelper.MkID("error expected"),
			ExpErr: testhelper.MkExpErr("part1", "another part"),
			v:      99,
		},
	}

	for range testCases { // for _, tc := range testCases {
		//
		// we can then run the test, collect any error returned and check
		// that it is as expected. With the ExpErr member left at its default
		// value no error is expecteed. If it has a value the error is
		// expected and it should contain all the supplied strings.
		//
		// With the testcase structure having an ID and an ExpErr you can
		// call
		//
		//     testhelper.CheckExpErr(t, err, tc)
		//
		// just passing the testing.T pointer (t), the error (err) and the
		// testcase (tc). This will report any missing or unexpected errors
		// or errors that have unexpected values. It wil return false if
		// there are any problems
	}
}

func ExampleExpPanic() {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		v int
	}{
		{
			ID: testhelper.MkID("No panic expected"),
			v:  1,
		},
		{
			ID:       testhelper.MkID("panic expected"),
			ExpPanic: testhelper.MkExpPanic("part1", "another part"),
			v:        2,
		},
	}

	for range testCases { // for _, tc := range testCases {
		//
		// to test panics we will need to write a helper function taking the
		// testcase values that parameterise the test. This should call our
		// code to be tested and recover from any panics. This should then
		// return true if a panic was seen and an interface value containing
		// the panic value. We then collect these returned values and check
		// that they are as expected. With the ExpPanic member left at its
		// default value no panic is expecteed. If it has a value the panic
		// is expected and it should contain all the supplied strings.
		//
		// With the testcase structure having an ID and an ExpPanic you can
		// call
		//
		//     testhelper.CheckExpPanic(t, panicked, panicVal, tc)
		//
		// just passing the testing.T pointer (t), the boolean indicating
		// whether a panic was seen , the panic value and the testcase
		// (tc). This will report any missing or unexpected panics or panics
		// that have unexpected values. It wil return true if there are any
		// problems. There is an alternative function that can be called
		// which allows you to pass a stack trace as a final parameter; this
		// is useful if you get an unexpected panic and want to find out
		// where it came from
	}
}

// ExamplePanicSafe demonstrates how the PanicSafe function may be used
func ExamplePanicSafe() {
	panicked, panicVal := testhelper.PanicSafe(func() { panic("As expected") })
	fmt.Println("Panicked:", panicked)
	fmt.Println("PanicVal:", panicVal)
	// Output:
	// Panicked: true
	// PanicVal: As expected
}
