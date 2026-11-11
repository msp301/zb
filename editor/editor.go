package editor

import (
	"errors"
	"os"

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

func Open(path string) error {
	editor := FindEditor()
	if editor == "" {
		return errors.New("no editor configured, set $EDITOR or run 'zb config editor <program>'")
	}

	command := util.Command(editor, path)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
