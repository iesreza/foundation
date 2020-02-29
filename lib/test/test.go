package main

// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
// to run:
// go run 03-live-progress-and-capture-v3.go

import (
	cmd2 "github.com/iesreza/foundation/lib/cmd"
	"os/exec"
	"time"
)

func Reciver(s string) {

}

func main() {

	cmd := exec.Command("longps", "4")

	for {
		cmd2.Info("Main proccess call")
		time.Sleep(3 * time.Second)
	}
}
