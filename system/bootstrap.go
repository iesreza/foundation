package system

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/iesreza/foundation/components/docker"
	"github.com/iesreza/foundation/language"
	"github.com/iesreza/foundation/lib"
	"github.com/iesreza/foundation/lib/log"
	"github.com/iesreza/foundation/lib/router"
	"github.com/ztrue/shutdown"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
)

func Boot() {
	SetupRouter()
	router.SessionAge = time.Duration(config.App.SessionAge) * time.Second
	StartWebServer()

	for _, item := range Components {
		item.Menu()
		item.Routers()
	}
}

func PreBoot() {
	cfg := GetConfig()

	if cfg.Tweaks.MaxProcessors < 1 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(cfg.Tweaks.MaxProcessors)
	}

	log.Register(&log.Logger{
		WriteToFile: cfg.Log.WriteFile,
		Concurrent:  true,
		Level:       log.ParseLevel(cfg.Log.Level),
		MaxAge:      cfg.Log.MaxAge,
		MaxSize:     cfg.Log.MaxSize,
		Path:        cfg.Log.Path,
	})

	if cfg.Tweaks.Ballast {
		size, err := lib.ParseSize(cfg.Tweaks.BallastSize)
		if err != nil {
			log.Error("Unable to parse ballast size: %s", cfg.Tweaks.BallastSize)
		} else {
			ballast := make([]byte, size)
			_ = ballast
		}
	}

	shutdown.Add(func() {
		log.Info("Shutting Down ...")
		//Shutdown event fire
		for _, item := range onShutdownCallbacks {
			item()
		}
	})

	RegisterCLI("help", &struct{}{}, func(command string, data interface{}) {
		list := make([]string, len(commandList))
		c := 0
		longest := 0
		for command, _ := range commandList {
			if len(command) > longest {
				longest = len(command)
			}
		}
		for command, v := range commandList {
			list[c] = command
			for i := len(command); i < longest+3; i++ {
				list[c] += " "
			}
			list[c] += v.Help
			c++
		}
		sort.Strings(list)
		fmt.Println("Available Commands:")
		for _, item := range list {
			fmt.Println(item)
		}

		fmt.Println("\r\nUsage:")
		fmt.Println("command [OPTIONS]")
		fmt.Println("\r\nHelp:")
		fmt.Println("command -h")
	}, "Show help")
	AliasCLI("help", "?")

	RegisterCLI("exit", &struct{}{}, func(command string, data interface{}) {
		response := WaitForConsoleTimeout("Do you want exit? (Y/N)", 5*time.Second, func() {
			fmt.Println("Confirmation Timed out ... back to normal ...")
		})
		response = strings.TrimSpace(strings.ToLower(response))
		if response == "y" || response == "yes" {
			for _, item := range onShutdownCallbacks {
				item()
			}
			os.Exit(2)
		}
	}, "Exit app")

	RegisterCLI("config.get", &struct{}{}, func(command string, data interface{}) {
		spew.Dump(GetConfig())
	}, "Show config")

	RegisterCLI("log.set", &cfg.Log, func(command string, data interface{}) {
		cfg := data.(*Log)
		log.SetSettings(&log.Logger{
			WriteToFile: cfg.WriteFile,
			Concurrent:  true,
			Level:       log.ParseLevel(cfg.Level),
			MaxAge:      cfg.MaxAge,
			MaxSize:     cfg.MaxSize,
			Path:        cfg.Path,
		})
		fmt.Println("Log level set to:" + cfg.Level)
	}, "Set log file verbosity")

	type logLines struct {
		Lines int    `short:"l" long:"lines" description:"maximum lines of logs to show" default:"100"`
		Level string `short:"v" long:"level" description:"maximum level of logs to show" choice:"critical" choice:"error" choice:"warning" choice:"info" choice:"notice" choice:"debug" default:"debug"`
		Test  string
	}
	logLineParam := &logLines{}
	RegisterCLI("log.read", logLineParam, func(command string, data interface{}) {
		log.Read(data.(*logLines).Lines, log.ParseLevel(data.(*logLines).Level))
	}, "Read log file")

	RegisterCLI("log.clear", &struct{}{}, func(command string, data interface{}) {
		fmt.Println("Log cleared ...")
		log.Clear()
	}, "Clear current day log")
	RegisterCLI("log.clearall", &struct{}{}, func(command string, data interface{}) {
		fmt.Println("All logs has been cleared ...")
		log.ClearAll()
	}, "Clear all logs")
}

func Essentials() {
	language.Register()
	docker.Register()
}
