package system

import (
	"gopkg.in/yaml.v2"
	"os"
	"runtime"
)
var configInstance *Config
type Config struct {
	Alarm struct{
		Processor float64 `yaml:"processor"`
		Memory float64 `yaml:"memory"`
		Reset  bool `yaml:"reset"`
	} `yaml:"alarm"`
	App struct{
		WorkingDir    string
		OS            string
		MaxProcessors int `yaml:"processors"`
		ProcessID    int
		LogoMini   string `yaml:"logomini"`
		LogoLarge  string `yaml:"logolarge"`
		Title      string `yaml:"title"`
		Path       string `yaml:"path"`
		Assets     string `yaml:"assets"`
		Static     string `yaml:"static"`
		SessionAge int    `yaml:"sessionage"`
	} `yaml:"app"`
	Server struct {
		Port string  `yaml:"port"`
		Host string  `yaml:"host"`
		Cert string  `yaml:"cert"`
		Key  string  `yaml:"key"`
		HTTPS bool `yaml:"https"`
	} `yaml:"server"`
	Database struct {
		Type string `yaml:"type"`
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
		Server string `yaml:"server"`
		Cache string `yaml:"cache"`
		CacheSize string `yaml:"cachesize"`
		Debug string `yaml:"debug"`
		Database string `yaml:"database"`
		SSLMode string `yaml:"sslmode"`
	} `yaml:"database"`
}

func GetConfig() Config {
	if configInstance == nil{
		configInstance = &Config{}
		f, err := os.Open("./config.yml")
		if err != nil {
			Critical(err)
		}
		decoder := yaml.NewDecoder(f)
		err = decoder.Decode(&configInstance)
		if err != nil {
			Critical(err)
		}

		configInstance.App.WorkingDir, err = os.Getwd()
		if err != nil {
			Critical(err)
		}
		configInstance.App.OS = runtime.GOOS
		configInstance.App.ProcessID = os.Getpid()

		if configInstance.App.MaxProcessors < 1{
			runtime.GOMAXPROCS(runtime.NumCPU())
		}else{
			runtime.GOMAXPROCS(configInstance.App.MaxProcessors)
		}
	}
	return *configInstance
}
