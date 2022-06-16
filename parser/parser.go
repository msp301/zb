package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/msp301/zb/util"
)

func Parse(filepath string) []Note {
	fileReader, err := os.Open(filepath)
	defer fileReader.Close()

	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(fileReader)

	fileScanner.Split(bufio.ScanLines)

	var notes []Note

	var content string
	var contentStart int
	var title string
	var ids []uint64
	tagMap := map[string]int{}
	var links []uint64

	lineNum := 1
	for fileScanner.Scan() {
		line := fileScanner.Text()

		if len(strings.TrimSpace(line)) == 0 {
			if contentStart > 0 {
				content += fmt.Sprintln(line)
			}
			continue
		}

		isNoteDivider, _ := regexp.Match(`^\s*-{3,}`, fileScanner.Bytes())
		if isNoteDivider {
			note := Note{
				Content: strings.TrimSpace(content),
				File:    filepath,
				Links:   links,
				Tags:    ConvertTagMapToTagNames(tagMap),
				Title:   title,
			}
			if len(ids) == 1 {
				note.Id = ids[0]
			}
			notes = append(notes, note)

			content = ""
			contentStart = 0
			ids = []uint64{}
			links = []uint64{}
			tagMap = map[string]int{}
			title = ""

			continue
		}

		if contentStart > 0 {
			content += fmt.Sprintln(line)
		}

		if len(title) == 0 && !util.IsMetadataString(line) {
			contentStart = lineNum
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
			for _, match := range util.TagRegex.FindAllStringSubmatch(line, -1) {
				tag := match[1]
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

	note := Note{
		Content: strings.TrimSpace(content),
		File:    filepath,
		Links:   links,
		Tags:    ConvertTagMapToTagNames(tagMap),
		Title:   title,
	}
	if len(ids) == 1 {
		note.Id = ids[0]
	}
	notes = append(notes, note)

	return notes
}
