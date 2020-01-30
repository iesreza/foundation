package lib

import (
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
