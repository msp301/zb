package editor

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/msp301/zb/util"
	"github.com/spf13/viper"
)

const filePlaceholder = "{file}"
const linePlaceholder = "{line}"

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
	command.Args = append(command.Args[:1], arguments(command.Args[1:], path, line)...)

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

func arguments(template []string, path string, line int) []string {
	if line < 1 {
		line = 1
	}

	replacer := strings.NewReplacer(filePlaceholder, path, linePlaceholder, strconv.Itoa(line))

	hasFile := false
	args := make([]string, 0, len(template)+1)
	for _, arg := range template {
		if strings.Contains(arg, filePlaceholder) {
			hasFile = true
		}
		args = append(args, replacer.Replace(arg))
	}

	if !hasFile {
		args = append(args, path)
	}

	return args
}
