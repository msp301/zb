package util

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func FileId(path string) (uint64, error) {
	fileName := filepath.Base(path)
	if !strings.HasSuffix(fileName, ".md") {
		return 0, fmt.Errorf("Not a markdown file")
	}

	nameWithoutSuffix := strings.TrimSuffix(fileName, ".md")
	if !IdRegex.MatchString(nameWithoutSuffix) {
		return 0, fmt.Errorf("Filename is not an ID")
	}

	fileId, err := strconv.ParseUint(nameWithoutSuffix, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse ID")
	}

	return fileId, nil
}
