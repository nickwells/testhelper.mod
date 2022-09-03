/*
Package testhelper contains helper functions for doing things you might want
to do in lots of your tests. For instance checking that an error value or a
panic status was as expected.

The functions where a 't *testing.T' is passed will mark themselves as
helpers by calling t.Helper(). These will also report the error themselves.

The testID parameter that many of the helpers take should be a string which
can identify the instance of the test that is being run. For instance, if you
have a slice of test cases you might want this string to include the index of
the test and possibly some descriptive name. It might be created and used as
follows:

	testcases := []struct{
	    name string
	    ...
	}{
	    {
	        name: "whatever",
	    },
	}
	...
	for i, tc := range testcases {
	    tcID := fmt.Sprintf("test %d: %s", i, tc.name)
	    ...
	    testhelper.SomeFunc(t, tcID, ...)
	    ...
	}

Alternatively you can construct the testcase struct with an embedded
testhelper.ID member. Then the testcase ID can be initialised with the MkID
func. The ID string can then be created with the IDStr() func defined on the
testhelper.ID as follows:

	testcases := []struct{
	    testhelper.ID
	    ...
	}{
	    {
	        ID: testhelper.MkID("whatever"),
	    },
	}
	...
	for _, tc := range testcases {
	    tcID := tc.IDStr()
	    ...
	    testhelper.SomeFunc(t, tcID, ...)
	    ...
	}

This way of constructing a test case struct has the advantage that the
constructed ID string gives more information about where the test was
constructed and several of the testhelper functions take a testhelper.ID (or
an interface which it satisfies).

The problem which the testhelper.ID solves is that reporting just the index
number of a test (as in the previous example) imposes work on the
tester such as counting tests to reach the failing test case. Using
testhelper.ID gives you the filename and line number where the ID was
instantiated.

There are additional mixin structs beside ID which allow you to record such
things as whether an error is expected and if so what the error string should
contain.
*/
package testhelper
