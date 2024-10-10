package stdlib

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TitleCase returns the given string in title case per English language rules.
//
// Ref: https://en.wikipedia.org/wiki/Title_case
func TitleCase(s string) string {
	return cases.Title(language.English, cases.Compact).String(s)
}

// ZeroPad returns the given integer as a zero-padded string.
func ZeroPad[T constraints.Integer](width T, v T) string {
	return fmt.Sprintf("%0*d", width, v)
}
