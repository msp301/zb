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
	TagRegex = regexp.MustCompile(`(?:^|\s+)(#+\.?[0-9a-zA-Z_-][#0-9a-zA-Z_-]+)`)
	LinkRegex = regexp.MustCompile(`\B\[\[([\d\s\/:]+)\]\]`)
}

func IsMetadataString(line string) bool {
	return IsMetadata([]byte(line))
}

func IsMetadata(line []byte) bool {
	if IdRegex.Match(line) {
		return true
	}
	var noTags = TagRegex.ReplaceAll(line, []byte{})
	var noLinks = LinkRegex.ReplaceAll(noTags, []byte{})
	var noSpaces = bytes.TrimSpace(noLinks)

	return len(noSpaces) == 0
}
