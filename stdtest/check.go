package stdtest

import (
	"testing"
	"testing/quick"
)

var _ Checker = (*Check)(nil)

// Checker defines common quick tests.
type Checker interface {
	Check(fn any) bool
	CheckEqual(fn1, fn2 any) bool
}

// Check implements property testing via the "testing/quick" package.
type Check struct {
	// tb is the active test/benchmark test being executed.
	tb testing.TB
	// Logf is called when condition of test assertion is not met.
	logf Logf
	// config modifies how quick test (property tests) are executed.
	config *quick.Config
}

// Check fails the test if the check function returns false.
func (c *Check) Check(fn any) bool {
	c.tb.Helper()
	if err := quick.Check(fn, c.config); err != nil {
		c.logf("\n\n\tgot: %v; \n\n\twant: no error\n", err)
		return false
	}
	return true
}

// CheckEqual fails the test if the check function returns false.
func (c *Check) CheckEqual(fn1, fn2 any) bool {
	c.tb.Helper()
	if err := quick.CheckEqual(fn1, fn2, c.config); err != nil {
		c.logf("\n\n\tgot: %v; \n\n\twant: no error\n", err)
		return false
	}
	return true
}
