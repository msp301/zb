package util

import (
	"os/exec"
	"strings"
)

func Command(commandLine string, args ...string) *exec.Cmd {
	if _, err := exec.LookPath(commandLine); err == nil {
		return exec.Command(commandLine, args...)
	}

	fields := strings.Fields(commandLine)
	if len(fields) == 0 {
		return exec.Command(commandLine, args...)
	}

	return exec.Command(fields[0], append(fields[1:], args...)...)
}
