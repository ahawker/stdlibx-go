package stdtest

import (
	"github.com/ahawker/stdlibx-go/stdlib"
	"testing"
)

var (
	_ Asserter = (*Test)(nil)
	_ Checker  = (*Test)(nil)
)

// BenchmarkTest creates a new *Benchmark configured for running only "benchmark" tests
// when the BENCHMARK_TEST environment variable is set.
func BenchmarkTest(t *testing.B, options ...stdlib.Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []stdlib.Option[*TestConfig]{
		WithTestPrecondition(testPreconditionEnvVarSet("BENCHMARK_TEST")),
	}
	return newTest(t, append(defaults, options...)...)
}

// FuzzTest creates a new *Benchmark configured for running only "fuzz" tests
// when the FUZZ_TEST environment variable is set.
func FuzzTest(t *testing.F, options ...stdlib.Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []stdlib.Option[*TestConfig]{
		WithTestLogf(t.Errorf),
		WithTestPrecondition(testPreconditionEnvVarSet("FUZZ_TEST")),
	}
	return newTest(t, append(defaults, options...)...)
}

// IntegrationTest creates a new *Test configured for running only "integration" tests
// when the INTEGRATION_TEST environment variable is set.
func IntegrationTest(t *testing.T, options ...stdlib.Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []stdlib.Option[*TestConfig]{
		WithTestPrecondition(testPreconditionEnvVarSet("INTEGRATION_TEST")),
	}
	return newTest(t, append(defaults, options...)...)
}

// PropertyTest creates a new *Test configured for running only "property" tests.
func PropertyTest(t testing.TB, options ...stdlib.Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []stdlib.Option[*TestConfig]{
		WithTestPrecondition(testPreconditionEnvVarSet("PROPERTY_TEST")),
	}
	return newTest(t, append(defaults, options...)...)
}

// UnitTest creates a new *Test configured for running only "unit" tests.
func UnitTest(t testing.TB, options ...stdlib.Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []stdlib.Option[*TestConfig]{
		WithTestPrecondition(testPreconditionEnvVarSet("UNIT_TEST")),
	}
	return newTest(t, append(defaults, options...)...)
}

// newTest creates a new *Test with the given options.
func newTest(t testing.TB, options ...stdlib.Option[*TestConfig]) *Test {
	t.Helper()

	defaults := []stdlib.Option[*TestConfig]{
		WithTestLogf(t.Fatalf),
	}
	config, err := NewTestConfig(append(defaults, options...)...)
	if err != nil {
		t.Fatal(err)
	}

	// Skip tests where precondition predicate is false.
	if ok, reason := config.Precondition(); !ok {
		t.Skip(reason)
	}

	// Parallelize test execution (multiple go-routines) when possible.
	if config.Parallel {
		if p, ok := t.(interface{ Parallel() }); ok {
			p.Parallel()
		}
	}

	// Set process environment variables for the duration of the test.
	// Note: This does not work with parallel tests.
	for k, v := range config.Env {
		t.Setenv(k, v)
	}

	return &Test{
		TB: t,
		Asserter: &Assert{
			tb:   t,
			logf: config.Logf,
		},
		Checker: &Check{
			tb:     t,
			logf:   config.Logf,
			config: config.QuickConfig,
		},
	}
}

// Test merges the stdlib 'testing.TB' with our test helpers
// into a single implementation.
//
// By using this 'Test' instead of the `testing` structs (T, B, F) we can automatically handle
// many common assertion patterns and preconditions for different test types (functional,
// fuzz, or integration).
type Test struct {
	testing.TB
	Asserter
	Checker
}

// Sub runs the given function as a subtest of the current test.
func (t *Test) Sub(
	name string,
	fn func(subtest *Test),
	options ...stdlib.Option[*TestConfig],
) bool {
	switch tb := t.TB.(type) {
	case *testing.T:
		return tb.Run(name, func(st *testing.T) {
			fn(newTest(st, options...))
		})
	case *testing.B:
		return tb.Run(name, func(bt *testing.B) {
			fn(newTest(bt, options...))
		})
	case *testing.F:
		t.Fatalf("subtest not support for fuzz tests name=%v", t.TB.Name())
		return false
	default:
		t.Fatalf("subtest not support for unknown test type name=%s type=%T", t.TB.Name(), t.TB)
		return false
	}
}
