package testhelper

import (
	"fmt"
	"reflect"
)

// DiffValErr holds an error reflecting the difference between two interface
// values. It is the type of error returned by the DiffVals func
type DiffValErr struct {
	dl  deepLoc
	msg string
}

// Error returns a string expressing the DiffValErr
func (dve DiffValErr) Error() string {
	return dve.dl.fullName + ": " + dve.msg
}

// visit represents a previously visited pair of pointers and allows us to
// avoid looping around the same circle of pointers for ever.
type visit struct {
	actPtr uintptr
	expPtr uintptr
	T      reflect.Type
}

// deepLoc represents the current location in the chain of values. It also
// holds the collection of locations to skip and the current location can be
// checked against this skip list. This is to allow comparisons to be made
// between values where some of the fields in a structure may be different
// but we don't care.
type deepLoc struct {
	depth      int
	fullName   string
	currentLoc []string
	skipLocs   [][]string
	loop       map[visit]bool
}

// makeDeepLoc returns a correctly formed deepLoc
func makeDeepLoc(sl [][]string) deepLoc {
	return deepLoc{
		fullName: "this",
		skipLocs: sl,
		loop:     map[visit]bool{},
	}
}

// String returns the full name of the location which includes index values
// which are not checked when deciding if the location should be skipped
func (dl deepLoc) String() string {
	return dl.fullName
}

// incr increments the number of times the code has called diffVals and
// panics if it's more than maxDepth
func (dl *deepLoc) incr() {
	const maxDepth = 1000

	dl.depth++

	if dl.depth > maxDepth {
		panic(DiffValErr{dl: *dl, msg: "undetected loop"})
	}
}

// isALoop returns true if the visit has already been seen, otherwise it adds
// it to the loop member
func (dl *deepLoc) isALoop(v visit) bool {
	if v.actPtr == uintptr(0) || v.expPtr == uintptr(0) {
		return false
	}

	if dl.loop[v] {
		return true
	}

	dl.loop[v] = true

	return false
}

// addName adds the name to the formatted name and to the currentLoc. It
// returns a copy of the amended value.
func addName(dl deepLoc, s string) deepLoc {
	dl.fullName += "." + s
	dl.currentLoc = append(dl.currentLoc, s)

	return dl
}

// addIdx adds the index to the formatted name. It returns a copy of the
// amended value.
func addIdx(dl deepLoc, i int) deepLoc {
	dl.fullName += fmt.Sprintf("[%d]", i)
	return dl
}

// addKey adds the key to the formatted name. It returns a copy of the
// amended value.
func addKey(dl deepLoc, k any) deepLoc {
	dl.fullName += fmt.Sprintf("[%v]", k)
	return dl
}

// skip returns true if the current location is in the list of locations to
// be skipped.
func (dl deepLoc) skip() bool {
	for _, sl := range dl.skipLocs {
		if len(sl) > len(dl.currentLoc) {
			continue
		}

		skip := true

		for i, s := range sl {
			if s != dl.currentLoc[i] {
				skip = false
				break
			}
		}

		if skip {
			return true
		}
	}

	return false
}

// DiffVals compares the actual and expected values and returns an error if
// they are different. This differs from the reflect package function
// DeepEqual in that the error shows which fields are different.
//
// The ignore argument holds the names of fields to not compare, the slice
// represents a chain of names in nested structs. So, for instance an ignore
// value containing the pair of values ["a", "b"] means to not compare the
// field called "b" in the sub-struct called "a".
func DiffVals(actVal, expVal any, ignore ...[]string) error {
	if actVal == nil && expVal == nil {
		return nil
	}

	dl := makeDeepLoc(ignore)
	if actVal == nil {
		return DiffValErr{
			dl:  dl,
			msg: "the actual value is nil, the expected value is not",
		}
	}

	if expVal == nil {
		return DiffValErr{
			dl:  dl,
			msg: "the expected value is nil, the actual value is not",
		}
	}

	return diffVals(reflect.ValueOf(actVal), reflect.ValueOf(expVal), dl)
}

