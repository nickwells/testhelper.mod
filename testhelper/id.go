package testhelper

import "fmt"

// ID holds common identifying information about a test
type ID struct {
	Name string
}

// MakeID returns an identifying string describing the test
func (id ID) MakeID(i int) string {
	return fmt.Sprintf("test %d: %s", i, id.Name)

}

// TestCase is an interface wrapping the MakeID method
type TestCase interface {
	MakeID(int) string
}
