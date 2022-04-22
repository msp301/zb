package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/msp301/zb/util"
)

func main() {

	filename := os.Args[1]

	fileReader, err := os.Open(filename)
	defer fileReader.Close()

	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(fileReader)

	fileScanner.Split(bufio.ScanLines)

	ids := []string{}
	tags := []string{}
	links := []string{}

	lineNum := 2
	for fileScanner.Scan() {
		if util.IdRegex.Match((fileScanner.Bytes())) {
			ids = append(ids, fileScanner.Text())
		}
		if strings.Contains(fileScanner.Text(), `#`) {
			tags = append(tags, util.TagRegex.FindAllString(fileScanner.Text(), -1)...)
		}
		if strings.Contains(fileScanner.Text(), `[[`) {
			links = append(links, util.LinkRegex.FindAllString(fileScanner.Text(), -1)...)
		}
		lineNum++
	}

	fmt.Println("ID: ", ids)
	fmt.Println("Tags: ", tags)
	fmt.Println("Links: ", links)
	fmt.Println("File:", filename)
}
