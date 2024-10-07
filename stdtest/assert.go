package stdtest

import (
	"errors"
	"fmt"
	"github.com/ahawker/stdlibx-go/stdlib"
	"math"
	"reflect"
	"regexp"
	"testing"
)

var _ Asserter = (*Assert)(nil)

// Comparer is a function that compares two values of any type and returns true
// if they are "equal".
type Comparer func(got, want any) bool

// Asserter defines common test assertions.
type Asserter interface {
	True(condition bool, format string, args ...any) bool
	False(condition bool, format string, args ...any) bool
	OK(err error) bool
	NotOK(err error) bool
	Match(got, want string) bool
	NotEqual(got, want any) bool
	Equal(got, want any) bool
	EqualPointer(got, want any) bool
	EqualError(got, want error) bool
	EqualType(got, want any) bool
	EqualComparer(got, want any, cmp Comparer) bool
	EqualFloat(got, want, epsilon float64) bool
	EqualPattern(got any, want string) bool
	EqualRegex(got any, want *regexp.Regexp) bool
	Panic(got func()) bool
}

// Assert implements helpers for common assertion patterns.
type Assert struct {
	// tb is the active test/benchmark test being executed.
	tb testing.TB
	// Logf is called when condition of test assertion is not met.
	logf Logf
}

// True fails the test if the condition is false.
func (a *Assert) True(condition bool, format string, args ...any) bool {
	a.tb.Helper()
	if !condition {
		msg := fmt.Sprintf(format, args...)
		a.logf("\n\n\tgot:  false\n\n\twant: true\n\n\tmsg: %s\n", msg)
		return false
	}
	return true
}

// False fails the test if the condition is true.
func (a *Assert) False(condition bool, format string, args ...any) bool {
	a.tb.Helper()
	if condition {
		msg := fmt.Sprintf(format, args...)
		a.logf("\n\n\tgot:  true\n\n\twant: false\n\n\tmsg: %s\n", msg)
		return false
	}
	return true
}

// OK fails the test if err is not nil.
func (a *Assert) OK(err error) bool {
	a.tb.Helper()
	if err != nil {
		a.logf("\n\n\tgot:  %v\n\n\twant: no error\n", err)
		return false
	}
	return true
}

// NotOK fails the test if err is nil.
func (a *Assert) NotOK(err error) bool {
	a.tb.Helper()
	if err == nil {
		a.logf("\n\n\tgot:  nil\n\n\twant: %v\n", err)
		return false
	}
	return true
}

// Match fails the test if 'got' does not match 'want' pattern.
func (a *Assert) Match(got, wantPattern string) bool {
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

// NotEqual fails the test if 'got' is equal to 'want' using reflect.DeepEqual.
func (a *Assert) NotEqual(got, want any) bool {
	a.tb.Helper()
	if reflect.DeepEqual(got, want) {
		a.logf("\n\n\tgot:  %#v\n\n\twant: %#v\n", got, want)
		return false
	}
	return true
}

// Equal fails the test if 'got' is not equal to 'want' using reflect.DeepEqual.
func (a *Assert) Equal(got, want any) bool {
	a.tb.Helper()
	if !reflect.DeepEqual(got, want) {
		a.logf("\n\n\tgot:  %#v\n\n\twant: %#v\n", got, want)
		return false
	}
	return true
}

// EqualPointer fails the test if 'got' is not equal to 'want' for pointers.
func (a *Assert) EqualPointer(got, want any) bool {
	a.tb.Helper()
	gotP := reflect.ValueOf(got).Pointer()
	wantP := reflect.ValueOf(want).Pointer()
	if gotP != wantP {
		a.logf("\n\n\tgot:  %#v\n\n\twant: %#v\n", gotP, wantP)
		return false
	}
	return true
}

// EqualError fails the test if 'got' is not equal to 'want' for errors.
func (a *Assert) EqualError(got, want error) bool {
	a.tb.Helper()
	if !errors.Is(got, want) && !errors.Is(want, got) {
		a.logf("\n\n\tgot:  %#v\n\n\twant: %#v\n", got, want)
		return false
	}
	return true
}

// EqualType fails the test if 'got' is not equal to 'want' for concrete types.
func (a *Assert) EqualType(got, want any) bool {
	a.tb.Helper()
	gotT := reflect.TypeOf(got)
	wantT := reflect.TypeOf(want)
	if gotT != wantT {
		a.logf("\n\n\tgot:  %#v\n\n\twant: %#v\n", gotT, wantT)
		return false
	}
	return true
}

// EqualComparer fails the test if the comparer is not true for the given values.
func (a *Assert) EqualComparer(got, want any, cmp Comparer) bool {
	a.tb.Helper()
	if !cmp(got, want) {
		a.logf("\n\n\tgot:  %#v\n\n\twant: %#v\n", got, want)
		return false
	}
	return true
}

// EqualFloat fails the test if 'got' is not equal to 'want' for float64 values.
func (a *Assert) EqualFloat(got, want, epsilon float64) bool {
	a.tb.Helper()
	if math.Abs(got-want) > epsilon {
		a.logf("\n\n\tgot:  %#v\n\n\twant: %#v\n\n\tepsilon: %#v", got, want, epsilon)
		return false
	}
	return true
}

// EqualPattern fails the test if 'got' does not match 'want' regular expression pattern.
func (a *Assert) EqualPattern(got any, want string) bool {
	a.tb.Helper()
	return a.EqualRegex(got, regexp.MustCompile(want))
}

// EqualRegex fails the test if 'got' does not match 'want' regular expression.
func (a *Assert) EqualRegex(got any, want *regexp.Regexp) bool {
	a.tb.Helper()
	gotStr := stdlib.MustE(func() (string, error) { return stdlib.ToString(got) })
	if !want.MatchString(gotStr) {
		a.logf("\n\n\tgot:  %#v\n\n\twant: %#v\n", gotStr, want.String())
		return false
	}
	return true
}

// Panic fails the test if the 'got' function does not panic.
func (a *Assert) Panic(got func()) bool {
	a.tb.Helper()
	defer func() {
		if r := recover(); r != nil {
			a.tb.Logf("\n\n\rpanic: %v", r)
		}
	}()
	got()
	a.logf("\n\n\tgot:  no panic\n\n\twant: panic\n")
	return false
}
