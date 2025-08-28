package util

import "unicode/utf8"

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) > n
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) < n
}
