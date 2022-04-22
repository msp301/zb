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

	ids := []uint64{}
	tags := []string{}
	links := []uint64{}

	lineNum := 1
	for fileScanner.Scan() {
		if util.IdRegex.Match((fileScanner.Bytes())) {
			id, err := strconv.ParseUint(fileScanner.Text(), 0, 64)
			if err == nil {
				ids = append(ids, id)
			}
		}
		if strings.Contains(fileScanner.Text(), `#`) {
			tags = append(tags, util.TagRegex.FindAllString(fileScanner.Text(), -1)...)
		}
		if strings.Contains(fileScanner.Text(), `[[`) {
			for _, str := range util.LinkRegex.FindAllString(fileScanner.Text(), -1) {
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
		Title: "",
	}
}
