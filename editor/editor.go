package editor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/msp301/zb/util"
	"github.com/spf13/viper"
)

func FindEditor() string {
	editor := viper.GetString("editor")
	if editor != "" {
		return editor
	}

	editor = os.Getenv("EDITOR")
	if editor != "" {
		return editor
	}

	return ""
}

func Command(path string, line int) (*exec.Cmd, error) {
	editor := FindEditor()
	if editor == "" {
		return nil, errors.New("no editor configured, set $EDITOR or run 'zb config editor <program>'")
	}

	command := util.Command(editor)
	command.Args = append(command.Args, locationArgs(command.Args[0], path, line)...)

	return command, nil
}

func Open(path string) error {
	command, err := Command(path, 0)
	if err != nil {
		return err
	}

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func locationArgs(program string, path string, line int) []string {
	if line < 1 {
		return []string{path}
	}

	switch programName(program) {
	case "vi", "vim", "nvim", "gvim", "nano", "emacs", "emacsclient", "micro", "kak", "joe", "mg", "ne":
		return []string{fmt.Sprintf("+%d", line), path}
	case "code", "code-insiders", "codium", "cursor", "windsurf":
		return []string{"--goto", fmt.Sprintf("%s:%d", path, line)}
	case "hx", "helix", "subl", "zed":
		return []string{fmt.Sprintf("%s:%d", path, line)}
	case "notepad++":
		return []string{fmt.Sprintf("-n%d", line), path}
	}

	return []string{path}
}

func programName(program string) string {
	name := strings.ToLower(filepath.Base(program))
	for _, suffix := range []string{".exe", ".cmd", ".bat"} {
		name = strings.TrimSuffix(name, suffix)
	}
	return name
}
