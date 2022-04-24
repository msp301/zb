package parser

import (
	"sort"
	"strings"
)

type Tag struct {
	Name    string
	LineNum int
}

type ByLineNum []Tag

func (l ByLineNum) Len() int { return len(l) }
func (l ByLineNum) Less(a, b int) bool {
	return l[a].LineNum < l[b].LineNum || l[a].LineNum == l[b].LineNum && strings.ToLower(l[a].Name) < strings.ToLower(l[b].Name)
}
func (l ByLineNum) Swap(a, b int) { l[a], l[b] = l[b], l[a] }

func ConvertTagMapToTagNames(tagMap map[string]int) []string {
	tags := make([]Tag, len(tagMap))
	index := 0
	for key, value := range tagMap {
		tags[index] = Tag{Name: key, LineNum: value}
		index++
	}
	sort.Sort(ByLineNum(tags))
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	return tagNames
}
