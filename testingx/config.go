package testingx

import (
	"fmt"
	"github.com/ahawker/stdlibx-go/stdlibx"
	"math/rand"
	"os"
	"strconv"
	"testing/quick"
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

// TestPrecondition is a func called prior to test execution to determine
// if the test should be skipped and the reason for it.
type TestPrecondition func() (bool, string)

// Logf is the func called when a boolean condition is not met
// as part of an assertion. It's responsible for generating output
// to a user indicating the test failure and either allowing the
// test to continue or immediately stop execution.
type Logf func(format string, args ...any)

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

// testPreconditionEnvVarSet returns true if the given environment variable name is present.
func testPreconditionEnvVarSet(name string) TestPrecondition {
	return func() (bool, string) {
		if os.Getenv(name) == "" {
			return false, fmt.Sprintf("%q environment variable required for this test type.", name)
		}
		return true, ""
	}
}

// NewTestConfig creates a new *TestConfig for the given functional opts
// and sane defaults.
func NewTestConfig(options ...stdlibx.Option[*TestConfig]) (*TestConfig, error) {
	config := &TestConfig{
		Logf:         defaultTestConfig.Logf,
		Parallel:     defaultTestConfig.Parallel,
		Precondition: defaultTestConfig.Precondition,
		QuickConfig: &quick.Config{
			MaxCount: defaultTestConfig.QuickConfig.MaxCount,
			Rand:     defaultTestConfig.QuickConfig.Rand,
		},
	}
	return stdlibx.OptionApply(config, options...)
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
func WithTestLogf(log Logf) stdlibx.Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.Logf = log
		return nil
	}
}

// WithTestParallel sets the config parallel.
func WithTestParallel(parallel bool) stdlibx.Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.Parallel = parallel
		return nil
	}
}

// WithTestPrecondition sets the config precondition.
func WithTestPrecondition(precondition TestPrecondition) stdlibx.Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.Precondition = precondition
		return nil
	}
}

// WithTestMaxCount sets the config property testing max count.
func WithTestMaxCount(maxCount int) stdlibx.Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.QuickConfig.MaxCount = maxCount
		return nil
	}
}

// WithTestMaxCountScale sets the config property testing max count scale.
func WithTestMaxCountScale(maxCountScale float64) stdlibx.Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.QuickConfig.MaxCountScale = maxCountScale
		return nil
	}
}

// WithTestRand sets the config rand number generator.
func WithTestRand(r *rand.Rand) stdlibx.Option[*TestConfig] {
	return func(t *TestConfig) error {
		t.QuickConfig.Rand = r
		return nil
	}
}
