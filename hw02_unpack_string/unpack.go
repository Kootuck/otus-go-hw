package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")
var ErrIncorrectCharacters = errors.New("non-latin characters are not supported")

func Unpack(inputStr string) (string, error) {
	var builder strings.Builder

	// check 1
	if len(inputStr) == 0 {
		return "", nil
	}
	// check 2
	if consecutiveDigitsExist(inputStr) {
		return "", ErrInvalidString
	}
	// check 3
	if startsWithDigit(inputStr) {
		return "", ErrInvalidString
	}
	// reduce a slice of substrings into string
	for _, str := range unpackIntoSlice(inputStr) {
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
