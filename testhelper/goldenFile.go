package testhelper

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
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
//
//    DirNames is a slice of strings holding the parts of the directory path
//    to the file
//    Pfx is an optional prefix - leave it as an empty string to exclude it
//    Sfx is an optional suffix - as for the prefix
type GoldenFileCfg struct {
	DirNames []string
	Pfx      string
	Sfx      string
}

// PathName will return the name of a golden file. It applies the directory
// names and any prefix or suffix to the supplied string to give a well-formed
// name using the appropriate filepath separators for the operating system. A
// suggested name to pass to this method might be the name of the current
// test as given by the Name() method on testing.T.
//
// Note that any supplied name is "cleaned" by removing any part prior to an
// embedded filepath.Separator.
func (gfc GoldenFileCfg) PathName(name string) string {
	fNameParts := make([]string, 0, 3)
	if gfc.Pfx != "" {
		fNameParts = append(fNameParts, gfc.Pfx)
	}
	fNameParts = append(fNameParts, filepath.Base(name))
	if gfc.Sfx != "" {
		fNameParts = append(fNameParts, gfc.Sfx)
	}
	fName := strings.Join(fNameParts, ".")

	pathParts := make([]string, 0, len(gfc.DirNames)+1)
	pathParts = append(pathParts, gfc.DirNames...)
	pathParts = append(pathParts, fName)

	return filepath.Join(pathParts...)
}

// CheckAgainstGoldenFile checks that the value given matches the contents of
// the golden file and returns true if it does, false otherwise. It will
// report any errors it finds including any problems reading from or writing
// to the golden file itself. If the updGF flag is set to true then the
// golden file will be updated with the supplied value. You can set this
// value through a command-line parameter to the test and then pass that to
// this function as follows
//
//    var upd = flag.Bool("upd-gf", false, "update the golden files")
//    gfc := testhelper.GoldenFileCfg{
//        DirNames: []string{"testdata"},
//        Pfx:      "values",
//        Sfx:      "txt",
//    }
//    ...
//    testhelper.CheckAgainstGoldenFile(t, ID, val, gfc.PathName(t.Name()), *upd)
//
// Then to update the golden files you would invoke the test command as follows
//
//    go test -upd-gf
//
// Give the -v argument to go test to see what is being updated
func CheckAgainstGoldenFile(t *testing.T, testID string, val []byte, gfName string, updGF bool) bool {
	t.Helper()

	if updGF {
		if !updateGoldenFile(t, gfName, val) {
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

// updateGoldenFile will attempt to update the golden file with the new
// content and return true if it succeeds or false otherwise. If there is an
// existing golden file it will try to preverve the contents so that they can
// be compared with the new file. It reports its progress; if the file hasn't
// changed it does nothing.
func updateGoldenFile(t *testing.T, gfName string, val []byte) bool {
	t.Helper()

	nameLogged := false
	origVal, err := ioutil.ReadFile(gfName) // nolint: gosec
	if err == nil {
		if bytes.Equal(val, origVal) {
			return true
		}

		t.Log("Updating golden file:", gfName)
		nameLogged = true
		origFileName := gfName + ".orig"
		err = ioutil.WriteFile(origFileName, origVal, 0644)
		if err != nil {
			t.Log("\t: Couldn't preserve the original contents")
			t.Error("\t: ", err)
		} else {
			t.Log("\t: Original contents have been preserved")
			t.Log("\t: in:", origFileName)
			t.Log("\t: Compare the files to see changes and then remove it")
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Log("Updating golden file:", gfName)
		nameLogged = true
		t.Log("\t: Couldn't read the original contents")
		t.Error("\t: ", err)
	}

	if !nameLogged {
		t.Log("Creating a new golden file:", gfName)
	}

	err = ioutil.WriteFile(gfName, val, 0644)
	if err != nil {
		t.Log("\t: Couldn't write to the golden file")
		t.Error("\t: ", err)
		return false
	}

	return true
}
