package util

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
	"regexp"
	"strings"
)

func Matches(str1 string, str2 string) bool {
	if len(str1) < 3 {
		lowerStr1 := strings.ToLower(str1)
		lowerStr2 := strings.ToLower(str2)
		if strings.HasPrefix(lowerStr2, lowerStr1) {
			return true
		}
		match, err := regexp.MatchString("^(\\W+)?"+lowerStr1, lowerStr2)
		if err == nil && match {
			return true
		}
		return false
	}
	if fuzzy.MatchFold(str1, str2) {
		return true
	} else {
		distance := fuzzy.LevenshteinDistance(strings.ToLower(str1), strings.ToLower(str2))
		if distance <= len(str2)/3 {
			return true
		}
	}
	return false
}
