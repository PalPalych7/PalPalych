package hw03frequencyanalysis

import (
	"regexp"
	"strings"
)

var patern = `([^a-zа-яё-]|-[^a-zа-яё])+`

var compRegexp = regexp.MustCompile(patern)

func Top10(inputStr string) []string {
	var wordSl []string
	wordDict := make(map[string]int)
	wordSl = compRegexp.Split(strings.ToLower(inputStr), -1)
	for _, v := range wordSl {
		val := wordDict[v]
		wordDict[v] = val + 1
	}
	var res []string
	for i := 1; i <= 10; i++ {
		maxVal := 0
		maxKey := ""
		for key, val := range wordDict {
			if val > maxVal || (val == maxVal && key < maxKey) {
				maxVal = val
				maxKey = key
			}
		}
		if maxKey != "" {
			res = append(res, maxKey)
			delete(wordDict, maxKey)
		}
	}
	return res
}
