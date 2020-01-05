package system

import (
	"github.com/iesreza/foundation/log"
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
