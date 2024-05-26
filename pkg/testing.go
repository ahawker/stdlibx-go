package stdlibx

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"testing/quick"
)

var (
	_ TestAsserter = (*Test)(nil)
	_ TestChecker  = (*Test)(nil)
)

func init() {
	// Global test parallel override via environment variable.
	if v, ok := os.LookupEnv("STDLIBX_TEST_DEFAULT_PARALLEL"); ok {
		b, err := strconv.ParseBool(v)
		if err == nil {
			defaultTestConfig.Parallel = b
		}
	}
	// Global override for testing/quick max count.
	iterations, err := strconv.ParseInt(os.Getenv("STDLIBX_TEST_PROPERTY_MAX_COUNT"), 10, 32)
	if err == nil {
		defaultTestConfig.QuickConfig.MaxCount = int(iterations)
	}
	// Global override for testing/quick test seed.
	seed, err := strconv.ParseInt(os.Getenv("STDLIBX_TEST_RANDOM_SEED"), 10, 64)
	if err == nil {
		defaultTestConfig.QuickConfig.Rand = rand.New(rand.NewSource(seed)) //nolint:gosec
	}
}

// defaultTestConfig contains default values for test configuration.
var defaultTestConfig = &TestConfig{
	// Logf is the default func called when condition of test assertion is not met.
	Logf: func(format string, args ...any) {
		_, _ = fmt.Fprintf(os.Stderr, format, args...)
	},
	// Parallel is the default value for attempting to run tests in parallel.
	// It can be overridden by setting the "STDLIBX_TEST_DEFAULT_PARALLEL" environment variable.
	Parallel: false,
	// Precondition is the default precondition check before running each test.
	Precondition: func() (bool, string) {
		return true, ""
	},
	// QuickConfig is the default config for property testing using "testing/quick".
	QuickConfig: &quick.Config{},
}

// TestPrecondition is a func called prior to test execution to determine
// if the test should be skipped and the reason for it.
type TestPrecondition func() (bool, string)

// TestAsserter defines common test assertions.
type TestAsserter interface {
	True(condition bool, format string, args ...any) bool
	False(condition bool, format string, args ...any) bool
	OK(err error) bool
	NotOK(err error) bool
	Match(got, want string) bool
	Equal(got, want any) bool
	PointerEqual(got, want any) bool
	ErrorEqual(got, want error) bool
	Panic(got func()) bool
}

// TestChecker defines common quick tests.
type TestChecker interface {
	Check(fn any) bool
	CheckEqual(fn1, fn2 any) bool
}

// Logf is the func called when a boolean condition is not met
// as part of an assertion. It's responsible for generating output
// to a user indicating the test failure and either allowing the
// test to continue or immediately stop execution.
type Logf func(format string, args ...any)

// TestAssert implements helpers for common assertion patterns.
type TestAssert struct {
	// tb is the active test/benchmark test being executed.
	tb testing.TB
	// Logf is called when condition of test assertion is not met.
	logf Logf
}

// True fails the test if the condition is false.
func (a *TestAssert) True(condition bool, format string, args ...any) bool {
	a.tb.Helper()
	if !condition {
		msg := fmt.Sprintf(format, args...)
		a.logf("\n\n\tgot:  false\n\n\twant: true\n\n\tmsg: %s\n", msg)
		return false
	}
	return true
}

// False fails the test if the condition is true.
func (a *TestAssert) False(condition bool, format string, args ...any) bool {
	a.tb.Helper()
	if condition {
		msg := fmt.Sprintf(format, args...)
		a.logf("\n\n\tgot:  true\n\n\twant: false\n\n\tmsg: %s\n", msg)
		return false
	}
	return true
}

// OK fails the test if err is not nil.
func (a *TestAssert) OK(err error) bool {
	a.tb.Helper()
	if err != nil {
		a.logf("\n\n\tgot:  %v\n\n\twant: no error\n", err)
		return false
	}
	return true
}

// NotOK fails the test if err is nil.
func (a *TestAssert) NotOK(err error) bool {
	a.tb.Helper()
	if err == nil {
		a.logf("\n\n\tgot:  nil\n\n\twant: error\n")
		return false
	}
	return true
}

