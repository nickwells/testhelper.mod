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
		input    string
		expOut   string
		expErr   string
		testFunc func()
	}{
		{
			ID:       testhelper.MkID("no input, no output"),
			testFunc: func() {},
		},
		{
			ID:       testhelper.MkID("no input, write to stdout"),
			testFunc: func() { fmt.Println("Hello, World!") },
			expOut:   "Hello, World!\n",
		},
		{
			ID:       testhelper.MkID("no input, write to stderr"),
			testFunc: func() { fmt.Fprintln(os.Stderr, "Hello, World!") },
			expErr:   "Hello, World!\n",
		},
		{
			ID: testhelper.MkID("small input, write to stdout"),
			testFunc: func() {
				b := make([]byte, 5)
				os.Stdin.Read(b)
				os.Stdout.Write(b)
			},
			input:  "Hello, World!\n",
			expOut: "Hello",
		},
		{
			ID: testhelper.MkID("small input, write to stdout"),
			testFunc: func() {
				b := make([]byte, 5)
				os.Stdin.Read(b)
				os.Stderr.Write(b)
			},
			input:  "Hello, World!\n",
			expErr: "Hello",
		},
		{
			ID: testhelper.MkID("very large input, write to stdout"),
			testFunc: func() {
				b := make([]byte, 5)
				os.Stdin.Read(b)
				os.Stdout.Write(b)
			},
			input:  strings.Repeat("Hello, World!\n", 1024*1024),
			expOut: "Hello",
		},
		{
			ID: testhelper.MkID("very large input, write to stdout"),
			testFunc: func() {
				b := make([]byte, 5)
				os.Stdin.Read(b)
				os.Stderr.Write(b)
			},
			input:  strings.Repeat("Hello, World!\n", 1024*1024),
			expErr: "Hello",
		},
		{
			ID: testhelper.MkID("very large input, close Stdin"),
			testFunc: func() {
				b := make([]byte, 5)
				os.Stdin.Read(b)
				os.Stdin.Close()
				os.Stderr.Write(b)
			},
			input:  strings.Repeat("Hello, World!\n", 1024*1024),
			expErr: "Hello",
		},
	}

	for _, tc := range testCases {
		fio, err := testhelper.NewStdioFromString(tc.input)
		if err != nil {
			t.Fatal("Unexpected error: ", err)
		}
		tc.testFunc()
		actOut, actErr, err := fio.Done()
		testhelper.DiffString(t, tc.IDStr(), "out", string(actOut), tc.expOut)
		testhelper.DiffString(t, tc.IDStr(), "err", string(actErr), tc.expErr)

	}
}
