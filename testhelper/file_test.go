package testhelper_test

import (
	"os"
	"testing"

	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestMakeTempDir(t *testing.T) {
	const tempDir = "testdata/temp.dir"
	const perms = os.FileMode(0555)

	var cleanup func()
	panicked, panicVal := testhelper.PanicSafe(func() {
		cleanup = testhelper.MakeTempDir(tempDir, perms)
	})
	if panicked {
		t.Fatalf("The MakeTempDir function panicked unexpectedly: %s",
			panicVal)
	}
	fi, err := os.Stat(tempDir)
	if err != nil {
		t.Errorf("Couldn't stat the temp directory %q: %s", tempDir, err)
	} else if !fi.IsDir() {
		t.Errorf("%q exists but is not a directory", tempDir)
	} else if fi.Mode()&os.ModePerm != perms {
		t.Logf("MakeTempDir: %s", tempDir)
		t.Logf("\t: expected permissions: %o", perms)
		t.Logf("\t:   actual permissions: %o", fi.Mode()&os.ModePerm)
		t.Errorf("\t: %q has the wrong permissions", tempDir)
	}
	panicked, panicVal = testhelper.PanicSafe(func() {
		testhelper.MakeTempDir(tempDir, perms)
	})
	if !panicked {
		t.Error(
			"The MakeTempDir function should have panicked - the dir exists")
	}
	err, ok := panicVal.(*os.PathError)
	if !ok {
		t.Errorf("Unexpected panic val when calling MakeTempDir twice: %s",
			panicVal)
	} else if !os.IsExist(err) {
		t.Errorf("Unexpected PathError when calling MakeTempDir twice: %s",
			err)
	}
	cleanup()
	_, err = os.Stat(tempDir)
	if err == nil {
		t.Errorf("The temp directory %q should not exist but does", tempDir)
	} else if !os.IsNotExist(err) {
		t.Errorf("Unexpected error from os.Stat(%q): %s", tempDir, err)
	}
}

func TestTempChmod(t *testing.T) {
	const fName = "testdata/file"
	const permMask = os.FileMode(0555)

	fi, err := os.Stat(fName)
	if err != nil {
		t.Fatalf("The file %q should exist but doesn't", fName)
	}

	originalPerms := fi.Mode() & os.ModePerm
	targetPerms := originalPerms ^ permMask

	var cleanup func()
	panicked, panicVal := testhelper.PanicSafe(func() {
		cleanup = testhelper.TempChmod(fName, targetPerms)
	})
	if panicked {
		t.Fatalf("The TempChmod function panicked unexpectedly: %s",
			panicVal)
	}
	fi, err = os.Stat(fName)
	if err != nil {
		t.Errorf("Couldn't stat the file %q: %s", fName, err)
	} else if fi.Mode()&os.ModePerm != targetPerms {
		t.Logf("TempChmod: %s", fName)
		t.Logf("\t: expected permissions: %o", targetPerms)
		t.Logf("\t:   actual permissions: %o", fi.Mode()&os.ModePerm)
		t.Errorf("\t: %q has the wrong permissions", fName)
	}
	cleanup()
	fi, err = os.Stat(fName)
	if err != nil {
		t.Errorf("The file %q should exist (post cleanup) but doesn't", fName)
	} else if fi.Mode()&os.ModePerm != originalPerms {
		t.Logf("TempChmod (post cleanup): %s", fName)
		t.Logf("\t: expected permissions: %o", originalPerms)
		t.Logf("\t:   actual permissions: %o", fi.Mode()&os.ModePerm)
		t.Errorf("\t: %q has the wrong permissions", fName)
	}
}

func TestTempChmodNoFile(t *testing.T) {
	const fName = "testdata/nosuchfile"

	panicked, panicVal := testhelper.PanicSafe(func() {
		testhelper.TempChmod(fName, os.ModePerm)
	})
	if !panicked {
		t.Fatalf("The TempChmod function was expected to panic but didn't")
	}
	err, ok := panicVal.(*os.PathError)
	if !ok {
		t.Errorf("Unexpected panic val when calling TempChmod"+
			" with a non-existent file: %s",
			panicVal)
	} else if !os.IsNotExist(err) {
		t.Errorf("Unexpected PathError when calling TempChmod"+
			" with a non-existent file: %s",
			err)
	}
}
