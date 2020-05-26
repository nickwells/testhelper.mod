package testhelper

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

const (
	pBits    = 0644
	dirPBits = 0755
)

// AddUpdateFlag adds a new flag to the standard flag package. The flag is
// used to control whether or not to update the Golden files with the new
// values rather than reporting differences as test errors. If there is
// already a Golden file present then this will be preserved in a file with
// the same name as the Golden file but with ".orig" as a suffix. This can
// then be set on the command line when testing and looked up by the
// GoldenFileCfg.Check method. The Check method will report the flag name to
// use if any is available.
func (gfc *GoldenFileCfg) AddUpdateFlag() {
	if gfc.updFlagAdded {
		return
	}

	gfGlob := gfc.PathName("*")
	if gfc.UpdFlagName == "" {
		panic(errors.New(
			"AddUpdateFlag has been called for files in " + gfGlob +
				" but the GoldenFileCfg has no flag name set"))
	}

	flag.BoolVar(&gfc.updFlag, gfc.UpdFlagName, false,
		"set this flag to update the golden files in "+gfGlob)

	gfc.updFlagAdded = true
}

// AddKeepBadResultsFlag adds a new flag to the standard flag package. The
// flag is used to control whether or not to keep the bad results in a
// file. The name of the file will be the name of the Golden file with
// ".badResults" as a suffix. These files can then be compared with the
// Golden files do see what changes have been made . This can then be set on
// the command line when testing and looked up by the GoldenFileCfg.Check
// method. The Check method will report the flag name to use if any is
// available.
func (gfc *GoldenFileCfg) AddKeepBadResultsFlag() {
	if gfc.keepBadResultsFlagAdded {
		return
	}

	gfGlob := gfc.PathName("*")
	if gfc.KeepBadResultsFlagName == "" {
		panic(errors.New(
			"AddKeepBadResultsFlag has been called for files in " + gfGlob +
				" but the GoldenFileCfg has no flag name set"))
	}

	flag.BoolVar(&gfc.keepBadResultsFlag, gfc.KeepBadResultsFlagName, false,
		"set this flag to keep bad results in"+gfGlob)

	gfc.keepBadResultsFlagAdded = true
}

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
//
//    Pfx is an optional prefix - leave it as an empty string to exclude it
//
//    Sfx is an optional suffix - as for the prefix
//
//    UpdFlagName is the name of a flag that will set a bool used to decide
//    whether or not to update the golden file. If it is not set then it is
//    ignored. If you have set this then you should also call the AddUpdateFlag
//    method (typically in an init() function) and then use the Check method
//    to compare with the file
//
//    KeepBadResultsFlagName is the name of a flag that will set a bool used
//    to decide whether or not to keep bad results. If it is not set then it
//    is ignored. If you have set this then you should also call the
//    AddKeepBadResultsFlag method (typically in an init() function) and then
//    use the Check method to compare with the file
type GoldenFileCfg struct {
	DirNames []string
	Pfx      string
	Sfx      string

	UpdFlagName  string
	updFlag      bool
	updFlagAdded bool

	KeepBadResultsFlagName  string
	keepBadResultsFlag      bool
	keepBadResultsFlagAdded bool
}

