package system

import (
	"github.com/iesreza/foundation/lib/log"
	"gopkg.in/yaml.v2"
	"os"
	"runtime"
)

var configInstance *Config

type Log struct {
	WriteFile bool   `yaml:"writefile"`
	MaxSize   int    `yaml:"size"` // megabytes
	MaxAge    int    `yaml:"age"`  //days
	Level     string `yaml:"level" short:"l" long:"level" description:"log level" choice:"critical" choice:"error" choice:"warning" choice:"info" choice:"notice" choice:"debug" default:"debug"`
	Path      string `yaml:"path"`
}
type Config struct {
	Alarm struct {
		Processor float64 `yaml:"processor"`
		Memory    float64 `yaml:"memory"`
		Reset     bool    `yaml:"reset"`
	} `yaml:"alarm"`
	Tweaks struct {
		Ballast       bool   `yaml:"ballast"`
		BallastSize   string `yaml:"ballastsize"`
		MaxProcessors int    `yaml:"processors"`
	} `yaml:"tweaks"`
	Log Log `yaml:"log"`
	App struct {
		WorkingDir string
		OS         string
		ProcessID  int
		LogoMini   string `yaml:"logomini"`
		LogoLarge  string `yaml:"logolarge"`
		Title      string `yaml:"title"`
		Path       string `yaml:"path"`
		Assets     string `yaml:"assets"`
		Static     string `yaml:"static"`
		SessionAge int    `yaml:"sessionage"`
		Language   string `yaml:"language"`
	} `yaml:"app"`
	Server struct {
		Port  string `yaml:"port"`
		Host  string `yaml:"host"`
		Cert  string `yaml:"cert"`
		Key   string `yaml:"key"`
		HTTPS bool   `yaml:"https"`
	} `yaml:"server"`
	Database struct {
		Enabled   bool   `yaml:"enabled"`
		Type      string `yaml:"type"`
		Username  string `yaml:"user"`
		Password  string `yaml:"pass"`
		Server    string `yaml:"server"`
		Cache     string `yaml:"cache"`
		CacheSize string `yaml:"cachesize"`
		Debug     string `yaml:"debug"`
		Database  string `yaml:"database"`
		SSLMode   string `yaml:"sslmode"`
		Params    string `yaml:"params"`
	} `yaml:"database"`
}

func GetConfig() Config {
	if configInstance == nil {
		configInstance = &Config{}
		f, err := os.Open("./config.yml")
		if err != nil {
			log.Critical(err)
		}
		decoder := yaml.NewDecoder(f)
		err = decoder.Decode(&configInstance)
		if err != nil {
			log.Critical(err)
		}

		configInstance.App.WorkingDir, err = os.Getwd()
		if err != nil {
			log.Critical(err)
		}
		configInstance.App.OS = runtime.GOOS
		configInstance.App.ProcessID = os.Getpid()

	}
	return *configInstance
}
