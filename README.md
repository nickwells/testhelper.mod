[![GoDoc](https://godoc.org/github.com/nickwells/testhelper.mod?status.png)](https://godoc.org/github.com/nickwells/testhelper.mod)

# testhelper.mod
Some useful functions and types to simplify and improve testing.

Several of these are intended to be used as unnamed members of a testcase
struct. The testcase idiom is a table-driven testing scheme where you first
define a slice of structs with the struct holding the parameters to the test
and the expected results. Then you iterate over the individual test cases
running the code you are testing and compare the results with the expected
results.

## the testhelper.ID type
This is intended to be used as an unnamed member of a testcase struct. Then
for each testcase instance you initialise the value using the constructor:

```go
testhelper.MkID("...")
```

This has the advantage over a simple string that the constructor records the
filename and linenumber where MkID(...)  was called which makes it easier to
find the precise test that is failing. It can be used to give a useful
description of the test case through the IDStr() method.

A struct with this embedded will satisfy the testhelper.TestCase interface.

## the testhelper.ExpErr type
This is intended to be used as an unnamed member of a testcase struct (though
if you want to check more than one error condition you can add more). It is
initialised using the constructor:

```go
testhelper.MkErr("part of the error message", "some more", "etc")
```

The default value expresses that no error is expected. If it is initialised a
non-nil error is expected and the strings passed are expected to be found in
the error message.

A struct with this embedded will satisfy the testhelper.TestErr interface and
if the struct also has a testhelper.ID embedded then it will satisfy the
testhelper.TestCaseWithErr interface. This can then be passed to
testhelper.CheckExpErr which will report a test error if the error is not as
expected.

## the testhelper.ExpPanic type
This is intended to be used as an unnamed member of a testcase struct (though
if you want to check more than one panic you can add more). It is initialised
using the constructor:

```go
testhelper.MkExpPanic("part of the panic message", "some more", "etc")
```

The default value expresses that no panic is expected. If it is initialised a
panic is expected and the strings passed are expected to be found in
the panic value (which is expected to be a string).

A struct with this embedded will satisfy the testhelper.TestPanic interface
and if the struct also has a testhelper.ID embedded then it will satisfy the
testhelper.TestCaseWithPanic interface. This can then be passed to
testhelper.CheckExpPanic which will report a test error if the panic is not
as expected.

## the StringSliceDiff func
This takes a pair of string slices and returns true if they differ, false
otherwise.

## the ShouldContain func
This takes a string and a slice of strings and reports test errors for each
string in the slice that isn't in the string. If any strings are missing it
will return true, otherwise false.

## the CheckAgainstGoldenFile func
This checks that the passed slice of bytes is the same as the value read from
the golden file. You can pass a flag to get the file created initially and to
update it when there has been a desired change in the output being
checked. It returns true if the passed bytes match the contents of the golden
file, false otherwise.