// Check confirms that the value given matches the contents of the golden
// file and returns true if it does, false otherwise. It will report any
// errors it finds including any problems reading from or writing to the
// golden file itself.
//
// If UpdFlagName is not empty and the AddUpdateFlag method
// has been called (typically in an init() function) then the corresponding
// flag value will be looked up and if the flag is set to true the golden
// file will be updated with the supplied value. You can set this value
// through a command-line parameter to the test and then pass that to this
// function as follows:
//
//    gfc := testhelper.GoldenFileCfg{
//        DirNames:    []string{"testdata"},
//        Pfx:         "values",
//        Sfx:         "txt",
//        UpdFlagName: "upd-gf",
//    }
//
//    func init() {
//        gfc.AddUpdateFlag()
//    }
//    ...
//    gfc.Check(t, "my value test", t.Name(), val)
//
// Then to update the golden files you would invoke the test command as follows
//
//    go test -upd-gf
//
// Similarly with the KeepBadResultsFlag
//
// Give the -v argument to go test to see what is being updated.
//
// An advantage of using this method (over using the
// testhelper.CheckAgainstGoldenFile function) is that this will show the
// name of the flag to use in order to update the files. You save the hassle
// of scanning the code to find out what you called the flag.
func (gfc GoldenFileCfg) Check(t *testing.T, id, gfName string, val []byte) bool {
	t.Helper()

	if gfc.UpdFlagName != "" && !gfc.updFlagAdded {
		panic(fmt.Errorf(
			"the name of the flag to update the golden files has been"+
				" given (%q) but the flag has not been added."+
				" You should call the AddUpdateFlag() method"+
				" (typically in an init() function)",
			gfc.UpdFlagName))
	}

	if gfc.KeepBadResultsFlagName != "" && !gfc.keepBadResultsFlagAdded {
		panic(fmt.Errorf(
			"the name of the flag to keep bad results has been"+
				" given (%q) but the flag has not been added."+
				" You should call the AddKeepBadResultsFlag() method"+
				" (typically in an init() function)",
			gfc.KeepBadResultsFlagName))
	}

	return gfc.checkFile(t, id, gfc.PathName(gfName), val)
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

// CheckAgainstGoldenFile confirms that the value given matches the contents
// of the golden file and returns true if it does, false otherwise. It will
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
//    testhelper.CheckAgainstGoldenFile(t,
//        "my value test",
//        val,
//        gfc.PathName(t.Name()),
//        *upd)
//
// Then to update the golden files you would invoke the test command as follows
//
//    go test -upd-gf
//
// Give the -v argument to go test to see what is being updated
//
// Deprecated: use the Check method on the GoldenFileCfg
func CheckAgainstGoldenFile(t *testing.T, testID string, val []byte, gfName string, updGF bool) bool {
	t.Helper()

	return checkFile(t, testID, gfName, val, updGF)
}

// getExpVal reads the contents of the golden file. If the updGF flag is set
// then if will write the contents of the file before reading it. It returns
// the contents and true if all went well, nil and false otherwise. It will
// report any errors it finds including any problems reading from or writing
// to the golden file itself.
func getExpVal(t *testing.T, id, gfName string, val []byte, updGF bool) ([]byte, bool) {
	t.Helper()

	if updGF {
		if !updateGoldenFile(t, gfName, val) {
			return nil, false
		}
	}

	expVal, err := ioutil.ReadFile(gfName) // nolint: gosec
	if err != nil {
		t.Log(id)
		t.Logf("\t: Problem with the golden file: %q", gfName)
		t.Errorf("\t: Couldn't read the expected value. Error: %s", err)
		return nil, false
	}
	return expVal, true
}

// checkFile confirms that the value given matches the contents of the golden
// file and returns true if it does, false otherwise. It will report any
// errors it finds including any problems reading from or writing to the
// golden file itself. If the updGF flag is set to true then the golden file
// will be updated with the supplied value.
func checkFile(t *testing.T, id, gfName string, val []byte, updGF bool) bool {
	t.Helper()

	expVal, ok := getExpVal(t, id, gfName, val, updGF)
	if !ok {
		t.Errorf("\t: Actual\n" + string(val))
		return false
	}

	return actEqualsExp(t, id, gfName, val, expVal)
}

// checkFile confirms that the value given matches the contents of the golden
// file and returns true if it does, false otherwise. It will report any
// errors it finds including any problems reading from or writing to the
// golden file itself. If the updGF flag is set to true then the golden file
// will be updated with the supplied value.
func (gfc GoldenFileCfg) checkFile(t *testing.T, id, gfName string, val []byte) bool {
	t.Helper()

	expVal, ok := getExpVal(t, id, gfName, val, gfc.updFlag)
	if !ok {
		if gfc.UpdFlagName != "" {
			t.Errorf("\t: To update the golden file with the new value"+
				" pass %q to the go test command", "-"+gfc.UpdFlagName)
		}
		t.Errorf("\t: Actual\n" + string(val))
		return false
	}

	if actEqualsExp(t, id, gfName, val, expVal) {
		return true
	}

	if gfc.UpdFlagName != "" {
		t.Errorf("\t: To update the golden file with the new value"+
			" pass %q to the go test command", "-"+gfc.UpdFlagName)
	}
	if gfc.keepBadResultsFlag {
		keepBadResults(t, gfName, val)
	} else if gfc.KeepBadResultsFlagName != "" {
		t.Errorf("\t: To keep the (bad) Actual results for later"+
			" investigation pass %q to the go test command",
			"-"+gfc.KeepBadResultsFlagName)
	}
	return false
}

// actEqualsExp compares the expected value against the actual and reports any
// difference. It will return true if they are equal and false otherwise
func actEqualsExp(t *testing.T, id, gfName string, actVal, expVal []byte) bool {
	t.Helper()

	if bytes.Equal(actVal, expVal) {
		return true
	}

	t.Log(id)
	t.Log("\t: Expected\n" + string(expVal))
	t.Log("\t: Actual\n" + string(actVal))
	t.Errorf("\t: The value given differs from the golden file value: %q",
		gfName)
	return false
}

// updateGoldenFile will attempt to update the golden file with the new
// content and return true if it succeeds or false otherwise. If there is an
// existing golden file it will try to preverve the contents so that they can
// be compared with the new file. It reports its progress; if the file hasn't
// changed it does nothing.
func updateGoldenFile(t *testing.T, gfName string, val []byte) bool {
	t.Helper()

	origVal, err := ioutil.ReadFile(gfName) // nolint: gosec
	if err == nil {
		if bytes.Equal(val, origVal) {
			return true
		}

		origFileName := gfName + ".orig"
		writeFile(t, origFileName, "original contents", origVal)
	} else if !os.IsNotExist(err) {
		t.Log("Couldn't preserve the original contents")
		t.Logf("\t: Couldn't read the golden file: %q", gfName)
		t.Error("\t: ", err)
	}

	if !writeFile(t, gfName, "golden", val) {
		return false
	}

	return true
}

// keepBadResults will attempt to write the bad results to a new file.
func keepBadResults(t *testing.T, gfName string, val []byte) {
	t.Helper()

	fName := gfName + ".badResults"
	writeFile(t, fName, "bad results", val)
}

// writeFile will write the values into the file. If the parent directories
// do not exist then it will create them and try again.
func writeFile(t *testing.T, fName, desc string, val []byte) (rval bool) {
	t.Helper()

	rval = true
	var err error
	defer func() {
		if err != nil {
			t.Logf("\t: Couldn't write to the %s file", desc)
			t.Error("\t: ", err)
			rval = false
		}
	}()

	t.Logf("Updating/Creating the %s file: %q", desc, fName)
	err = ioutil.WriteFile(fName, val, pBits)
	if os.IsNotExist(err) {
		dir := path.Dir(fName)
		if dir == "." {
			return
		}
		err = os.MkdirAll(dir, dirPBits)
		if err != nil {
			return
		}

		err = ioutil.WriteFile(fName, val, pBits)
	}
	return
}
