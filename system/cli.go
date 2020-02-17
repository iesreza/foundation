package system

import (
	"bufio"
	"fmt"
	"github.com/iesreza/foundation/lib/cmd"
	"github.com/iesreza/go-flags"
	"os"
	"reflect"
	"strings"
	"time"
)

type cli struct {
	Command   string
	Structure interface{}
	OnCall    func(command string, data interface{})
	Help      string
}

var commandList = map[string]*cli{}

func RegisterCLI(command string, structure interface{}, onCall func(command string, data interface{}), help string) {

	commandList[command] = &cli{
		Command:   command,
		Structure: structure,
		OnCall:    onCall,
		Help:      help,
	}

}

func AliasCLI(command string, alias ...string) {

	for _, item := range alias {
		commandList[item] = commandList[command]
	}

}

func ListenCLI() {
	fmt.Println("Listen for commands")
	fmt.Println("Type help to see commands")

	for {
		cmdin := WaitForConsole(GetConfig().App.Title + ">")
		if strings.TrimSpace(cmdin) != "" {
			if !TryParseCommand(cmdin) {
				cmd.Error("Invalid command: " + cmdin)
			}
		}
	}
}

func TryParseCommand(cmd string) bool {
	cmd = strings.TrimSpace(cmd)
	for command, cli := range commandList {
		if strings.HasPrefix(cmd, command) {

			opt := cli.Structure
			var parser = flags.NewParser(opt, flags.Default)
			fields := strings.Fields(strings.Replace(cmd, command, "", 1))
			containArgs := false
			for _, item := range fields {
				if item[0] == '-' {
					containArgs = true
					break
				}
			}

			v := reflect.ValueOf(opt)
			i := reflect.Indirect(v)
			s := i.Type()

			fieldsNeededReview := 0
			for r := 0; r < s.NumField(); r++ {
				if s.Field(r).Tag.Get("short") != "" && s.Field(r).Tag.Get("default") == "" {
					fieldsNeededReview++
				}
			}

			if !containArgs && fieldsNeededReview > 0 {
				return false
			}
			_, err := parser.ParseArgs(cli.Command, fields)
			if err == nil {
				cli.OnCall(command, opt)
			}
			return true
		}
	}
	return false
}

func WaitForConsole(hint string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n" + hint)
	cmd, _ := reader.ReadString('\n')
	return cmd
}

func WaitForConsoleTimeout(hint string, timeout time.Duration, onTimeout func()) string {

	var cmd string
	ch := make(chan int)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\n" + hint)
		cmd, _ = reader.ReadString('\n')
		ch <- 1
	}()

	select {
	case <-ch:
		return cmd
	case <-time.After(timeout):
		if onTimeout != nil {
			onTimeout()
		}
		return ""
	}
}
