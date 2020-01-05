package system

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/iesreza/foundation/lib"
	"github.com/iesreza/foundation/lib/router"
	"github.com/iesreza/foundation/log"
	"github.com/ztrue/shutdown"
	"os"
	"runtime"
	"time"
)

func Boot() {

	SetupDatabase()
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
		log.Warning("Shutting Down ...")
		//Shutdown event fire
		for _, item := range onShutdownCallbacks {
			item()
		}
	})

	RegisterCLI("help", &struct{}{}, func(command string, data interface{}) {
		for command, _ := range commandList {
			fmt.Println(command)
		}
	})
	RegisterCLI("exit", &struct{}{}, func(command string, data interface{}) {
		for _, item := range onShutdownCallbacks {
			item()
		}
		os.Exit(2)
	})

	RegisterCLI("config.get", &struct{}{}, func(command string, data interface{}) {
		spew.Dump(GetConfig())
	})

}
