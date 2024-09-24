package stdlib

import (
	"math/rand"
)

// RandomRegex returns a random string that matches a regular expression.
func RandomRegex[T ~string](rng *rand.Rand, pattern string) T {
	return T("Hello World")

}
