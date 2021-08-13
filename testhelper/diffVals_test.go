package testhelper_test

import (
	"testing"
	"unsafe"

	"github.com/nickwells/testhelper.mod/testhelper"
)

// myFunc exists only to test the behaviour of func comparisons
func myFunc() int {
	return 0
}

// otherFunc exists only to test the behaviour of func comparisons
func otherFunc() int {
	return 0
}

// myStructLoopy exists only to test the behaviour of struct comparisons
type myStructLoopy struct {
	i   int
	f   float64
	msl *myStructLoopy
}

// myStructSimple exists only to test the behaviour of struct comparisons
type myStructSimple struct {
	i int
	f float64
}

// myStructNestedUnnamed exists only to test the behaviour of struct comparisons
type myStructNestedUnnamed struct {
	s string
	myStructSimple
}

// myStructNestedNamed exists only to test the behaviour of struct comparisons
type myStructNestedNamed struct {
	s   string
	mss myStructSimple
}

func TestDiffVals(t *testing.T) {
	i := 0
	j := 0
	k := 1

	loopyMap := make(map[string]interface{})
	loopyMap["loop"] = loopyMap
	var nilMap map[string]interface{}

	var loopyArray1 [3]interface{}
	loopyArray1[0] = &loopyArray1
	var loopyArray2 [3]interface{}
	loopyArray2[0] = &loopyArray2

	loopySlice1 := make([]interface{}, 0)
	loopySlice1 = append(loopySlice1, &loopySlice1)
	loopySlice2 := make([]interface{}, 0)
	loopySlice2 = append(loopySlice2, &loopySlice2)
	var nilSlice []interface{}

	msl1 := myStructLoopy{i: 42, f: 3.14159}
	msl1.msl = &msl1
	msl2 := myStructLoopy{i: 42, f: 3.14159}
	msl2.msl = &msl2

	mssi42 := myStructSimple{i: 42, f: 3.14159}
	mssi99 := myStructSimple{i: 99, f: 3.14159}

	msnui42 := myStructNestedUnnamed{s: "hello", myStructSimple: mssi42}
	msnui99 := myStructNestedUnnamed{s: "hello", myStructSimple: mssi99}

	msnni42 := myStructNestedNamed{s: "hello", mss: mssi42}
	msnni99 := myStructNestedNamed{s: "hello", mss: mssi99}

	chan1 := make(chan bool)
	chan2 := make(chan bool)

	var nilChan chan bool

	testCases := []struct {
		testhelper.ID
		actVal interface{}
		expVal interface{}
		ignore [][]string
		testhelper.ExpErr
	}{
		{
			ID: testhelper.MkID("both nil"),
		},
		{
			ID:     testhelper.MkID("exp nil, act not"),
			actVal: 42,
			ExpErr: testhelper.MkExpErr("the expected value is nil," +
				" the actual value is not"),
		},
		{
			ID:     testhelper.MkID("act nil, exp not"),
			expVal: 42,
			ExpErr: testhelper.MkExpErr("the actual value is nil," +
				" the expected value is not"),
		},
		{
			ID:     testhelper.MkID("types differ, array by length"),
			actVal: [...]int{1, 2, 3},
			expVal: [...]int{1, 2, 3, 4},
			ExpErr: testhelper.MkExpErr(`this: types differ.`,
				"Actual: [3]int, expected: [4]int"),
		},
		{
			ID:     testhelper.MkID("types differ"),
			actVal: 3.14159,
			expVal: 42,
			ExpErr: testhelper.MkExpErr(`this: types differ.`,
				"Actual: float64, expected: int"),
		},
		{
			ID:     testhelper.MkID("same val, bool"),
			actVal: true,
			expVal: true,
		},
		{
			ID:     testhelper.MkID("vals differ, bool"),
			actVal: true,
			expVal: false,
			ExpErr: testhelper.MkExpErr(`this: bool values differ.`,
				"Actual: true, expected: false"),
		},
		{
			ID:     testhelper.MkID("same val, int"),
			actVal: 42,
			expVal: 42,
		},
		{
			ID:     testhelper.MkID("vals differ, int"),
			actVal: 42,
			expVal: 43,
			ExpErr: testhelper.MkExpErr(`this: int values differ.`,
				"Actual: 42, expected: 43"),
		},
		{
			ID:     testhelper.MkID("same val, uint"),
			actVal: uint(42),
			expVal: uint(42),
		},
		{
			ID:     testhelper.MkID("vals differ, uint"),
			actVal: uint(42),
			expVal: uint(43),
			ExpErr: testhelper.MkExpErr(`this: uint values differ.`,
				"Actual: 42, expected: 43"),
		},
		{
			ID:     testhelper.MkID("same val, float"),
			actVal: 3.14159,
			expVal: 3.14159,
		},
		{
			ID:     testhelper.MkID("vals differ, float"),
			actVal: 3.14159,
			expVal: 4.14159,
			ExpErr: testhelper.MkExpErr(`this: float values differ.`,
				"Actual: 3.14159, expected: 4.14159"),
		},
		{
			ID:     testhelper.MkID("same val, complex"),
			actVal: complex(1, 2),
			expVal: complex(1, 2),
		},
		{
			ID:     testhelper.MkID("vals differ, complex"),
			actVal: complex(1, 2),
			expVal: complex(3, 4),
			ExpErr: testhelper.MkExpErr(`this: complex values differ.`,
				"Actual: (1+2i), expected: (3+4i)"),
		},
		{
			ID:     testhelper.MkID("same val, map"),
			actVal: map[string]interface{}{"a": "A", "b": 42},
			expVal: map[string]interface{}{"a": "A", "b": 42},
		},
		{
			ID:     testhelper.MkID("map, act nil"),
			actVal: nilMap,
			expVal: map[string]interface{}{},
		},
		{
			ID:     testhelper.MkID("loopy map"),
			actVal: loopyMap,
			expVal: loopyMap,
		},
		{
			ID:     testhelper.MkID("value diff by len, map"),
			actVal: map[string]interface{}{"a": "A", "b": 42},
			expVal: map[string]interface{}{"a": "A"},
			ExpErr: testhelper.MkExpErr(`this: map lengths differ.`,
				`Actual: 2, expected: 1`),
		},
		{
			ID:     testhelper.MkID("value diff by keys, map"),
			actVal: map[string]interface{}{"a": "A", "b": 42},
			expVal: map[string]interface{}{"a": "A", "c": 42},
			ExpErr: testhelper.MkExpErr(`this[b]:` +
				` the expected value is invalid, the actual value is not`),
		},
		{
			ID:     testhelper.MkID("value diff by value, map"),
			actVal: map[string]interface{}{"a": "A", "b": 42},
			expVal: map[string]interface{}{"a": "Not-A", "b": 42},
			ExpErr: testhelper.MkExpErr(`this[a]: strings differ.`,
				`Actual: "A", expected: "Not-A"`),
		},
		{
			ID:     testhelper.MkID("value diff by type, map"),
			actVal: map[string]interface{}{"a": "A", "b": 42},
			expVal: map[string]interface{}{"a": 3.14159, "b": 42},
			ExpErr: testhelper.MkExpErr(`this[a]: types differ.`,
				"Actual: string, expected: float64"),
		},
		{
			ID:     testhelper.MkID("same val, array"),
			actVal: [...]interface{}{42, 3.14159, "Hello", "World", myFunc},
			expVal: [...]interface{}{42, 3.14159, "Hello", "World", myFunc},
		},
		{
			ID:     testhelper.MkID("loopy array"),
			actVal: loopyArray1,
			expVal: loopyArray2,
		},
		{
			ID:     testhelper.MkID("value diff by value, array"),
			actVal: [...]interface{}{42, 3.14159, "Hello", "World", myFunc},
			expVal: [...]interface{}{1, 3.14159, "Hello", "World", myFunc},
			ExpErr: testhelper.MkExpErr(`this[0]: int values differ.`,
				"Actual: 42, expected: 1"),
		},
		{
			ID:     testhelper.MkID("value diff by type, array"),
			actVal: [...]interface{}{42, 3.14159, "Hello", "World", myFunc},
			expVal: [...]interface{}{1.2, 3.14159, "Hello", "World", myFunc},
			ExpErr: testhelper.MkExpErr(`this[0]: types differ.`,
				"Actual: int, expected: float64"),
		},
		{
			ID:     testhelper.MkID("same val, slice"),
			actVal: []interface{}{42, 3.14159, "Hello", "World", myFunc},
			expVal: []interface{}{42, 3.14159, "Hello", "World", myFunc},
		},
		{
			ID:     testhelper.MkID("slice, act nil"),
			actVal: nilSlice,
			expVal: []interface{}{},
		},
		{
			ID:     testhelper.MkID("loopy slice"),
			actVal: loopySlice1,
			expVal: loopySlice2,
		},
		{
			ID:     testhelper.MkID("value diff by length, slice"),
			actVal: []interface{}{42, 3.14159, "Hello", "World", myFunc},
			expVal: []interface{}{42, 3.14159, "Hello", "World"},
			ExpErr: testhelper.MkExpErr(`this: slice lengths differ.`,
				"Actual: 5, expected: 4"),
		},
		{
			ID:     testhelper.MkID("value diff by value, slice"),
			actVal: []interface{}{42, 3.14159, "Hello", "World", myFunc},
			expVal: []interface{}{1, 3.14159, "Hello", "World", myFunc},
			ExpErr: testhelper.MkExpErr(`this[0]: int values differ.`,
				"Actual: 42, expected: 1"),
		},
		{
			ID:     testhelper.MkID("value diff by type, slice"),
			actVal: []interface{}{42, 3.14159, "Hello", "World", myFunc},
			expVal: []interface{}{1.2, 3.14159, "Hello", "World", myFunc},
			ExpErr: testhelper.MkExpErr(`this[0]: types differ.`,
				"Actual: int, expected: float64"),
		},
		{
			ID:     testhelper.MkID("same val, func"),
			actVal: myFunc,
			expVal: myFunc,
		},
		{
			ID:     testhelper.MkID("vals differ, func"),
			actVal: myFunc,
			expVal: otherFunc,
			ExpErr: testhelper.MkExpErr(`this: funcs differ.`,
				"Actual instance is not equal to expected"),
		},
		{
			ID:     testhelper.MkID("same val, string"),
			actVal: "Hello",
			expVal: "Hello",
		},
		{
			ID:     testhelper.MkID("vals differ, string"),
			actVal: "Hello",
			expVal: "Goodbye",
			ExpErr: testhelper.MkExpErr(`this: strings differ.`,
				`Actual: "Hello", expected: "Goodbye"`),
		},
		{
			ID:     testhelper.MkID("same pointer, unsafePointer"),
			actVal: unsafe.Pointer(&i),
			expVal: unsafe.Pointer(&i),
		},
		{
			ID:     testhelper.MkID("diff pointer, unsafePointer"),
			actVal: unsafe.Pointer(&i),
			expVal: unsafe.Pointer(&j),
			ExpErr: testhelper.MkExpErr(`this: pointers differ`),
		},
		{
			ID:     testhelper.MkID("same pointer, uintptr"),
			actVal: uintptr(unsafe.Pointer(&i)),
			expVal: uintptr(unsafe.Pointer(&i)),
		},
		{
			ID:     testhelper.MkID("diff pointer, uintptr"),
			actVal: uintptr(unsafe.Pointer(&i)),
			expVal: uintptr(unsafe.Pointer(&j)),
			ExpErr: testhelper.MkExpErr(`this: uintptr pointers differ`),
		},
		{
			ID:     testhelper.MkID("same pointer, ptr"),
			actVal: &i,
			expVal: &i,
		},
		{
			ID:     testhelper.MkID("diff pointer, same val, ptr"),
			actVal: &i,
			expVal: &j,
		},
		{
			ID:     testhelper.MkID("diff pointer, diff value, ptr"),
			actVal: &i,
			expVal: &k,
			ExpErr: testhelper.MkExpErr(`this: int values differ.`,
				"Actual: 0, expected: 1"),
		},
		{
			ID:     testhelper.MkID("same val, simple struct"),
			actVal: mssi42,
			expVal: mssi42,
		},
		{
			ID:     testhelper.MkID("loopy struct"),
			actVal: msl1,
			expVal: msl2,
		},
		{
			ID:     testhelper.MkID("vals differ, simple struct"),
			actVal: mssi42,
			expVal: mssi99,
			ExpErr: testhelper.MkExpErr(`this.i: int values differ.`,
				"Actual: 42, expected: 99"),
		},
		{
			ID:     testhelper.MkID("same val, nested, unnamed"),
			actVal: msnui42,
			expVal: msnui42,
		},
		{
			ID:     testhelper.MkID("vals differ, nested, unnamed"),
			actVal: msnui42,
			expVal: msnui99,
			ExpErr: testhelper.MkExpErr(
				`this.myStructSimple.i: int values differ.`,
				"Actual: 42, expected: 99"),
		},
		{
			ID:     testhelper.MkID("vals differ, ignored, nested, unnamed"),
			actVal: msnui42,
			expVal: msnui99,
			ignore: [][]string{{"myStructSimple", "i"}},
		},
		{
			ID:     testhelper.MkID("same val, nested, named"),
			actVal: msnni42,
			expVal: msnni42,
		},
		{
			ID:     testhelper.MkID("vals differ, nested, named"),
			actVal: msnni42,
			expVal: msnni99,
			ExpErr: testhelper.MkExpErr(`this.mss.i: int values differ.`,
				"Actual: 42, expected: 99"),
		},
		{
			ID:     testhelper.MkID("vals differ, nested, named"),
			actVal: msnni42,
			expVal: msnni99,
			ExpErr: testhelper.MkExpErr(
				`this.mss.i: int values differ.`,
				"Actual: 42, expected: 99"),
		},
		{
			ID:     testhelper.MkID("vals differ, other ignored, nested, named"),
			actVal: msnni42,
			expVal: msnni99,
			ignore: [][]string{{"mss", "f"}},
			ExpErr: testhelper.MkExpErr(
				`this.mss.i: int values differ.`,
				"Actual: 42, expected: 99"),
		},
		{
			ID:     testhelper.MkID("vals differ, ignored, nested, named"),
			actVal: msnni42,
			expVal: msnni99,
			ignore: [][]string{{"mss", "i"}},
		},
		{
			ID:     testhelper.MkID("chan, same val"),
			actVal: chan1,
			expVal: chan1,
		},
		{
			ID:     testhelper.MkID("chan, vals differ"),
			actVal: chan1,
			expVal: chan2,
			ExpErr: testhelper.MkExpErr(`this: chans differ.`,
				`Actual instance is not equal to expected`),
		},
		{
			ID:     testhelper.MkID("chan, vals differ, act nil"),
			actVal: nilChan,
			expVal: chan2,
			ExpErr: testhelper.MkExpErr(`this: chans differ.`,
				`Actual instance is not equal to expected`),
		},
		{
			ID:     testhelper.MkID("chan, vals differ, exp nil"),
			actVal: chan1,
			expVal: nilChan,
			ExpErr: testhelper.MkExpErr(`this: chans differ.`,
				`Actual instance is not equal to expected`),
		},
	}

	for _, tc := range testCases {
		err := testhelper.DiffVals(tc.actVal, tc.expVal, tc.ignore...)
		testhelper.CheckExpErr(t, err, tc)
	}
}
