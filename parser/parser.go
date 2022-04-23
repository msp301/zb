package parser

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/msp301/zb/util"
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

func Parse(filepath string) Note {
	fileReader, err := os.Open(filepath)
	defer fileReader.Close()

	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(fileReader)

	fileScanner.Split(bufio.ScanLines)

	var content string
	var content_start int
	var title string
	ids := []uint64{}
	tagMap := map[string]int{}
	links := []uint64{}

	lineNum := 1
	for fileScanner.Scan() {
		line := fileScanner.Text()

		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		if content_start > 0 {
			content += fmt.Sprintln(line)
		}

		if len(title) == 0 && !util.IsMetadataString(line) {
			content_start = lineNum
			title = line
			if strings.HasPrefix(line, "# ") {
				title = line[2:]
				continue
			}
			content += fmt.Sprintln(line)
		}

		if util.IdRegex.MatchString(line) {
			id, err := strconv.ParseUint(line, 0, 64)
			if err == nil {
				ids = append(ids, id)
			}
		}

		if strings.Contains(line, `#`) {
			for _, tag := range util.TagRegex.FindAllString(line, -1) {
				if tagMap[tag] == 0 && !strings.Contains(line, fmt.Sprintf("\\%s", tag)) {
					tagMap[tag] = lineNum
				}
			}
		}

		if strings.Contains(line, `[[`) {
			for _, str := range util.LinkRegex.FindAllStringSubmatch(line, -1) {
				link, err := strconv.ParseUint(str[1], 0, 64)
				if err == nil {
					links = append(links, link)
				}
			}
		}
		lineNum++
	}

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

	return Note{
		Content: strings.TrimSpace(content),
		File:    filepath,
		Id:      ids[0],
		Links:   links,
		Tags:    tagNames,
		Title:   title,
	}
}
