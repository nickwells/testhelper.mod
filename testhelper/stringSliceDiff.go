package testhelper

// StringSliceDiff will compare two slices of strings for equality. If they
// are of different lengths they are taken to be different. A nil slice and
// an empty slice are taken to be the same.
func StringSliceDiff(a, b []string) bool {
	if len(a) != len(b) {
		return true
	}

	for i, v := range a {
		if b[i] != v {
			return true
		}
	}
	return false
}
