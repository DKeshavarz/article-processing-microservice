package tagextractor

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"unicode"
)

type CharTransform func(rune) rune

func ExtractTags(body string, n int) ([]string, error) {
	body = Normalize(body, ToLower, ReplacePunctuationWithSpace)
	words := strings.Fields(body)
	words = RemoveStopwords(words)
	slices.Sort(words)

	type wordFreq struct{
		Word string
		count int
	}
	wordsFreq := make([]wordFreq, 0)

	cnt := 1
	for i := 1 ; i <= len(words) ; i++ {
		if i == len(words) || words[i] != words[i-1] {
			wordsFreq = append(wordsFreq, wordFreq{Word: words[i-1], count: cnt})
			cnt = 1
		} else {
			cnt++
		}
	}
	
	sort.Slice(wordsFreq, func(i, j int) bool {
		return wordsFreq[i].count > wordsFreq[j].count
	})

	if len(wordsFreq) < n {
		return nil, ErrWrongHighTagCount
	}
	for i := range n {
		words[i] = wordsFreq[i].Word
	}

	return words[:n], nil
}

func Normalize(input string, transforms ...CharTransform) string {
	runes := []rune(input)
	for i, r := range runes {
		for _, transform := range transforms {
			r = transform(r)
		}
		runes[i] = r
	}
	return string(runes)
}

func ToLower(r rune) rune {
	return unicode.ToLower(r)
}

func ReplacePunctuationWithSpace(r rune) rune {
	if unicode.IsPunct(r) {
		return ' '
	}
	return r
}

//TODO:document this function
func RemoveStopwords(words []string) []string{
	stopwords := map[string]bool{
		"the": true,
		"and": true,
		"of": true,
		"in": true,
		"to": true,
		"is": true,
	}	
	left, right := 0, len(words) - 1
	for left <= right {
		if !stopwords[words[left]] {
			left++
		} else {
			words[left], words[right] = words[right], words[left]
			right--
		}
	}
	fmt.Println(words, words[:left])
	return words[:left]
}
