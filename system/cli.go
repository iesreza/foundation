package system

import (
	"bufio"
	"fmt"
	"github.com/iesreza/foundation/log"
	"github.com/jessevdk/go-flags"
	"os"
	"strings"
)

type cli struct {
	Command   string
	Structure interface{}
	OnCall    func(command string, data interface{})
}

var commandList = map[string]*cli{}

func RegisterCLI(command string, structure interface{}, onCall func(command string, data interface{})) {

	commandList[command] = &cli{
		Command:   command,
		Structure: structure,
		OnCall:    onCall,
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
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\n>")
		cmd, _ := reader.ReadString('\n')
		if strings.TrimSpace(cmd) != "" {
			if !TryParseCommand(cmd) {
				log.Error("Invalid command: " + cmd)
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
			_, err := parser.ParseArgs(fields)
			if err == nil {
				cli.OnCall(command, opt)
			} else {
				fmt.Println(err)
			}
			return true
		}
	}
	return false
}
