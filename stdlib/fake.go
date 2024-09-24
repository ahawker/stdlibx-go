//go:generate go run github.com/abice/go-enum --marshal --names
package stdlib

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync/atomic"
	"time"
)

// RNG a default global random number generator.
var RNG *rand.Rand

func init() {
	if seed, ok := os.LookupEnv("STDLIB_TEST_SEED"); ok {
		RNG = rand.New(rand.NewSource(int64(MustInt(seed))))
	} else {
		RNG = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
}

// FakeStrategy represents a strategy for choosing a fake value.
//
// unspecified: No strategy specified.
// random: Randomly generate a value with no bounds or limitations.
// random_range: Randomly select a value within the bounds.
// random_select: Randomly select a value from loaded set of values.
// distribution_normal: Select a value from a normal distribution.
// distribution_uniform: Select a value from a uniform distribution.
// stateful: Select a value based on some state/previous value.
//
// ENUM(unspecified, random, random_range, random_pattern, random_select, distribution_normal, distribution_uniform, stateful).
type FakeStrategy string

// FakeState holds persistent values for some strategies.
type FakeState[T any] struct {
	// Generation is the numeric value of the previous generation (incrementing).
	Generation uint64
	// Init is the initial stored value from the first generation.
	Init T
	// Curr is the stored value from the previous generation.
	Curr T
}

// FakeConstraints define constraints for new fake values.
type FakeConstraints[T any] struct {
	// Cardinality is the number of distinct values.
	Cardinality uint64
	// Dataset is a set of values to choose from.
	Dataset []T
	// Min is the minimum value.
	Min T
	// Max is the maximum value.
	Max T
	// Pattern is a regular expression that the value must match.
	Pattern string
}

// FakeOptions for generating fake values.
type FakeOptions[T any] struct {
	// Strategy to use for selecting fake values.
	Strategy FakeStrategy
	// Constraints that define constraints for new fake values.
	Constraints *FakeConstraints[T]
	// State holds necessary persistent values for some strategies.
	State *FakeState[T]
	// RandomFn is a function that generates a random value without limitations.
	RandomFn func(ctx context.Context, options *FakeOptions[T]) T
	// RangeFn is a function that generates a random value within the bounds.
	RangeFn func(ctx context.Context, options *FakeOptions[T], min, max T) T
	// SelectFn is a function that generates a random value from a set of values.
	SelectFn func(ctx context.Context, options *FakeOptions[T], items []T) T
	// StateFn is a function that generates a value based on some state/previous value.
	StateFn func(ctx context.Context, options *FakeOptions[T])
}

// Generate generates a fake value based on the configured options.
func (o *FakeOptions[T]) Generate(ctx context.Context) T {
	switch o.Strategy {
	case FakeStrategyRandom:
		return o.RandomFn(ctx, o)
	case FakeStrategyRandomRange:
		return o.RangeFn(ctx, o, o.Constraints.Min, o.Constraints.Max)
	case FakeStrategyRandomSelect:
		return o.SelectFn(ctx, o, o.Constraints.Dataset)
	case FakeStrategyStateful:
		atomic.AddUint64(&o.State.Generation, 1)
		o.StateFn(ctx, o)
		return o.State.Curr
	default:
		panic(fmt.Sprintf("strategy %s not supported", o.Strategy))
	}
}