// Match fails the test if got does not match want pattern.
func (a *TestAssert) Match(got, wantPattern string) bool {
	a.tb.Helper()
	compiled, err := regexp.Compile(wantPattern)
	if err != nil {
		a.logf("\n\n\tmatch want pattern compile error: %v\n")
		return false
	}
	if !compiled.MatchString(got) {
		a.logf("\n\n\tgot:  %#v\n\n\twant pattern match: %s\n", got, wantPattern)
		return false
	}
	return true
}

// Equal fails the test if got is not equal to want using reflect.DeepEqual.
func (a *TestAssert) Equal(got, want any) bool {
	a.tb.Helper()
	if !reflect.DeepEqual(got, want) {
		a.logf("\n\n\tgot:  %#v\n\n\twant: %#v\n", got, want)
		return false
	}
	return true
}

// PointerEqual fails the test if got is not equal to want for pointers.
func (a *TestAssert) PointerEqual(got, want any) bool {
	a.tb.Helper()
	gotP := reflect.ValueOf(got).Pointer()
	wantP := reflect.ValueOf(want).Pointer()
	if gotP != wantP {
		a.logf("\n\n\tgot:  %#v\n\n\twant: %#v\n", gotP, wantP)
		return false
	}
	return true
}

// ErrorEqual fails the test if got is not equal to want for errors.
func (a *TestAssert) ErrorEqual(got, want error) bool {
	a.tb.Helper()
	if !errors.Is(want, got) {
		a.logf("\n\n\tgot:  %#v\n\n\twant: %#v\n", got, want)
		return false
	}
	return true
}

// Panic fails the test if the got function does not panic.
func (a *TestAssert) Panic(got func()) bool {
	a.tb.Helper()
	defer func() {
		_ = recover()
	}()
	got()
	a.logf("\n\n\tgot:  no panic\n\n\twant: panic\n")
	return false
}

// TestCheck implements property testing via the "testing/quick" package.
type TestCheck struct {
	// tb is the active test/benchmark test being executed.
	tb testing.TB
	// Logf is called when condition of test assertion is not met.
	logf Logf
	// config modifies how quick test (property tests) are executed.
	config *quick.Config
}

// Check fails the test if the check function returns false.
func (c *TestCheck) Check(fn any) bool {
	c.tb.Helper()

	if err := quick.Check(fn, c.config); err != nil {
		c.logf("\n\n\tgot: %v; \n\n\twant: no error\n", err)
		return false
	}
	return true
}

// CheckEqual fails the test if the check function returns false.
func (c *TestCheck) CheckEqual(fn1, fn2 any) bool {
	c.tb.Helper()

	if err := quick.CheckEqual(fn1, fn2, c.config); err != nil {
		c.logf("\n\n\tgot: %v; \n\n\twant: no error\n", err)
		return false
	}
	return true
}

// NewTestConfig creates a new *TestConfig for the given functional opts
// and sane defaults.
func NewTestConfig(options ...Option[*TestConfig]) (*TestConfig, error) {
	config := &TestConfig{
		Logf:         defaultTestConfig.Logf,
		Parallel:     defaultTestConfig.Parallel,
		Precondition: defaultTestConfig.Precondition,
		QuickConfig: &quick.Config{
			MaxCount: defaultTestConfig.QuickConfig.MaxCount,
			Rand:     defaultTestConfig.QuickConfig.Rand,
		},
	}
	return OptionApply(config, options...)
}

// TestConfig defines config options for test.
type TestConfig struct {
	// Logf is called when condition of test assertion is not met.
	Logf Logf
	// Parallel enables/disables Parallel test execution.
	Parallel bool
	// Precondition is called before test execution to determine if test should
	// be skipped and the reason for it.
	Precondition TestPrecondition
	// QuickConfig modifies of how quick test (property tests) are executed.
	QuickConfig *quick.Config
}

// WithTestLogf sets the config logf.
func WithTestLogf(log Logf) Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.Logf = log
		return nil
	}
}

// WithTestParallel sets the config parallel.
func WithTestParallel(parallel bool) Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.Parallel = parallel
		return nil
	}
}

// WithTestPrecondition sets the config precondition.
func WithTestPrecondition(precondition TestPrecondition) Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.Precondition = precondition
		return nil
	}
}

// WithTestMaxCount sets the config property testing max count.
func WithTestMaxCount(maxCount int) Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.QuickConfig.MaxCount = maxCount
		return nil
	}
}

