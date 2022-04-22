package util

import (
	"bytes"
	"regexp"
)

var IdRegex *regexp.Regexp
var TagRegex *regexp.Regexp
var LinkRegex *regexp.Regexp

func init() {
	IdRegex = regexp.MustCompile(`^(\d{12})$`)
	TagRegex = regexp.MustCompile(`\B(#\.?\w[^\s]+)`)
	LinkRegex = regexp.MustCompile(`\B\[\[([\d\s\/:]+)\]\]`)
}

func IsMetadata(line []byte) bool {
	if IdRegex.Match(line) {
		return true
	}
	var noTags = TagRegex.ReplaceAll(line, []byte{})
	var noLinks = LinkRegex.ReplaceAll(noTags, []byte{})
	var noSpaces = bytes.Trim(noLinks, "\r\n ")

	return len(noSpaces) == 0
}
