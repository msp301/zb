package parser

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/msp301/zb/util"
)

func Parse(filepath string) Note {
	fileReader, err := os.Open(filepath)
	defer fileReader.Close()

	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(fileReader)

	fileScanner.Split(bufio.ScanLines)

	var title string
	ids := []uint64{}
	tags := []string{}
	links := []uint64{}

	lineNum := 1
	for fileScanner.Scan() {
		line := fileScanner.Text()

		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		if len(title) == 0 && !util.IsMetadataString(line) {
			title = line
			if strings.HasPrefix(line, "# ") {
				title = line[2:]
			}
		}

		if util.IdRegex.MatchString(line) {
			id, err := strconv.ParseUint(line, 0, 64)
			if err == nil {
				ids = append(ids, id)
			}
		}

		if strings.Contains(line, `#`) {
			tags = append(tags, util.TagRegex.FindAllString(line, -1)...)
		}

		if strings.Contains(line, `[[`) {
			for _, str := range util.LinkRegex.FindAllString(line, -1) {
				link, err := strconv.ParseUint(str, 0, 64)
				if err == nil {
					links = append(links, link)
				}
			}
		}
		lineNum++
	}

	return Note{
		File:  filepath,
		Id:    ids[0],
		Links: links,
		Tags:  tags,
		Title: title,
	}
}
