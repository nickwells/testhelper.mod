package testhelper

import "os"

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
