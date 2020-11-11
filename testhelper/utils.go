package testhelper

// ErrSliceToStrSlice creates a string slice from the supplied slice of
// errors - each string is the result of the Error func being called on the
// corresponding error
func ErrSliceToStrSlice(es []error) []string {
	ss := make([]string, 0, len(es))

	for _, err := range es {
		ss = append(ss, err.Error())
	}

	return ss
}
