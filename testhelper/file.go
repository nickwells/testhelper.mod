package testhelper

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// MakeTempDir creates a directory with the given name and permissions and
// returns a function that will delete the directory. This can be used as a
// defer func to tidy up after the tests are complete. This might be useful
// where you want to create an empty or unreadable directory for a test.
//
// It will panic if the directory cannot be created or the permissions set.
//
// Only the Permission bits of the perms are used, any other values are
// masked out before use.
func MakeTempDir(name string, perms os.FileMode) func() {
	err := os.Mkdir(name, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.Chmod(name, perms&os.ModePerm)
	if err != nil {
		_ = os.Remove(name)
		panic(err)
	}

	return func() { _ = os.Remove(name) }
}

// TempChmod will change the FileMode of the named file and return a function
// that will restore the original perms.
//
// It will panic if the permission bits can't be set or the original values
// obtained.
//
// Only the Permission bits of the perms are used, any other values are
// masked out before use.
func TempChmod(name string, perms os.FileMode) func() {
	fi, err := os.Stat(name)
	if err != nil {
		panic(err)
	}

	err = os.Chmod(name, perms&os.ModePerm)
	if err != nil {
		panic(err)
	}

	return func() { _ = os.Chmod(name, fi.Mode()&os.ModePerm) }
}

// MakeTempDirCopy creates a temporary directory and copies the contents of
// fromDir into it. It returns the name of the temporary directory, a
// function to clean up (remove the temp dir) and any error encountered.
// Only regular files and sub-directories can be copied. Any other
// types of file in the from-directory will cause an error.
func MakeTempDirCopy(fromDir string) (string, func() error, error) {
	tmpDir, err := os.MkdirTemp("", "testdir.")
	if err != nil {
		return "", func() error { return nil }, err
	}

	err = copyDirFromTo(fromDir, tmpDir)
	return tmpDir, func() error { return os.RemoveAll(tmpDir) }, err
}

// copyDirFromTo copies the contents of "from" into "to". It will call itself
// recursively for subdirectories and return the first error encountered.
// Both "from" and "to" directories should exist before it is called.
func copyDirFromTo(from, to string) error {
	d, err := os.Open(from)
	if err != nil {
		return err
	}
	defer d.Close()

	toBeCopied, err := d.Readdir(0)
	if err != nil {
		return err
	}

	for _, fi := range toBeCopied {
		if fi.IsDir() {
			var (
				newFromDir = filepath.Join(from, fi.Name())
				newToDir   = filepath.Join(to, fi.Name())
			)
			err = os.Mkdir(newToDir, fi.Mode()&fs.ModePerm)
			if err != nil {
				return err
			}
			err = copyDirFromTo(newFromDir, newToDir)
			if err != nil {
				return err
			}
		} else if fi.Mode().IsRegular() {
			var (
				fromFile = filepath.Join(from, fi.Name())
				toFile   = filepath.Join(to, fi.Name())
			)
			fromBytes, err := os.ReadFile(fromFile)
			if err != nil {
				return err
			}
			err = os.WriteFile(toFile, fromBytes, fi.Mode()&fs.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf(
				"only dirs & regular files can be copied, %q is neither",
				fi.Name())
		}
	}
	return nil
}
