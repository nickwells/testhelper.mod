package testhelper

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
)

// almostEqual returns true if a and b are within epsilon of one
// another. Copied from github.com/nickwells/mathutil.mod/mathutil
func almostEqual(a, b, epsilon float64) bool {
	if a == b {
		return true
	}

	return math.Abs(a-b) < epsilon
}

// CmpValFloat64
//
// Deprecated: use DiffFloat64
func CmpValFloat64(t *testing.T, id, name string, act, exp, epsilon float64) {
	t.Helper()
	DiffFloat64(t, id, name, act, exp, epsilon)
}

// DiffFloat64 compares the actual against the expected value and reports
// an error if they differ by more than epsilon
func DiffFloat64(t *testing.T, id, name string, act, exp, epsilon float64) bool {
	t.Helper()
	if !almostEqual(act, exp, epsilon) {
		t.Log(id)
		t.Logf("\t: expected %s: %5g\n", name, exp)
		t.Logf("\t:   actual %s: %5g\n", name, act)
		charCnt := len(name) + len("expected") + 1
		t.Logf("\t: %*s: %5g\n", charCnt, "diff", math.Abs(act-exp))
		t.Errorf("\t: %s is incorrect\n", name)
		return true
	}
	return false
}

// CmpValFloat32
//
// Deprecated: use DiffFloat32
func CmpValFloat32(t *testing.T, id, name string, act, exp, epsilon float32) {
	t.Helper()
	DiffFloat32(t, id, name, act, exp, epsilon)
}

// DiffFloat32 compares the actual against the expected value and reports
// an error if they differ by more than epsilon
func DiffFloat32(t *testing.T, id, name string, act, exp, epsilon float32) bool {
	t.Helper()
	return DiffFloat64(t, id, name, float64(act), float64(exp), float64(epsilon))
}

// CmpValInt64
//
// Deprecated: use DiffInt64
func CmpValInt64(t *testing.T, id, name string, act, exp int64) {
	t.Helper()
	DiffInt64(t, id, name, act, exp)
}

// DiffInt64 compares the actual against the expected value and reports an
// error if they differ
func DiffInt64(t *testing.T, id, name string, act, exp int64) bool {
	t.Helper()
	if act != exp {
		t.Log(id)
		t.Logf("\t: expected %s: %5d\n", name, exp)
		t.Logf("\t:   actual %s: %5d\n", name, act)
		charCnt := len(name) + len("expected") + 1
		t.Logf("\t: %*s: %5d\n", charCnt, "diff",
			int64(math.Abs(float64(act-exp))))
		t.Errorf("\t: %s is incorrect\n", name)
		return true
	}
	return false
}

// CmpValInt32
//
// Deprecated: use DiffInt32
func CmpValInt32(t *testing.T, id, name string, act, exp int32) {
	t.Helper()
	DiffInt32(t, id, name, act, exp)
}

// DiffInt32 compares the actual against the expected value and reports an
// error if they differ
func DiffInt32(t *testing.T, id, name string, act, exp int32) bool {
	t.Helper()
	return DiffInt64(t, id, name, int64(act), int64(exp))
}

// CmpValInt
//
// Deprecated: use DiffInt
func CmpValInt(t *testing.T, id, name string, act, exp int) {
	t.Helper()
	DiffInt(t, id, name, act, exp)
}

// DiffInt compares the actual against the expected value and reports an
// error if they differ
func DiffInt(t *testing.T, id, name string, act, exp int) bool {
	t.Helper()
	return DiffInt64(t, id, name, int64(act), int64(exp))
}

// CmpValUint64
//
// Deprecated: use DiffUint64
func CmpValUint64(t *testing.T, id, name string, act, exp uint64) {
	t.Helper()
	DiffUint64(t, id, name, act, exp)
}

// DiffUint64 compares the actual against the expected value and reports an
// error if they differ
func DiffUint64(t *testing.T, id, name string, act, exp uint64) bool {
	t.Helper()
	if act != exp {
		t.Log(id)
		t.Logf("\t: expected %s: %5d\n", name, exp)
		t.Logf("\t:   actual %s: %5d\n", name, act)
		t.Errorf("\t: %s is incorrect\n", name)
		return true
	}
	return false
}

// CmpValUint32
//
// Deprecated: use DiffUint32
func CmpValUint32(t *testing.T, id, name string, act, exp uint32) {
	t.Helper()
	DiffUint32(t, id, name, act, exp)
}

// DiffUint32 compares the actual against the expected value and reports an
// error if they differ
func DiffUint32(t *testing.T, id, name string, act, exp uint32) bool {
	t.Helper()
	return DiffUint64(t, id, name, uint64(act), uint64(exp))
}

