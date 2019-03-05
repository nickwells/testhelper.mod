package testhelper

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

// GoldenFileCfg holds common configuration details for a collection of
// golden files. It helps with consistent naming of golden files without
// having to repeat common parts throughout the code.
//
// A golden file is a file that holds expected output (typically lengthy)
// that can be compared as part of a test. It avoids the need to have a long
// string in the body of a test.
type GoldenFileCfg struct {
	DirName string
	Pfx     string
	Sfx     string
}

// FileName will return the name of a golden file. It applies the directory
// name and any prefix or suffix to the supplied string to give a well-formed
// name using the appropriate filepath parts for the operating system. A
// suggested name to pass to this method might be the name of the current
// test as given by the Name() method on testing.T.
//
// Note that any supplied name is "cleaned" by removing any part prior to an
// embedded filepath.Separator.
func (gfc GoldenFileCfg) FileName(name string) string {
	return filepath.Join(gfc.DirName,
		strings.Join([]string{gfc.Pfx, filepath.Base(name), gfc.Sfx}, "."))
}

// CheckAgainstGoldenFile checks that the value given matches the contents of
// the golden file and returns true if it does, false otherwise. It will
// report any errors it finds including any problems reading from or writing
// to the golden file itself. If the updGF flag is set to true then the
// golden file will be updated with the supplied value. You can set this
// value through a command-line parameter to the test and then pass that to
// this function as follows
//
//    var upd = flag.Bool("upd", false, "update the golden files")
//    gfc := testhelper.GoldenFileCfg{
//        DirName: "testdata",
//        Pfx:     "values",
//        Sfx:     "txt",
//    }
//    ...
//    testhelper.CheckAgainstGoldenFile(t, ID, val, gfc.FileName(t.Name()), *upd)
func CheckAgainstGoldenFile(t *testing.T, testID string, val []byte, gfName string, updGF bool) bool {
	t.Helper()

	if updGF {
		err := ioutil.WriteFile(gfName, val, 0644)
		if err != nil {
			t.Errorf("Couldn't update the golden file: %s", err)
			return false
		}
	}

	expVal, err := ioutil.ReadFile(gfName) // nolint: gosec
	if err != nil {
		t.Errorf("couldn't read the expected value from the golden file: %s",
			err)
		return false
	}

	if !bytes.Equal(val, expVal) {
		t.Log(testID)
		t.Log("\t: Expected")
		t.Log(string(expVal))
		t.Log("\t: Actual")
		t.Log(string(val))
		t.Errorf("\t: Unexpected value")
		return false
	}
	return true
}
