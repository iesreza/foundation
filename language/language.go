package language

import (
	"encoding/json"
	"fmt"
	"github.com/iesreza/foundation/lib/log"
	"github.com/iesreza/foundation/system"
	"github.com/iesreza/gutil/path"
	"io/ioutil"
	"strconv"
	"strings"
)

var Available []*Translator
var Default *Translator
var DIR = system.DIR + "/language"

type Translator struct {
	Name      string `json:"Name"`
	Flag      string `json:"Flag"`
	LocalName string `json:"LocalName"`
	Default   bool
	Table     map[string]string
}

func Register() {
	//load languages
	LoadLanguages()
	system.RegisterCLI("language.default", &struct{}{}, func(command string, data interface{}) {
		fmt.Println(Default.Name)
	}, "Show Default Language")

	system.RegisterCLI("language.available", &struct{}{}, func(command string, data interface{}) {
		for i, item := range Available {
			fmt.Printf("%d. %s\r\n", i+1, item.Name)
		}
	}, "Show available languages")

}

func LoadLanguages() {
	files, err := ioutil.ReadDir(DIR)
	if err != nil {
		log.Critical(err)
	}
	for _, item := range files {
		if item.IsDir() {
			translator := &Translator{
				Table: map[string]string{},
			}
			readLangInfo(item.Name(), translator)
			parseLangFiles(item.Name(), translator)
			Available = append(Available, translator)
			if translator.Name == system.GetConfig().App.Language {
				translator.SetDefault()
			}
		}
	}

}

func parseLangFiles(lang string, translator *Translator) {
	files, err := path.Dir(DIR + "/" + lang).Find("*.json")
	if err != nil {
		log.Critical(err)
	}
	for _, item := range files {
		if !strings.HasSuffix(item, "language.json") {
			data, err := path.File(item).Content()
			if err != nil {
				log.Critical(err)
			}
			structure := map[string]interface{}{}
			err = json.Unmarshal([]byte(data), &structure)
			if err != nil {
				log.Critical(fmt.Errorf("Malformed Language File:" + item))
			}
			parseLangStructure("", structure, translator)
		}
	}

}

func parseLangStructure(key string, structure map[string]interface{}, translator *Translator) {
	for k, v := range structure {
		nkey := strings.Trim(key+"."+k, ".")
		switch v := v.(type) {
		case int:
			translator.Table[nkey] = strconv.Itoa(v)
		case string:
			translator.Table[nkey] = v
		case map[string]interface{}:
			parseLangStructure(nkey, v, translator)

		}

	}
}

func readLangInfo(lang string, translator *Translator) {
	data, err := path.File(DIR + "/" + lang + "/language.json").Content()
	if err != nil {
		log.Critical(err)
	}
	err = json.Unmarshal([]byte(data), translator)
	if err != nil {
		log.Critical(fmt.Errorf("Malformed Language File:" + DIR + "/" + lang + "/language.json"))
	}
}

func T(msg string, params ...interface{}) string {
	return Default.T(msg, params...)
}

func (t *Translator) SetDefault() {
	for _, item := range Available {
		item.Default = false
	}
	t.Default = true
	Default = t
}

func (t *Translator) T(msg string, params ...interface{}) string {
	if t.Table == nil {
		return msg
	}

	if item, ok := t.Table[msg]; ok {
		return fmt.Sprintf(item, params...)
	}
	return msg
}
