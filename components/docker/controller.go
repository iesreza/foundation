package docker

import (
	"bufio"
	"fmt"
	"github.com/iesreza/foundation/lib/log"
	"github.com/iesreza/foundation/system"
	"github.com/iesreza/gutil/path"
	"regexp"
	"strings"
)

type Controller struct{}

var gitRegex = regexp.MustCompile(`(?m)url\s*=\s*(https*\:\/\/.+\.git)`)
var gitPathRegex = regexp.MustCompile(`(?m)https*:\/\/(.+).git`)
var gitRepo = ""
var gitPath = ""
var gitName = ""
var gitOwner = ""

func RegisterCLI() {

	system.RegisterCLI("docker", &struct{}{}, func(command string, data interface{}) {
		version, err := system.Exec("docker -v", nil)
		if err != nil || !strings.Contains(version, "Docker version") {
			log.Error("Docker is not installed on this system")
			return
		} else {
			log.Notice("Installed Docker Version:")
			log.Notice(version)
			log.Notice("Type exit to exit")
		}

		git := path.File(system.DIR + "/.git/config")
		if !git.Exist() {
			log.Error(".git/config not found. To use docker first you have to put project on git.")
			return
		}
		content, _ := git.Content()
		chunks := gitRegex.FindStringSubmatch(content)
		if len(chunks) != 2 {
			log.Error("proper .git/config not found. To use docker first you have to put project on git.")
			return
		}
		gitRepo = chunks[1]
		gitPath = gitPathRegex.FindStringSubmatch(gitRepo)[1]
		chunks = strings.Split(gitPath, "/")
		gitName = chunks[len(chunks)-1]
		if len(chunks) > 1 {
			gitOwner = chunks[len(chunks)-2]
		} else {
			gitOwner = "foundation"
		}
		for {
			cmd := system.WaitForConsole(system.GetConfig().App.Title + " Docker>")
			lower := strings.ToLower(cmd)
			if strings.TrimSpace(lower) == "exit" {
				break
			}
			if strings.TrimSpace(cmd) != "" {

				switch strings.TrimSpace(cmd) {
				case "init", "create":
					createDockerFile()
					break
				case "remove":
					removeDockerFile()
					break
				case "build", "deploy":
					buildDockerFile()
					break
				}

			}

		}
	}, "Docker command line tool")

}

func buildDockerFile() {
	system.Exec("docker build -t "+gitOwner+"/"+gitName+" "+system.DIR, &system.ExecParams{ReturnResult: false, OnStdScan: func(s string) {
		fmt.Print(s)
	}, SplitFunc: bufio.ScanRunes})
	fmt.Println(" ")
}

func removeDockerFile() {
	dockerfile.Remove()
	fmt.Println("Dockerfile removed")
}

var dockerfile = path.File(system.DIR + "/Dockerfile")

func createDockerFile() {
	if dockerfile.Exist() {
		log.Error("Dockerfile exist. try remove docker file using remove command first")
		return
	}

	dockerfile.Create("#Autogenreated Dockerfile\n")
	dockerfile.Append("FROM iron/go:dev\n")
	dockerfile.Append("WORKDIR /app\n")
	dockerfile.Append("ENV SRC_DIR=" + gitPath + "/\n")
	dockerfile.Append("ENV go get " + gitPath + "\n")
	dockerfile.Append("ADD ./config.yml /app/\n")
	dockerfile.Append("RUN cd $SRC_DIR; go build -o /app/" + gitName + ";\n")
	dockerfile.Append("ENTRYPOINT [\"/app/" + gitName + "\"]\n")

	fmt.Println("Dockerfile created successfully")

}
