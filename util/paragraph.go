package util

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
	"regexp"
	"strings"
)

var noiseRegex = regexp.MustCompile(`[^\s\w#.-]`)

func ParagraphMatches(content string, query string) bool {
	content = noiseRegex.ReplaceAllString(content, " ")
	tokens := strings.Fields(content)
	for _, token := range tokens {
		if len(query) > 3 && strings.HasPrefix(token, query) {
			return true
		}

		var distance int
		hasUppercase := regexp.MustCompile("[A-Z]")
		if hasUppercase.MatchString(query) {
			distance = fuzzy.RankMatchNormalized(query, token)
		} else {
			distance = fuzzy.RankMatchNormalizedFold(query, token)
		}

		if distance == -1 {
			continue
		}

		if distance == 0 {
			return true
		}

		distancePercent := (float64(distance) / float64(len(token))) * 100
		thresholdPercent := 50.0
		if distancePercent < thresholdPercent || (distancePercent == thresholdPercent && len(token) == 2) {
			return true
		}
	}
	return false
}
