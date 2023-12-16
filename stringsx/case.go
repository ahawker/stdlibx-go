package stringsx

import (
    "golang.org/x/text/cases"
    "golang.org/x/text/language"
)

// TitleCase returns the given string in title case per English language rules.
//
// https://en.wikipedia.org/wiki/Title_case
func TitleCase(s string) string {
    return cases.Title(language.English, cases.Compact).String(s)
}