// diffVals returns an error if the two values differ, either in type or value
func diffVals(actVal, expVal reflect.Value, dl deepLoc) error { //nolint:cyclop
	if dl.skip() {
		return nil
	}

	dl.incr()

	if !actVal.IsValid() && !expVal.IsValid() {
		return nil
	}

	if !actVal.IsValid() {
		return DiffValErr{
			dl:  dl,
			msg: "the actual value is invalid, the expected value is not",
		}
	}

	if !expVal.IsValid() {
		return DiffValErr{
			dl:  dl,
			msg: "the expected value is invalid, the actual value is not",
		}
	}

	actType := actVal.Type()
	expType := expVal.Type()

	if actType != expType {
		return DiffValErr{
			dl: dl,
			msg: fmt.Sprintf("types differ. Actual: %s, expected: %s",
				actType, expType),
		}
	}

	switch actType.Kind() {
	case reflect.Invalid:
		return nil
	case reflect.Bool:
		return diffValsBool(actVal, expVal, dl)
	case reflect.Int,
		reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return diffValsInt(actVal, expVal, dl)
	case reflect.Uint,
		reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		return diffValsUint(actVal, expVal, dl)
	case reflect.Float32, reflect.Float64:
		return diffValsFloat(actVal, expVal, dl)
	case reflect.Complex64, reflect.Complex128:
		return diffValsComplex(actVal, expVal, dl)
	case reflect.String:
		return diffValsString(actVal, expVal, dl)
	case reflect.Interface:
		return diffVals(actVal.Elem(), expVal.Elem(), dl)
	case reflect.Array:
		return diffValsArray(actVal, expVal, dl)
	case reflect.Slice:
		if valsMustBeEqual(actVal.Pointer(), expVal.Pointer(), actType, dl) {
			return nil
		}

		return diffValsSlice(actVal, expVal, dl)
	case reflect.Func:
		return diffValsFunc(actVal, expVal, dl)
	case reflect.Chan:
		return diffValsChan(actVal, expVal, dl)
	case reflect.UnsafePointer:
		return diffValsPointer(actVal, expVal, dl)
	case reflect.Uintptr:
		return diffValsUintptr(actVal, expVal, dl)
	case reflect.Ptr:
		if valsMustBeEqual(actVal.Pointer(), expVal.Pointer(), actType, dl) {
			return nil
		}

		return diffVals(actVal.Elem(), expVal.Elem(), dl)
	case reflect.Map:
		if valsMustBeEqual(actVal.Pointer(), expVal.Pointer(), actType, dl) {
			return nil
		}

		return diffValsMap(actVal, expVal, dl)
	case reflect.Struct:
		return diffValsStruct(actVal, expVal, dl)
	}

	panic(DiffValErr{
		dl:  dl,
		msg: fmt.Sprintf("unchecked value kind: %s", actType.Kind()),
	})
}

// valsMustBeEqual returns true if the pointers are the same or if they have
// been visited before
func valsMustBeEqual(actPtr, expPtr uintptr, t reflect.Type, dl deepLoc) bool {
	if actPtr == expPtr {
		return true
	}

	if dl.isALoop(visit{actPtr, expPtr, t}) {
		return true
	}

	return false
}

// diffValsMap returns an error if the two map values differ
func diffValsMap(actVal, expVal reflect.Value, dl deepLoc) error {
	aMapKeys := actVal.MapKeys()
	eMapKeys := expVal.MapKeys()
	aLen := len(aMapKeys)
	eLen := len(eMapKeys)

	if aLen != eLen {
		return DiffValErr{
			dl: dl,
			msg: fmt.Sprintf("map lengths differ. Actual: %d, expected: %d",
				aLen, eLen),
		}
	}

	for i := range aLen {
		k := aMapKeys[i]

		err := diffVals(actVal.MapIndex(k), expVal.MapIndex(k), addKey(dl, k))
		if err != nil {
			return err
		}
	}

	return nil
}

// diffValsStruct returns an error if the two struct values differ
func diffValsStruct(actVal, expVal reflect.Value, dl deepLoc) error {
	// we know that both the expected and actual values have the same number
	// of fields as their types are the same
	fields := actVal.NumField()

	for i := range fields {
		fldName := actVal.Type().Field(i).Name

		err := diffVals(actVal.Field(i), expVal.Field(i), addName(dl, fldName))
		if err != nil {
			return err
		}
	}

	return nil
}

