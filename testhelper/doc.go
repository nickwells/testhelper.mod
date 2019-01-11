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

    ...
    for i, tc := range testcases {
        tcID := fmt.Sprintf("test %d: %s", i, tc.name)
        ...
        testhelper.SomeFunc(t, tcID, ...)
        ...
    }

*/
package testhelper
