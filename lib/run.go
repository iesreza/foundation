package lib

import (
	"bufio"
	"fmt"
	"github.com/iesreza/foundation/lib/log"
	"os/exec"
	"strings"
)

func Run(command string, args ...string) string {
	command = strings.TrimSpace(command)
	if strings.Contains(command, " ") {
		args = append(strings.Fields(command), args...)
		return Run(args[0], args[1:]...)
	}
	c := exec.Command(command, args...)
	out, err := c.Output()
	if err != nil {
		log.Error(err, "Unable to run %s %v", command, args)
	}
	return string(out)
}

type RunControl struct {
	OnFinish func()
	OnStdout func()
}
type cmdOut struct {
	output []byte
	error  error
}

func RunInBackground(ctl RunControl, command string, args ...string) {
	command = strings.TrimSpace(command)
	if strings.Contains(command, " ") {
		args = append(strings.Fields(command), args...)
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Error(err)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
