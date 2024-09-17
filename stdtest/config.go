package stdtest

import (
	"context"
	"fmt"
	"github.com/ahawker/stdlibx-go/stdlib"
	"math/rand"
	"os"
	"strconv"
	"testing/quick"
	"time"
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

// ErrTestParallelWithSetEnv is returned when attempting to create
// a parallel test that overrides process environment variables.
var ErrTestParallelWithSetEnv = stdlib.Error{
	Code:      "parallel_with_setenv",
	Message:   "test cannot be parallel and also modify environment variables",
	Namespace: stdlib.ErrorNamespaceDefault,
}

// TestPrecondition is a func called prior to test execution to determine
// if the test should be skipped and the reason for it.
type TestPrecondition func() (bool, string)

// Logf is the func called when a boolean condition is not met
// as part of an assertion. It's responsible for generating output
// to a user indicating the test failure and either allowing the
// test to continue or immediately stop execution.
type Logf func(format string, args ...any)

// TestPreconditionEnvVarSet returns true if the given environment variable name is present.
func TestPreconditionEnvVarSet(name string) TestPrecondition {
	return func() (bool, string) {
		if os.Getenv(name) == "" {
			return false, fmt.Sprintf("%q env var required for this test type", name)
		}
		return true, ""
	}
}

// testPreconditionNoOp is a dummy test precondition that is always true.
var testPreconditionNoOp = func() (bool, string) {
	return true, ""
}

// defaultTestConfig contains default values for test configuration.
var defaultTestConfig = &TestConfig{
	// Context is the default context for an individual test.
	Context: context.Background(),
	// Logf is the default func called when condition of test assertion is not met.
	Logf: func(format string, args ...any) {
		_, _ = fmt.Fprintf(os.Stderr, format, args...)
	},
	// Parallel is the default value for attempting to run tests in parallel.
	// It can be overridden by setting the "STDLIBX_TEST_DEFAULT_PARALLEL" environment variable.
	Parallel: false,
	// Preconditions being always true is the default before running each test.
	Preconditions: []TestPrecondition{testPreconditionNoOp},
	// QuickConfig is the default config for property testing using "testing/quick".
	QuickConfig: &quick.Config{},
	// Skip is the default reason for skipping a test.
	Skip: "",
	// Timeout is the default duration for an individual test.
	Timeout: 5 * time.Minute,
}

// NewTestConfig creates a new *TestConfig for the given functional opts
// and sane defaults.
func NewTestConfig(options ...stdlib.Option[*TestConfig]) (*TestConfig, error) {
	config := &TestConfig{
		Context:       defaultTestConfig.Context,
		Logf:          defaultTestConfig.Logf,
		Parallel:      defaultTestConfig.Parallel,
		Preconditions: defaultTestConfig.Preconditions,
		QuickConfig: &quick.Config{
			MaxCount: defaultTestConfig.QuickConfig.MaxCount,
			Rand:     defaultTestConfig.QuickConfig.Rand,
		},
		Skip:    defaultTestConfig.Skip,
		Timeout: defaultTestConfig.Timeout,
	}
	return stdlib.OptionApply(config, options...)
}

// TestConfig defines config options for test.
type TestConfig struct {
	// Context is context for an individual test.
	Context context.Context
	// Env contains key/value pairs to be set with 'os.SetEnv' for the scope of the test.
	Env map[string]string
	// Logf is called when condition of test assertion is not met.
	Logf Logf
	// Parallel enables/disables Parallel test execution.
	Parallel bool
	// Preconditions are called before test execution to determine if test should
	// be skipped and the reason for it.
	Preconditions []TestPrecondition
	// QuickConfig modifies of how quick test (property tests) are executed.
	QuickConfig *quick.Config
	// Skip test for this reason.
	Skip string
	// Timeout is timeout duration for test execution.
	Timeout time.Duration
}

// WithTestContext sets the config ctx.
func WithTestContext(ctx context.Context) stdlib.Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.Context = ctx
		return nil
	}
}

// WithTestEnv sets the config env.
func WithTestEnv(env map[string]string) stdlib.Option[*TestConfig] {
	return func(t *TestConfig) error {
		if t.Parallel && len(env) > 0 {
			return ErrTestParallelWithSetEnv
		}
		t.Env = env
		return nil
	}
}

// WithTestLogf sets the config logf.
func WithTestLogf(log Logf) stdlib.Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.Logf = log
		return nil
	}
}

// WithTestParallel sets the config parallel.
func WithTestParallel(parallel bool) stdlib.Option[*TestConfig] {
	return func(t *TestConfig) error {
		if parallel && len(t.Env) > 0 {
			return ErrTestParallelWithSetEnv
		}
		t.Parallel = parallel
		return nil
	}
}

// WithTestPrecondition sets the config preconditions.
func WithTestPrecondition(preconditions ...TestPrecondition) stdlib.Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.Preconditions = append(t.Preconditions, preconditions...)
		return nil
	}
}

// WithTestMaxCount sets the config property testing max count.
func WithTestMaxCount(maxCount int) stdlib.Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.QuickConfig.MaxCount = maxCount
		return nil
	}
}

// WithTestMaxCountScale sets the config property testing max count scale.
func WithTestMaxCountScale(maxCountScale float64) stdlib.Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.QuickConfig.MaxCountScale = maxCountScale
		return nil
	}
}

// WithTestRand sets the config rand number generator.
func WithTestRand(r *rand.Rand) stdlib.Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.QuickConfig.Rand = r
		return nil
	}
}

// WithTestSkip sets the config skip reason.
func WithTestSkip(skip string) stdlib.Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.Skip = skip
		return nil
	}
}

// WithTestTimeout sets the config timeout.
func WithTestTimeout(to time.Duration) stdlib.Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.Timeout = to
		return nil
	}
}
