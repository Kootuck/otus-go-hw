package main

import (
	"regexp"
	"sort"
	"strings"
)

var (
	regexRemoveSurroundingPunct = regexp.MustCompile(`^\p{P}+|\p{P}+$`)
	onlyLatinOrCyrillicSymbols  = regexp.MustCompile(`^[\p{Latin}\p{Cyrillic}\s]*$`)
)

type keyValue struct {
	word  string
	count int
}

func Top10(input string) []string {
	var sortedWordsResult []string
	// 0. input string -> slice of words
	words := strings.Fields(input)
	// 1. Sanitize words
	sanitizedWords := sanititzeWords(words)
	// 2. Convert input text into map [word(string) : count(int)]
	wordsCount := countWords(sanitizedWords)
	// 3. Convert a map into a sorted kv slice
	sortedWordsResult = sortWords(wordsCount)
	// 4. topTen := make([]string, 10)
	if len(sortedWordsResult) == 0 {
		return sortedWordsResult
	}

	return sortedWordsResult[:10]
}

func sortWords(input map[string]int) []string {
	// 1. Convert map into kv slice
	wordCount := make([]keyValue, 0, len(input))
	for w, c := range input {
		wordCount = append(wordCount, keyValue{word: w, count: c})
	}
	// 2. Sort kv slice by value
	sort.Slice(wordCount, func(i, j int) bool {
		// 2.1. Lexigraphic sort for elements of equal length
		if wordCount[i].count == wordCount[j].count {
			return wordCount[i].word < wordCount[j].word
		}
		// 2.2. Sort by occurrences from highest to lowest
		return wordCount[i].count > wordCount[j].count
	})

	result := make([]string, len(wordCount))
	for i, wordCountElem := range wordCount {
		result[i] = wordCountElem.word
	}

	return result
}

func sanititzeWords(words []string) (szdWords []string) {
	szdWords = make([]string, len(words))
	k := 0

	for _, w := range words {
		sanitized := removePunctuation(w)
		// "-" is not a word
		if len(sanitized) == 1 && rune(sanitized[0]) == '-' {
			continue
		}
		// empty string is not a word
		if len(sanitized) == 0 {
			continue
		}
		// only valid words are counted
		if !onlyLatinOrCyrillicSymbols.MatchString(sanitized) {
			continue
		}

		szdWords[k] = sanitized
		k++
	}

	return szdWords[:k]
}

func countWords(words []string) (wordCount map[string]int) {
	wordCount = make(map[string]int)

	for _, w := range words {
		wordCount[w]++
	}
	return wordCount
}

func removePunctuation(input string) string {
	if len(input) > 1 {
		noPunctuation := regexRemoveSurroundingPunct.ReplaceAllString(input, "")
		return strings.ToLower(noPunctuation)
	}
	return strings.ToLower(input)
}
