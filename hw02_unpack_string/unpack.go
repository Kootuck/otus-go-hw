package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidString       = errors.New("invalid string")
	ErrIncorrectCharacters = errors.New("non-latin characters are not supported")
)

func Unpack(inputStr string) (string, error) {
	var builder strings.Builder

	// 1. Check inputs before unpacking.
	if len(inputStr) == 0 {
		return "", nil
	}
	// 1.1. There must be no 2+ adjacent digits in input string.
	if consecutiveDigitsExist(inputStr) {
		return "", ErrInvalidString
	}
	// 1.2. Input string must start with a alphabetic character.
	if startsWithDigit(inputStr) {
		return "", ErrInvalidString
	}
	// 1.3. Method supports only latin characters.
	if containsNotAllowedCharacters(inputStr) {
		return "", ErrIncorrectCharacters
	}

	// 2. Break into a slice.
	inputSlice := unpackIntoSlice(inputStr)

	// 3. Reduce a slice of substrings into string.

	for _, str := range inputSlice {
		builder.WriteString(str)
	}

	return builder.String(), nil
}

func startsWithDigit(inputStr string) bool {
	return isDigit(inputStr[0])
}

func unpackIntoSlice(inputStr string) []string {
	unpacked := make([]string, len(inputStr))
	escapeMode := false
	// iterate,  if digit     -> replace previous character (with itself *digit times)
	// 			 if character -> add to slice
	for i, char := range []byte(inputStr) {
		if !escapeMode && isEscapeCharacter(char) {
			escapeMode = true
			continue
		}
		if !escapeMode && isDigit(char) {
			unpacked[i-1] = strings.Repeat(unpacked[i-1], int(byteToDigit(char)))
		} else {
			unpacked[i] = string(char)
		}
		escapeMode = false
	}
	return unpacked
}

func consecutiveDigitsExist(in string) bool {
	isDigitPrevious := isDigit(in[0])

	for i := 1; i < len(in); i++ {
		isDigitCurrent := isDigit(in[i]) && !isEscapeCharacter(in[i-1])
		if isDigitPrevious && isDigitCurrent {
			return true
		}
		isDigitPrevious = isDigitCurrent
	}
	return false
}

func isDigit(in byte) bool {
	if _, err := strconv.Atoi(string(in)); err == nil {
		return true
	}
	return false
}

func byteToDigit(in byte) int8 {
	value, _ := strconv.Atoi(string(in))
	return int8(value)
}

func isEscapeCharacter(char byte) bool {
	return char == '\\'
}

// Allowed: a Latin lettes, digits, forward slash or backslash.
func containsNotAllowedCharacters(str string) bool {
	for _, char := range str {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == 47 ||
			char == 92 {
			continue
		}
		return true
	}
	return false
}
