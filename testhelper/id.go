package testhelper

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// ID holds common identifying information about a test. Several of the
// testhelper functions take an ID (or an interface which it satisfies) which
// can simplify the tests. This is often combined with other testhelper mixin
// structs to record standard details for common tests.
type ID struct {
	Name       string
	At         string
	AtFullName string
}

// MkID is a constructor for the ID type. It will record where it was called
// from and the resulting ID, when used to report an error, will give
// information on the location which should speed up locating the test setup
// for the failing test.
func MkID(name string) ID {
	id := ID{Name: name}
	if _, file, line, ok := runtime.Caller(1); ok {
		id.At = fmt.Sprintf("%s:%d", filepath.Base(file), line)
		id.AtFullName = fmt.Sprintf("%s:%d", file, line)
	}
	return id
}

// IDStr returns an identifying string describing the test
func (id ID) IDStr() string {
	if id.At == "" {
		return "test: " + id.Name
	}
	return "test: " + id.At + ":" + id.Name
}

// IDStrFullName returns an identifying string describing the test, using the
// full pathname of the file where the MkID func was called rather than just
// the last part of the path. You might want to use this if your test cases
// are initialised in a more complex way and there is some ambiguity as to
// the location of a source file.
func (id ID) IDStrFullName() string {
	if id.AtFullName == "" {
		return "test: " + id.Name
	}
	return "test: " + id.AtFullName + ":" + id.Name
}

// TestCase is an interface wrapping the IDStr methods
type TestCase interface {
	IDStr() string
	IDStrFullName() string
}
