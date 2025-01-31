package stdtest

import (
	"github.com/ahawker/stdlibx-go/stdlib"
	"io"
	"testing"
)

var (
	_ testing.TB = (*Test)(nil)
	_ Asserter   = (*Test)(nil)
	_ Checker    = (*Test)(nil)
)

// BenchmarkTest creates a new *Benchmark configured for running only "benchmark" tests
// when the BENCHMARK_TEST environment variable is set.
func BenchmarkTest(t *testing.B, options ...stdlib.Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []stdlib.Option[*TestConfig]{
		WithTestPrecondition(TestPreconditionEnvVarSet("BENCHMARK_TEST")),
	}
	return NewTest(t, append(defaults, options...)...)
}

// FuzzTest creates a new *Benchmark configured for running only "fuzz" tests
// when the FUZZ_TEST environment variable is set.
func FuzzTest(t *testing.F, options ...stdlib.Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []stdlib.Option[*TestConfig]{
		WithTestLogf(t.Errorf),
		WithTestPrecondition(TestPreconditionEnvVarSet("FUZZ_TEST")),
	}
	return NewTest(t, append(defaults, options...)...)
}

// IntegrationTest creates a new *Test configured for running only "integration" tests
// when the INTEGRATION_TEST environment variable is set.
func IntegrationTest(t *testing.T, options ...stdlib.Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []stdlib.Option[*TestConfig]{
		WithTestPrecondition(TestPreconditionEnvVarSet("INTEGRATION_TEST")),
	}
	return NewTest(t, append(defaults, options...)...)
}

// PropertyTest creates a new *Test configured for running only "property" tests.
func PropertyTest(t testing.TB, options ...stdlib.Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []stdlib.Option[*TestConfig]{
		WithTestPrecondition(TestPreconditionEnvVarSet("PROPERTY_TEST")),
	}
	return NewTest(t, append(defaults, options...)...)
}

// UnitTest creates a new *Test configured for running only "unit" tests.
func UnitTest(t testing.TB, options ...stdlib.Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []stdlib.Option[*TestConfig]{
		WithTestPrecondition(TestPreconditionEnvVarSet("UNIT_TEST")),
	}
	return NewTest(t, append(defaults, options...)...)
}

// NewTest creates a new *Test with the given options.
func NewTest(t testing.TB, options ...stdlib.Option[*TestConfig]) *Test {
	t.Helper()

	defaults := []stdlib.Option[*TestConfig]{
		WithTestLogf(t.Fatalf),
	}
	config, err := NewTestConfig(append(defaults, options...)...)
	if err != nil {
		t.Fatal(err)
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

	// Skip the test if a reason is set.
	// Commonly used for tests that are not yet implemented.
	if config.Skip != "" {
		t.Skip(config.Skip)
	}

	// Skip tests where precondition predicate is false.
	for _, precondition := range config.Preconditions {
		if ok, reason := precondition(); !ok {
			t.Skip(reason)
		}
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
		Config: config,
	}
}

// Test merges the stdlib 'testing.TB' with our test helpers
// into a single implementation.
//
// By using this 'Test' instead of the `testing` structs (T, B, F) we can automatically handle
// many common assertion patterns and preconditions for different test types (functional,
// fuzz, or integration).
type Test struct {
	// TB is the golang 'testing' implementation for (T, B, F) tests.
	testing.TB
	// Asserter handles common test assertions.
	Asserter
	// Checker handles property tests.
	Checker
	// Config stores configuration specific to an individual test.
	Config *TestConfig
}

// Fuzz runs the given function as a fuzztest of the current test.
func (t *Test) Fuzz(fn any) bool {
	switch tb := t.TB.(type) {
	case *testing.F:
		tb.Fuzz(fn)
		return true
	default:
		t.Fatalf("fuzztest not supported for %T name=%s", t.TB, t.TB.Name())
		return false
	}
}

// Closeup registers a cleanup configuration for the closer.
func (t *Test) Closeup(closer io.Closer) {
	t.Helper()

	if closer == nil {
		return
	}
	t.Cleanup(func() {
		if err := closer.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}
	})
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
			fn(NewTest(st, options...))
		})
	case *testing.B:
		return tb.Run(name, func(bt *testing.B) {
			fn(NewTest(bt, options...))
		})
	case *testing.F:
		t.Fatalf("subtest not supported for *testing.F name=%v", t.TB.Name())
		return false
	default:
		t.Fatalf("subtest not supported for %T name=%s", t.TB, t.TB.Name())
		return false
	}
}
