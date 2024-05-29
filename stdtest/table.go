package stdtest

import (
	"github.com/ahawker/stdlibx-go/stdlib"
	"testing"
)

// TestFunc is a function that executes a single test for the given testcase.
//
// Note: The benefit of passing in an isolated function for the test is that we can avoid
// issues with closures when t.Run(...) is called in parallel with the testcase.
type TestFunc[TGot any, TWant any] func(t *Test, tc Testcase[TGot, TWant])

// Testcase encapsulates the inputs and expectations of a single testcase.
type Testcase[TGot any, TWant any] struct {
	// Got stores inputs for the testcase.
	Got TGot
	// Want stores the expectations for the testcase.
	Want TWant
	// WantErr stores the optional error expectation for the testcase.
	WantErr error
	// Options are custom options for test execution specific to this case.
	Options []stdlib.Option[*TestConfig]
}

// Table represents a collection of "table tests" in the form of
// named test cases.
type Table[TGot any, TWant any] map[string]Testcase[TGot, TWant]

// Run will execute all table tests using the given test function.
func (table Table[TGot, TWant]) Run(
	tb testing.TB,
	testFn TestFunc[TGot, TWant],
	options ...stdlib.Option[*TestConfig],
) {
	tb.Helper()
	test := newTest(tb, options...)
	for name, testcase := range table {
		subtestFn := func(subtest *Test) {
			testFn(subtest, testcase)
		}
		test.Sub(name, subtestFn, append(options, testcase.Options...)...)
	}
}
