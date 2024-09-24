//go:generate go run github.com/abice/go-enum --marshal --names
package stdlib

import (
	"context"
	"fmt"
	"golang.org/x/exp/constraints"
	"math/rand"
	"regexp"
	"sync/atomic"
	"time"
)

var FakeRequestKey = NewContextKey[*FakeRequest](
	"stdlib.FakeRequest",
	&FakeRequest{
		Generations: 1,
		Rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	},
)

func FakeContext(ctx context.Context, count uint64) context.Context {
	return FakeRequestKey.WithValue(ctx, &FakeRequest{Generations: count})
}

type FakeRequest struct {
	Generations uint64
	Rand        *rand.Rand
}

// FakeStrategy represents a strategy for choosing a fake value.
//
// unspecified: No strategy specified.
// random: Randomly generate a value with no bounds or limitations.
// random_range: Randomly select a value within the bounds.
// random_select: Randomly select a value from loaded set of values.
// random_pattern: Randomly select a value from a regex pattern.
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

type FakeConstraintsNumeric[T constraints.Integer | constraints.Float] struct {
	// Cardinality is the number of distinct values.
	Cardinality uint64
	// Dataset is a set of values to choose from.
	Dataset []T
	// Min is the minimum value.
	Min T
	// Max is the maximum value.
	Max T
}

type FakeConstraintsTextual[T ~string] struct {
	// Cardinality is the number of distinct values.
	Cardinality uint64
	// Dataset is a set of values to choose from.
	Dataset   []T
	MinLength uint64
	MaxLength uint64
}

// FakeNumber represents a fake number generator.
type FakeNumber[T constraints.Integer | constraints.Float] struct {
	// Strategy to use for selecting fake values.
	Strategy FakeStrategy
	// State holds necessary persistent values for some strategies.
	State *FakeState[T]

	// Cardinality is the number of distinct values relative to the total number of generations.
	Cardinality uint64
	// Possible is a fixed set of values to choose from.
	Possible []T
	// Min is the minimum value.
	Min T
	// Max is the maximum value.
	Max T

	// RandomFn is a function that generates a random value without limitations.
	RandomFn func(ctx context.Context, n *FakeNumber[T]) T
	// RangeFn is a function that generates a random value within the bounds.
	RangeFn func(ctx context.Context, n *FakeNumber[T]) T
	// SelectFn is a function that generates a random value from a set of values.
	SelectFn func(ctx context.Context, n *FakeNumber[T]) T
	// StateFn is a function that generates a value based on some state/previous value.
	StateFn func(ctx context.Context, n *FakeNumber[T])
}

// Generate generates a fake number based on the configured options.
func (n *FakeNumber[T]) Generate(ctx context.Context) T {
	req := FakeRequestKey.Value(ctx)

	// Defaults.
	if n.RandomFn == nil {
		n.RandomFn = func(ctx context.Context, n *FakeNumber[T]) T {
			return RandomNumber[T](req.Rand)
		}
	}
	if n.RangeFn == nil {
		n.RangeFn = func(ctx context.Context, n *FakeNumber[T]) T {
			return RandomNumberRange[T](req.Rand, n.Min, n.Max)
		}
	}
	if n.SelectFn == nil {
		n.SelectFn = func(ctx context.Context, n *FakeNumber[T]) T {
			return RandomSelection(req.Rand, n.Possible)
		}
	}

	switch n.Strategy {
	case FakeStrategyRandom:
		return n.RandomFn(ctx, n)
	case FakeStrategyRandomRange:
		return n.RangeFn(ctx, n)
	case FakeStrategyRandomSelect:
		return n.SelectFn(ctx, n)
	case FakeStrategyStateful:
		atomic.AddUint64(&n.State.Generation, 1)
		n.StateFn(ctx, n)
		return n.State.Curr
	default:
		panic(fmt.Sprintf("number strategy %s not supported", n.Strategy))
	}
}

type FakeText[T ~string] struct {
	// Strategy to use for selecting fake values.
	Strategy FakeStrategy
	// State holds necessary persistent values for some strategies.
	State *FakeState[T]

	// Cardinality is the number of distinct values relative to the total number of generations.
	Cardinality uint64
	// Possible is a fixed set of values to choose from.
	Possible []T
	// Regex is a regular expression that the value must match.
	Regex *regexp.Regexp
	// MinLength is the minimum text length.
	MinLength uint64
	// MaxLength is the maximum text value.
	MaxLength uint64

	// RandomFn is a function that generates a random value without limitations.
	RandomFn func(ctx context.Context, t *FakeText[T]) T
	// RangeFn is a function that generates a random value within the bounds.
	RangeFn func(ctx context.Context, t *FakeText[T]) T
	// PatternFn is a function that generates a random value that matches a regular expression pattern.
	PatternFn func(ctx context.Context, t *FakeText[T]) T
	// SelectFn is a function that generates a random value from a set of values.
	SelectFn func(ctx context.Context, t *FakeText[T]) T
	// StateFn is a function that generates a value based on some state/previous value.
	StateFn func(ctx context.Context, t *FakeText[T])
}

// Generate generates fake text based on the configured options.
func (t *FakeText[T]) Generate(ctx context.Context) T {
	req := FakeRequestKey.Value(ctx)

	// Defaults.
	if t.RandomFn == nil {
		t.RandomFn = func(ctx context.Context, t *FakeText[T]) T {
			return RandomString[T](req.Rand, t.MinLength, t.MaxLength)
		}
	}
	if t.RangeFn == nil {
		t.RangeFn = func(ctx context.Context, t *FakeText[T]) T {
			return RandomString[T](req.Rand, t.MinLength, t.MaxLength)
		}
	}
	if t.PatternFn == nil {
		t.PatternFn = func(ctx context.Context, t *FakeText[T]) T {
			return RandomRegex[T](req.Rand, t.Regex.String())
		}
	}
	if t.SelectFn == nil {
		t.SelectFn = func(ctx context.Context, t *FakeText[T]) T {
			return RandomSelection(req.Rand, t.Possible)
		}
	}

	switch t.Strategy {
	case FakeStrategyRandom:
		return t.RandomFn(ctx, t)
	case FakeStrategyRandomRange:
		return t.RangeFn(ctx, t)
	case FakeStrategyRandomPattern:
		return t.PatternFn(ctx, t)
	case FakeStrategyRandomSelect:
		return t.SelectFn(ctx, t)
	case FakeStrategyStateful:
		atomic.AddUint64(&t.State.Generation, 1)
		t.StateFn(ctx, t)
		return t.State.Curr
	default:
		panic(fmt.Sprintf("text strategy %s not supported", t.Strategy))
	}
}
