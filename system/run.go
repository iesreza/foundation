package system

import (
	"bufio"
	"github.com/iesreza/foundation/lib/log"
	"io"
	"os"
	"os/exec"
	"runtime"
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

type ExecParams struct {
	ReturnResult bool
	OnStdScan    func(s string)
	SplitFunc    bufio.SplitFunc
	PrintOutput  bool
}

var defaultExecParams = ExecParams{
	true, nil, bufio.ScanRunes, false,
}

func Exec(command string, params *ExecParams) (string, error) {
	if params == nil {
		params = &defaultExecParams
	}
	output := ""
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}

	if params.PrintOutput {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Start()
		if err != nil {
			return output, err
		}

	} else {
		stderr, _ := cmd.StderrPipe()
		stdout, _ := cmd.StdoutPipe()

		err := cmd.Start()
		if err != nil {
			return output, err
		}

		scanner := bufio.NewScanner(io.MultiReader(stderr, stdout))
		scanner.Split(params.SplitFunc)

		for scanner.Scan() {
			if params.ReturnResult {
				output += scanner.Text()
			}
			if params.OnStdScan != nil {
				params.OnStdScan(scanner.Text())
			}
		}
	}

	cmd.Wait()
	return output, nil
}
