package system

import (
	"github.com/iesreza/foundation/lib"
	"github.com/iesreza/foundation/lib/gpath"
	"github.com/iesreza/foundation/lib/log"
	"github.com/iesreza/foundation/lib/request"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
	"runtime"
)

var configInstance *Config
var Debug bool

type Log struct {
	WriteFile bool   `yaml:"write-file"`
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
		BallastSize   string `yaml:"ballast-size"`
		MaxProcessors int    `yaml:"processors"`
	} `yaml:"tweaks"`
	Log Log `yaml:"log"`
	App struct {
		WorkingDir    string
		OS            string
		ProcessID     int
		LogoMini      string `yaml:"logo-mini"`
		LogoLarge     string `yaml:"logo-large"`
		Title         string `yaml:"title"`
		Path          string `yaml:"path"`
		Assets        string `yaml:"assets"`
		Static        string `yaml:"static"`
		SessionAge    int    `yaml:"session-age"`
		Language      string `yaml:"language"`
		MaxUploadSize string `yaml:"max-upload-size"`
		Debug         bool   `yaml:"debug"`
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
		CacheSize string `yaml:"cache-size"`
		Debug     string `yaml:"debug"`
		Database  string `yaml:"database"`
		SSLMode   string `yaml:"ssl-mode"`
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

		request.MaxUploadSize, _ = lib.ParseSize(configInstance.App.MaxUploadSize)
		Debug = configInstance.App.Debug
	}
	return *configInstance
}

func LoadConfig(path string, out interface{}) {
	if reflect.TypeOf(out).String()[0] != '*' {
		log.Critical("Passed object to system.LoadConfig is not pointer")
		return
	}
	var err error
	var f *os.File
	if path == "" || gpath.IsFileExist(path) {
		f, err = os.Open(path)
	} else {
		f, err = os.Open("./config.yml")
	}
	if err != nil {
		log.Critical(err)
		return
	}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(out)
	if err != nil {
		log.Critical(err)
		return
	}

}