// WithTestMaxCountScale sets the config property testing max count scale.
func WithTestMaxCountScale(maxCountScale float64) Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.QuickConfig.MaxCountScale = maxCountScale
		return nil
	}
}

// WithTestRand sets the config rand number generator.
func WithTestRand(r *rand.Rand) Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.QuickConfig.Rand = r
		return nil
	}
}

// testPreconditionEnvVarSet returns true if the given environment variable name is present.
func testPreconditionEnvVarSet(name string) TestPrecondition {
	return func() (bool, string) {
		if os.Getenv(name) == "" {
			return false, fmt.Sprintf("%q environment variable required for this test type.", name)
		}
		return true, ""
	}
}

// BenchmarkTest creates a new *Benchmark configured for running only "benchmark" tests
// when the BENCHMARK_TEST environment variable is set.
func BenchmarkTest(t *testing.B, options ...Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []Option[*TestConfig]{
		WithTestPrecondition(testPreconditionEnvVarSet("BENCHMARK_TEST")),
	}
	return newTest(t, append(defaults, options...)...)
}

// FuzzTest creates a new *Benchmark configured for running only "fuzz" tests
// when the FUZZ_TEST environment variable is set.
func FuzzTest(t *testing.F, options ...Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []Option[*TestConfig]{
		WithTestLogf(t.Errorf),
		WithTestPrecondition(testPreconditionEnvVarSet("FUZZ_TEST")),
	}
	return newTest(t, append(defaults, options...)...)
}

// IntegrationTest creates a new *Test configured for running only "integration" tests
// when the INTEGRATION_TEST environment variable is set.
func IntegrationTest(t *testing.T, options ...Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []Option[*TestConfig]{
		WithTestPrecondition(testPreconditionEnvVarSet("INTEGRATION_TEST")),
	}
	return newTest(t, append(defaults, options...)...)
}

// PropertyTest creates a new *Test configured for running only "property" tests.
func PropertyTest(t testing.TB, options ...Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []Option[*TestConfig]{
		WithTestPrecondition(testPreconditionEnvVarSet("PROPERTY_TEST")),
	}
	return newTest(t, append(defaults, options...)...)
}

// UnitTest creates a new *Test configured for running only "unit" tests.
func UnitTest(t testing.TB, options ...Option[*TestConfig]) *Test {
	t.Helper()
	defaults := []Option[*TestConfig]{
		WithTestPrecondition(testPreconditionEnvVarSet("UNIT_TEST")),
	}
	return newTest(t, append(defaults, options...)...)
}

// newTest creates a new *Test with the given options.
func newTest(t testing.TB, options ...Option[*TestConfig]) *Test {
	t.Helper()

	defaults := []Option[*TestConfig]{
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

	return &Test{
		TB: t,
		TestAsserter: &TestAssert{
			tb:   t,
			logf: config.Logf,
		},
		TestChecker: &TestCheck{
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
	TestAsserter
	TestChecker
}

// Sub runs the given function as a subtest of the current test.
func (t *Test) Sub(
	name string,
	fn func(subtest *Test),
	options ...Option[*TestConfig],
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

// TestFn is a function that executes a single test for the given testcase.
//
// Note: The benefit of passing in an isolated function for the test is that we can avoid
// issues with closures when t.Run(...) is called in parallel with the testcase.
type TestFn[TGot any, TWant any] func(t *Test, tc TestCase[TGot, TWant])

// TestCase encapsulates the inputs and expectations of a single testcase.
type TestCase[TGot any, TWant any] struct {
	// Got stores inputs for the testcase.
	Got TGot
	// Want stores the expectations for the testcase.
	Want TWant
	// WantErr stores the optional error expectation for the testcase.
	WantErr error
	// Options are custom options for test execution specific to this case.
	Options []Option[*TestConfig]
}

// TestTable represents a collection of "table tests" in the form of
// named test cases.
type TestTable[TGot any, TWant any] map[string]TestCase[TGot, TWant]

// TestTableRun will execute all testcases as unit tests with the given test function.
func TestTableRun[TGot any, TWant any](
	t testing.TB,
	fn TestFn[TGot, TWant],
	table TestTable[TGot, TWant],
	options ...Option[*TestConfig],
) {
	t.Helper()
	test := newTest(t, options...)
	for name, testcase := range table {
		test.Sub(name, func(subtest *Test) {
			fn(subtest, testcase)
		},
			options...,
		)
	}
}
