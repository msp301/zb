package pager

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/spf13/viper"
)

func FindPager() string {
	pager := viper.GetString("pager")
	if pager != "" {
		return pager
	}

	pager = os.Getenv("PAGER")
	if pager != "" {
		return pager
	}

	return ""
}

type Pager struct {
	pipe    io.WriteCloser
	command *exec.Cmd
	wg      sync.WaitGroup
}

func Open() (*Pager, error) {
	pagerCmd := FindPager()

	var pager = Pager{}
	pager.wg.Add(1)

	if pagerCmd == "" {
		pager.pipe = os.Stdout
	} else {
		pager.command = exec.Command(pagerCmd)
		pager.command.Stdout = os.Stdout
		pager.command.Stderr = os.Stderr

		var err error
		pager.pipe, err = pager.command.StdinPipe()
		if err != nil {
			return nil, err
		}
	}

	return &pager, nil
}

func (pager *Pager) Close() {
	pager.pipe.Close()
	pager.wg.Done()
}

func (pager *Pager) Wait() error {
	if pager.command == nil {
		pager.wg.Wait()
		return nil
	}

	return pager.command.Run()
}

func (pager *Pager) Writef(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(pager.pipe, format, args...)
}
