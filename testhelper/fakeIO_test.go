package testhelper_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestFakeIO(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		ignoreErr string
		input     string
		expStdout string
		expStderr string
		testFunc  func()
	}{
		{
			ID:       testhelper.MkID("no input, no output"),
			testFunc: func() {},
		},
		{
			ID:        testhelper.MkID("no input, write to stdout"),
			testFunc:  func() { fmt.Println("Hello, World!") },
			expStdout: "Hello, World!\n",
		},
		{
			ID:        testhelper.MkID("no input, write to stderr"),
			testFunc:  func() { fmt.Fprintln(os.Stderr, "Hello, World!") },
			expStderr: "Hello, World!\n",
		},
		{
			ID: testhelper.MkID("small input, write to stdout"),
			testFunc: func() {
				b := make([]byte, 5)
				_, _ = os.Stdin.Read(b)
				os.Stdout.Write(b)
			},
			input:     "Hello, World!\n",
			expStdout: "Hello",
		},
		{
			ID: testhelper.MkID("small input, write to stdout"),
			testFunc: func() {
				b := make([]byte, 5)
				_, _ = os.Stdin.Read(b)
				os.Stderr.Write(b)
			},
			input:     "Hello, World!\n",
			expStderr: "Hello",
		},
		{
			ID: testhelper.MkID("very large input, write to stdout"),
			testFunc: func() {
				b := make([]byte, 5)
				_, _ = os.Stdin.Read(b)
				os.Stdout.Write(b)
			},
			input:     strings.Repeat("Hello, World!\n", 1024*1024),
			expStdout: "Hello",
		},
		{
			ID: testhelper.MkID("very large input, write to stderr"),
			testFunc: func() {
				b := make([]byte, 5)
				_, _ = os.Stdin.Read(b)
				os.Stderr.Write(b)
			},
			input:     strings.Repeat("Hello, World!\n", 1024*1024),
			expStderr: "Hello",
		},
		{
			ID:        testhelper.MkID("very large input, close Stdin"),
			ignoreErr: "Error writing to stdin: broken pipe",
			testFunc: func() {
				b := make([]byte, 5)
				_, _ = os.Stdin.Read(b)
				os.Stdin.Close()
				os.Stderr.Write(b)
			},
			input:     strings.Repeat("Hello, World!\n", 1024*1024),
			expStderr: "Hello",
		},
	}

	for _, tc := range testCases {
		fio, err := testhelper.NewStdioFromString(tc.input)
		if err != nil {
			t.Fatal("Unexpected error (NewStdioFromString): ", err)
		}

		tc.testFunc()

		actOut, actErr, err := fio.Done()
		if fmt.Sprint(err) != tc.ignoreErr {
			testhelper.CheckExpErr(t, err, tc)
		}

		testhelper.DiffString(t, tc.IDStr(), "stdout",
			string(actOut), tc.expStdout)
		testhelper.DiffString(t, tc.IDStr(), "stderr",
			string(actErr), tc.expStderr)
	}
}
