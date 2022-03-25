package testhelper

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"golang.org/x/exp/constraints"
)

// almostEqual returns true if a and b are within epsilon of one
// another. Copied from github.com/nickwells/mathutil.mod/mathutil
func almostEqual[T constraints.Float](a, b, epsilon T) bool {
	if a == b {
		return true
	}

	return math.Abs(float64(a-b)) < float64(epsilon)
}

// reportFloatDiff reports the difference between two float values
func reportFloatDiff[T constraints.Float](t *testing.T, name string,
	act, exp T,
) {
	t.Helper()

	t.Logf("\t: expected %s: %5g\n", name, exp)
	t.Logf("\t:   actual %s: %5g\n", name, act)
	charCnt := len(name) + len("expected") + 1
	t.Logf("\t: %*s: %5g\n", charCnt, "diff", math.Abs(float64(act-exp)))
	t.Errorf("\t: %s is incorrect\n", name)
}

// DiffFloat compares the actual against the expected value and reports
// an error if they differ by more than epsilon
func DiffFloat[T constraints.Float](t *testing.T, id, name string,
	act, exp, epsilon T,
) bool {
	t.Helper()
	if !almostEqual(act, exp, epsilon) {
		t.Log(id)
		reportFloatDiff(t, name, act, exp)
		return true
	}
	return false
}

// DiffInt compares the actual against the expected value and reports an
// error if they differ
func DiffInt[T constraints.Integer](t *testing.T, id, name string,
	act, exp T,
) bool {
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

// reportStringDiff reports the difference between two strings
func reportStringDiff(t *testing.T, name, act, exp string) {
	t.Helper()
	t.Logf("\t: expected %s (length: %4d): %q\n", name, len(exp), exp)
	t.Logf("\t:   actual %s (length: %4d): %q\n", name, len(act), act)
	t.Errorf("\t: %s is incorrect\n", name)
}

// DiffString compares the actual against the expected value and reports an
// error if they differ
func DiffString(t *testing.T, id, name, act, exp string) bool {
	t.Helper()
	if act != exp {
		t.Log(id)
		reportStringDiff(t, name, act, exp)
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

// DiffErr compares the actual against the expected value and reports an
// error if they differ. Note that it compares the string representation and
// not the error type so there might be a mismatch.
func DiffErr(t *testing.T, id, name string, act, exp error) bool {
	t.Helper()

	if act == nil && exp == nil {
		return false
	}

	if act == nil && exp != nil ||
		act != nil && exp == nil ||
		act.Error() != exp.Error() {
		t.Log(id)
		t.Logf("\t: expected %s: %v\n", name, exp)
		t.Logf("\t:   actual %s: %v\n", name, act)
		t.Errorf("\t: %s is incorrect\n", name)
		return true
	}
	return false
}

const MaxReportedDiffs = 5

// reportDiffCount reports the number of differences found
func reportDiffCount(t *testing.T, name string, diffCount int) {
	t.Helper()
	if diffCount > 0 {
		diffStr := "differences"
		if diffCount == 1 {
			diffStr = "difference"
		}
		t.Logf("\t: %d %s found\n", diffCount, diffStr)
		t.Errorf("\t: %s is incorrect\n", name)
	}
}

// reportMaxDiffsShown logs an elipsis to show that more differences have
// been found but that they have been elided. This is only done for the first
// elided difference.
func reportMaxDiffsShown(t *testing.T, diffCount int) {
	t.Helper()
	if diffCount == (MaxReportedDiffs + 1) {
		t.Log("\t: ...\n")
	}
}

// reportSliceLenDiff reports if the lengths of the slices differ. It returns
// true if so and false otherwise.
func reportSliceLenDiff(t *testing.T, id, name string, act, exp int) bool {
	t.Helper()
	if act != exp {
		t.Log(id)
		t.Logf("\t: expected %s length: %4d\n", name, exp)
		t.Logf("\t:   actual %s length: %4d\n", name, act)
		t.Errorf("\t: %s is incorrect\n", name)
		return true
	}
	return false
}

// DiffInt64Slice compares the actual against the expected value and reports
// an error if they differ. At most MaxReportedDiffs are reported.
func DiffInt64Slice(t *testing.T, id, name string, act, exp []int64) bool {
	t.Helper()
	if reportSliceLenDiff(t, id, name, len(act), len(exp)) {
		return true
	}

	diffCount := 0
	for i, v := range act {
		if v != exp[i] {
			diffCount++
			if diffCount > MaxReportedDiffs {
				reportMaxDiffsShown(t, diffCount)
				continue
			}
			if diffCount == 1 {
				t.Log(id)
			}
			t.Logf("\t: expected %s [%d]: %d\n", name, i, exp[i])
			t.Logf("\t:   actual %s [%d]: %d\n", name, i, v)
		}
	}
	reportDiffCount(t, name, diffCount)

	return diffCount > 0
}

// DiffFloat64Slice compares the actual against the expected value and reports
// an error if they differ. At most MaxReportedDiffs are reported.
func DiffFloat64Slice(t *testing.T, id, name string, act, exp []float64, epsilon float64) bool {
	t.Helper()
	if reportSliceLenDiff(t, id, name, len(act), len(exp)) {
		return true
	}

	diffCount := 0
	for i, v := range act {
		if !almostEqual(v, exp[i], epsilon) {
			diffCount++
			if diffCount > MaxReportedDiffs {
				reportMaxDiffsShown(t, diffCount)
				continue
			}
			if diffCount == 1 {
				t.Log(id)
			}
			reportFloatDiff(t, fmt.Sprintf("%s [%d]", name, i), v, exp[i])
		}
	}
	reportDiffCount(t, name, diffCount)

	return diffCount > 0
}

// DiffStringSlice compares the actual against the expected value and reports
// an error if they differ. At most MaxReportedDiffs are reported.
func DiffStringSlice(t *testing.T, id, name string, act, exp []string) bool {
	t.Helper()
	if reportSliceLenDiff(t, id, name, len(act), len(exp)) {
		return true
	}

	diffCount := 0
	for i, s := range act {
		if s != exp[i] {
			diffCount++
			if diffCount > MaxReportedDiffs {
				reportMaxDiffsShown(t, diffCount)
				continue
			}
			if diffCount == 1 {
				t.Log(id)
			}
			reportStringDiff(t, fmt.Sprintf("%s [%d]", name, i), s, exp[i])
		}
	}
	reportDiffCount(t, name, diffCount)
	return diffCount > 0
}
