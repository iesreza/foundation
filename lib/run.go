package lib

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func Run(command string, args ...string) (string, error) {
	command = strings.TrimSpace(command)
	if strings.Contains(command, " ") {
		args = append(strings.Fields(command), args...)
		return Run(args[0], args[1:]...)
	}
	c := exec.Command(command, args...)
	out, err := c.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

type RunControl struct {
	OnFinish func()
	OnStdout func()
}

func RunInBackground(command string, args ...string) error {
	command = strings.TrimSpace(command)
	if strings.Contains(command, " ") {
		args = append(strings.Fields(command), args...)
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
