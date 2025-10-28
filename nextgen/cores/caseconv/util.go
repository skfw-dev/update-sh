package caseconv

import (
	"iter"
	"strings"
	"unicode"
)

// eachWords ...
func eachWords(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	// Convert string to []rune for safe index-based look ahead/look back
	// Ensuring string is full UTF-8 encode/decode compatibility.
	runes := []rune(s)

	var words []string
	var currentWord strings.Builder

	for i, r := range runes {

		// Boundary Check 1: Delimiters (space, hyphen, underscore)
		if unicode.IsSpace(r) || r == '-' || r == '_' {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}

			continue
		}

		// Boundary Check 2: Case changes (only check after the first rune)
		if unicode.IsUpper(r) && i > 0 {
			prevIndex, nextIndex := i-1, i+1

			prevRune := runes[prevIndex]
			isPrevUpper := unicode.IsUpper(prevRune)
			isPrevLower := unicode.IsLower(prevRune)

			// Case A: lowercase followed by uppercase (e.g., "fooBar" -> "foo", "Bar").
			// Case B: Acronym boundary (e.g., "HTTPApi" -> "HTTP", "Api").
			if isPrevLower || (isPrevUpper && nextIndex < len(runes) && unicode.IsLower(runes[nextIndex])) {
				if currentWord.Len() > 0 {
					words = append(words, currentWord.String())
					currentWord.Reset()
				}
			}
		}

		currentWord.WriteRune(r)
	}

	// Add the last word if any was being built
	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}

// toUniformCase ...
func toUniformCase(s, sep string, upper bool) string {
	words := eachWords(s)
	if len(words) == 0 {
		return ""
	}

	transform := SelectTransform(upper, strings.ToUpper, strings.ToLower)
	return strings.Join(MapFunc(words, transform), sep)
}

// toMixedCase ...
func toMixedCase(s, sep string, lower bool) string {
	words := eachWords(s)
	if len(words) == 0 {
		return ""
	}

	words = MapFunc(words, toTitleUnlessAcronym)
	if len(words) > 0 && lower {
		words[0] = toLowerFirstChar(words[0])
	}

	return strings.Join(words, sep)
}

// toLowerFirstChar ...
func toLowerFirstChar(s string) string {
	if len(s) == 0 {
		return ""
	}

	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// toTitleUnlessAcronym ...
func toTitleUnlessAcronym(s string) string {
	if !isAcronym(s) {
		return toTitle(s)
	}

	return s
}

// isAcronym ...
func isAcronym(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if !unicode.IsUpper(char) {
			return false
		}
	}
	return true
}

// toTitle ...
func toTitle(s string) string {
	size := len(s)
	temp := []rune(s)
	if size > 0 {
		if size > 1 {
			return string(unicode.ToUpper(temp[0])) + strings.ToLower(string(temp[1:]))
		}

		return string(unicode.ToUpper(temp[0]))
	}

	return ""
}

// IterateSlice ...
func IterateSlice[T any](data []T, start, end, step int) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := start; i < end; i += step {
			yield(data[i])
		}
	}
}

// SelectTransform ...
func SelectTransform[E, R any](choice1 bool, transform1 func(E) R, transform2 func(E) R) func(E) R {
	if choice1 {
		return transform1
	}
	return transform2
}

// MapFunc ...
func MapFunc[E, R any](data []E, transform func(E) R) []R {
	if transform != nil {
		result := make([]R, len(data))
		for i, v := range data {
			result[i] = transform(v)
		}
		return result
	}

	return nil
}
