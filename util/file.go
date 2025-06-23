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
		return 0, fmt.Errorf("'%s' is not a markdown file", path)
	}

	nameWithoutSuffix := strings.TrimSuffix(fileName, ".md")
	if !IdRegex.MatchString(nameWithoutSuffix) {
		return 0, fmt.Errorf("'%s' filename is not an ID", path)
	}

	fileId, err := strconv.ParseUint(nameWithoutSuffix, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("'%s' failed to parse ID", path)
	}

	return fileId, nil
}