// CmpValUint
//
// Deprecated: use DiffUint
func CmpValUint(t *testing.T, id, name string, act, exp uint) {
	t.Helper()
	DiffUint(t, id, name, act, exp)
}

// DiffUint compares the actual against the expected value and reports an
// error if they differ
func DiffUint(t *testing.T, id, name string, act, exp uint) bool {
	t.Helper()
	return DiffUint64(t, id, name, uint64(act), uint64(exp))
}

// CmpValString
//
// Deprecated: use DiffString
func CmpValString(t *testing.T, id, name, act, exp string) {
	t.Helper()
	DiffString(t, id, name, act, exp)
}

// DiffString compares the actual against the expected value and reports an
// error if they differ
func DiffString(t *testing.T, id, name, act, exp string) bool {
	t.Helper()
	if act != exp {
		t.Log(id)
		t.Logf("\t: expected %s (length: %4d): %q\n", name, len(exp), exp)
		t.Logf("\t:   actual %s (length: %4d): %q\n", name, len(act), act)
		t.Errorf("\t: %s is incorrect\n", name)
		return true
	}
	return false
}

// DiffStringer compares the actual against the expected value and reports an
// error if they differ. It will report them as different if one is nil or
// has a nil value and the other isn't/doesn't or if they are both non-nil
// and the string values differ.
func DiffStringer(t *testing.T, id, name string, actS, expS fmt.Stringer) bool {
	t.Helper()

	actIsNil := actS == nil || reflect.ValueOf(actS).IsNil()
	expIsNil := expS == nil || reflect.ValueOf(expS).IsNil()

	if actIsNil && expIsNil {
		return false
	}

	if actIsNil {
		t.Log(id)
		t.Logf("\t: expected %s is non-nil\n", name)
		t.Logf("\t:   actual %s is nil\n", name)
		t.Errorf("\t: %s is incorrect\n", name)
		return true
	}

	if expIsNil {
		t.Log(id)
		t.Logf("\t: expected %s is nil\n", name)
		t.Logf("\t:   actual %s is non-nil\n", name)
		t.Errorf("\t: %s is incorrect\n", name)
		return true
	}

	return DiffString(t, id, name, actS.String(), expS.String())
}

// DiffStringSlice compares the actual against the expected value and reports an
// error if they differ
func DiffStringSlice(t *testing.T, id, name string, act, exp []string) bool {
	t.Helper()
	if len(act) != len(exp) {
		t.Log(id)
		t.Logf("\t: expected %s length: %4d\n", name, len(exp))
		t.Logf("\t:   actual %s length: %4d\n", name, len(act))
		t.Errorf("\t: %s is incorrect\n", name)
		return true
	}

	diffCount := 0
	for i, s := range act {
		if s != exp[i] {
			diffCount++
			if diffCount == 1 {
				t.Log(id)
			}
			t.Logf("\t: expected %s [%d]: %q\n", name, i, exp[i])
			t.Logf("\t:   actual %s [%d]: %q\n", name, i, s)
		}
	}
	if diffCount > 0 {
		diffStr := "differences"
		if diffCount == 1 {
			diffStr = "difference"
		}
		t.Logf("\t: %d %s found\n", diffCount, diffStr)
		t.Errorf("\t: %s is incorrect\n", name)
	}
	return diffCount > 0
}

// CmpValBool
//
// Deprecated: use DiffBool
func CmpValBool(t *testing.T, id, name string, act, exp bool) {
	t.Helper()
	DiffBool(t, id, name, act, exp)
}

// DiffBool compares the actual against the expected value and reports an
// error if they differ
func DiffBool(t *testing.T, id, name string, act, exp bool) bool {
	t.Helper()
	if act != exp {
		t.Log(id)
		t.Logf("\t: expected %s: %v\n", name, exp)
		t.Logf("\t:   actual %s: %v\n", name, act)
		t.Errorf("\t: %s is incorrect\n", name)
		return true
	}
	return false
}

// DiffTime compares the actual against the expected value and reports an
// error if they differ
func DiffTime(t *testing.T, id, name string, act, exp time.Time) bool {
	t.Helper()
	if d := act.Sub(exp); d != 0 {
		t.Log(id)
		t.Logf("\t: expected %s: %v\n", name, exp)
		t.Logf("\t:   actual %s: %v\n", name, act)
		t.Logf("\t: difference: %v\n", d)
		t.Errorf("\t: %s is incorrect\n", name)
		return true
	}
	return false
}
