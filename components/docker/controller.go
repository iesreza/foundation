package docker

import (
	"fmt"
	"github.com/iesreza/foundation/system"
	"github.com/iesreza/gutil/path"
	"regexp"
	"strings"
)

type Controller struct{}

var gitRegex = regexp.MustCompile(`(?m)url\s*=\s*(https{0,1}\:\/\/.+)(\.git){0,1}`)
var gitPathRegex = regexp.MustCompile(`(?m)https*:\/\/(.+)`)
var gitRepo = ""
var gitPath = ""
var gitName = ""
var gitOwner = ""

func RegisterCLI() {
	type dockerCLI struct {
		Action string `short:"a" long:"action" description:"docker action" choice:"init" choice:"remove" choice:"build"`
	}

	system.RegisterCLI("docker", &dockerCLI{}, func(command string, data interface{}) {
		version, err := system.Exec("docker -v", nil)
		if err != nil || !strings.Contains(version, "Docker version") {
			fmt.Println("Docker is not installed on this system")
			return
		} else {
			fmt.Println("Installed Docker Version:")
			fmt.Println(version)
		}

		git := path.File(system.DIR + "/.git/config")
		if !git.Exist() {
			fmt.Println(".git/config not found. To use docker first you have to put project on git.")
			return
		}
		content, _ := git.Content()
		chunks := gitRegex.FindStringSubmatch(content)

		if len(chunks) < 2 {
			fmt.Println("proper .git/config not found. To use docker first you have to put project on git.")
			return
		}
		gitRepo = chunks[1]
		gitPath = strings.TrimRight(gitPathRegex.FindStringSubmatch(gitRepo)[1], ".git")
		chunks = strings.Split(gitPath, "/")
		gitName = chunks[len(chunks)-1]
		if len(chunks) > 1 {
			gitOwner = chunks[len(chunks)-2]
		} else {
			gitOwner = "foundation"
		}

		switch data.(*dockerCLI).Action {
		case "init":
			createDockerFile()
			break
		case "remove":
			removeDockerFile()
			break
		case "build":
			buildDockerFile()
			break
		}

	}, "Dockerfile (create,build,remove)")

}

func buildDockerFile() {
	system.Exec("docker build -t "+gitOwner+"/"+gitName+" "+system.DIR, &system.ExecParams{ReturnResult: false, PrintOutput: true})
	fmt.Println(" ")
}

func removeDockerFile() {
	dockerfile.Remove()
	fmt.Println("Dockerfile removed")
}

var dockerfile = path.File(system.DIR + "/Dockerfile")

func createDockerFile() {
	if dockerfile.Exist() {
		fmt.Println("Dockerfile exist. try remove docker file using remove command first")
		return
	}

	dockerfile.Create("#Autogenreated Dockerfile\n")
	dockerfile.Append("FROM iron/go:dev\n")
	dockerfile.Append("WORKDIR /app\n")
	dockerfile.Append("ENV SRC_DIR=$GOPATH/src/" + gitPath + "/\n")
	dockerfile.Append("RUN go get " + gitPath + "\n")
	dockerfile.Append("ADD ./config.yml /app/\n")
	dockerfile.Append("RUN cd $SRC_DIR; go build -o /app/" + gitName + ";\n")
	dockerfile.Append("ENTRYPOINT [\"/app/" + gitName + "\"]\n")

	fmt.Println("Dockerfile created successfully")

}
