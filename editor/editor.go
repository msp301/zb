package editor

import (
	"github.com/spf13/viper"
	"os"
	"os/exec"
)

func findEditor() string {
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
	editor := findEditor()

	command := exec.Command(editor, path)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Start()
	if err != nil {
		return err
	}

	err = command.Wait()
	if err != nil {
		return err
	}

	return nil
}