// diffValsSlice returns an error if the two slice values differ
func diffValsSlice(actVal, expVal reflect.Value, dl deepLoc) error {
	aLen := actVal.Len()
	eLen := expVal.Len()

	if aLen != eLen {
		return DiffValErr{
			dl: dl,
			msg: fmt.Sprintf("slice lengths differ. Actual: %d, expected: %d",
				aLen, eLen),
		}
	}

	for i := range aLen {
		err := diffVals(actVal.Index(i), expVal.Index(i), addIdx(dl, i))
		if err != nil {
			return err
		}
	}

	return nil
}

// diffValsArray returns an error if the two array values differ
func diffValsArray(actVal, expVal reflect.Value, dl deepLoc) error {
	// we know both the expected and actual arrays have the same length
	// as their types are the same
	aLen := actVal.Len()

	for i := range aLen {
		err := diffVals(actVal.Index(i), expVal.Index(i), addIdx(dl, i))
		if err != nil {
			return err
		}
	}

	return nil
}

// diffValsPointer returns an error if the two pointer values differ
func diffValsPointer(actVal, expVal reflect.Value, dl deepLoc) error {
	if actVal.Pointer() != expVal.Pointer() {
		return DiffValErr{
			dl:  dl,
			msg: "pointers differ",
		}
	}

	return nil
}

// diffValsUintptr returns an error if the two uintptr values differ
func diffValsUintptr(actVal, expVal reflect.Value, dl deepLoc) error {
	if actVal.Uint() != expVal.Uint() {
		return DiffValErr{
			dl:  dl,
			msg: "uintptr pointers differ",
		}
	}

	return nil
}

// diffValsFunc returns an error if the two func values differ
func diffValsFunc(actVal, expVal reflect.Value, dl deepLoc) error {
	if actVal.Pointer() != expVal.Pointer() {
		return DiffValErr{
			dl:  dl,
			msg: "funcs differ. Actual instance is not equal to expected",
		}
	}

	return nil
}

// diffValsChan returns an error if the two chan values differ
func diffValsChan(actVal, expVal reflect.Value, dl deepLoc) error {
	if actVal.Pointer() != expVal.Pointer() {
		return DiffValErr{
			dl:  dl,
			msg: "chans differ. Actual instance is not equal to expected",
		}
	}

	return nil
}

// diffValsBool returns an error if the two bool values differ
func diffValsBool(actVal, expVal reflect.Value, dl deepLoc) error {
	if actVal.Bool() != expVal.Bool() {
		return DiffValErr{
			dl: dl,
			msg: fmt.Sprintf("bool values differ. Actual: %t, expected: %t",
				actVal.Bool(), expVal.Bool()),
		}
	}

	return nil
}

// diffValsString returns an error if the two string values differ
func diffValsString(actVal, expVal reflect.Value, dl deepLoc) error {
	if actVal.String() != expVal.String() {
		return DiffValErr{
			dl: dl,
			msg: fmt.Sprintf("strings differ. Actual: %q, expected: %q",
				actVal.String(), expVal.String()),
		}
	}

	return nil
}

// diffValsInt returns an error if the two int values differ
func diffValsInt(actVal, expVal reflect.Value, dl deepLoc) error {
	if actVal.Int() != expVal.Int() {
		return DiffValErr{
			dl: dl,
			msg: fmt.Sprintf("int values differ. Actual: %d, expected: %d",
				actVal.Int(), expVal.Int()),
		}
	}

	return nil
}

// diffValsUint returns an error if the two uint values differ
func diffValsUint(actVal, expVal reflect.Value, dl deepLoc) error {
	if actVal.Uint() != expVal.Uint() {
		return DiffValErr{
			dl: dl,
			msg: fmt.Sprintf("uint values differ. Actual: %d, expected: %d",
				actVal.Uint(), expVal.Uint()),
		}
	}

	return nil
}

// diffValsFloat returns an error if the two float values differ
func diffValsFloat(actVal, expVal reflect.Value, dl deepLoc) error {
	if actVal.Float() != expVal.Float() {
		return DiffValErr{
			dl: dl,
			msg: fmt.Sprintf("float values differ. Actual: %g, expected: %g",
				actVal.Float(), expVal.Float()),
		}
	}

	return nil
}

// diffValsComplex returns an error if the two complex values differ
func diffValsComplex(actVal, expVal reflect.Value, dl deepLoc) error {
	if actVal.Complex() != expVal.Complex() {
		return DiffValErr{
			dl: dl,
			msg: fmt.Sprintf("complex values differ. Actual: %g, expected: %g",
				actVal.Complex(), expVal.Complex()),
		}
	}

	return nil
}
