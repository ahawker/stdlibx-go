//go:generate go run github.com/abice/go-enum --marshal --names
package stdlib

// FakeStrategy represents a strategy for choosing a fake value.
//
// unspecified: No strategy specified.
// random: Randomly select a value within the bounds/pattern.
// random_selection: Randomly select a value from loaded set of values.
// normal_distribution: Select a value from a normal distribution.
// uniform_distribution: Select a value from a uniform distribution.
// stateful: Select a value based on some state/previous value.
//
// ENUM(unspecified, random, random_selection, normal_distribution, uniform_distribution, stateful).
type FakeStrategy string

type FakeState[T any] struct {
	// Generation is the numeric value of the previous generation (incrementing).
	Generation uint64
	// Init is the initial stored value from the first generation.
	Init T
	// Prev is the stored value from the previous generation.
	Prev T
}

// FakeRules define constraints for new fake values.
type FakeRules[T any] struct {
	// Cardinality is the number of distinct values.
	Cardinality uint64
	// LowerBound is the minimum value.
	LowerBound T
	// UpperBound is the maximum value.
	UpperBound T
	// Matches is a regular expression that the value must match.
	Matches string
	// DistributionMean is the mean value for the distribution.
	DistributionMean float64
	// Incrementing is true if the value should always increment and never go backwards.
	Incrementing bool
}

type FakeOptions[T any] struct {
	// Strategy to use for selecting fake values.
	Strategy FakeStrategy
	// Rules that define constraints for new fake values.
	Rules *FakeRules[T]
	// Next function mutates the state to the next generation.
	Next func(o *FakeOptions[T])
	// state holds necessary persistent values for some strategies.
	state *FakeState[T]
}
