package testhelper

import (
	"math"
	"testing"
)

// almostEqual returns true if a and b are within epsilon of one
// another. Copied from github.com/nickwells/mathutil.mod/mathutil
func almostEqual(a, b, epsilon float64) bool {
	if a == b {
		return true
	}

	return math.Abs(a-b) < epsilon
}

// CmpValFloat64 compares the actual against the expected value and reports
// an error if they differ by more than epsilon
func CmpValFloat64(t *testing.T, id, name string, act, exp, epsilon float64) {
	t.Helper()
	if !almostEqual(act, exp, epsilon) {
		t.Log(id)
		t.Logf("\t: expected %s: %5g\n", name, exp)
		t.Logf("\t:   actual %s: %5g\n", name, act)
		charCnt := len(name) + len("expected") + 1
		t.Logf("\t: %*s: %5g\n", charCnt, "diff", math.Abs(act-exp))
		t.Errorf("\t: %s is incorrect\n", name)
	}
}

// CmpValFloat32 compares the actual against the expected value and reports
// an error if they differ by more than epsilon
func CmpValFloat32(t *testing.T, id, name string, act, exp, epsilon float32) {
	t.Helper()
	CmpValFloat64(t, id, name, float64(act), float64(exp), float64(epsilon))
}

// CmpValInt64 compares the actual against the expected value and reports an
// error if they differ
func CmpValInt64(t *testing.T, id, name string, act, exp int64) {
	t.Helper()
	if act != exp {
		t.Log(id)
		t.Logf("\t: expected %s: %5d\n", name, exp)
		t.Logf("\t:   actual %s: %5d\n", name, act)
		t.Errorf("\t: %s is incorrect\n", name)
	}
}

// CmpValInt32 compares the actual against the expected value and reports an
// error if they differ
func CmpValInt32(t *testing.T, id, name string, act, exp int32) {
	t.Helper()
	CmpValInt64(t, id, name, int64(act), int64(exp))
}

// CmpValInt compares the actual against the expected value and reports an
// error if they differ
func CmpValInt(t *testing.T, id, name string, act, exp int) {
	t.Helper()
	CmpValInt64(t, id, name, int64(act), int64(exp))
}

// CmpValUint64 compares the actual against the expected value and reports an
// error if they differ
func CmpValUint64(t *testing.T, id, name string, act, exp uint64) {
	t.Helper()
	if act != exp {
		t.Log(id)
		t.Logf("\t: expected %s: %5d\n", name, exp)
		t.Logf("\t:   actual %s: %5d\n", name, act)
		t.Errorf("\t: %s is incorrect\n", name)
	}
}

// CmpValUint32 compares the actual against the expected value and reports an
// error if they differ
func CmpValUint32(t *testing.T, id, name string, act, exp uint32) {
	t.Helper()
	CmpValUint64(t, id, name, uint64(act), uint64(exp))
}

// CmpValUint compares the actual against the expected value and reports an
// error if they differ
func CmpValUint(t *testing.T, id, name string, act, exp uint) {
	t.Helper()
	CmpValUint64(t, id, name, uint64(act), uint64(exp))
}

// CmpValString compares the actual against the expected value and reports an
// error if they differ
func CmpValString(t *testing.T, id, name, act, exp string) {
	t.Helper()
	if act != exp {
		t.Log(id)
		t.Logf("\t: expected %s: %q\n", name, exp)
		t.Logf("\t:   actual %s: %q\n", name, act)
		t.Errorf("\t: %s is incorrect\n", name)
	}
}

// CmpValBool compares the actual against the expected value and reports an
// error if they differ
func CmpValBool(t *testing.T, id, name string, act, exp bool) {
	t.Helper()
	if act != exp {
		t.Log(id)
		t.Logf("\t: expected %s: %v\n", name, exp)
		t.Logf("\t:   actual %s: %v\n", name, act)
		t.Errorf("\t: %s is incorrect\n", name)
	}
}
