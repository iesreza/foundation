package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/iesreza/foundation/lib/router"
	"github.com/iesreza/gutil/path"
	"html/template"
)

var Assets = GetConfig().App.Assets
var Templates = map[string]*Template{}

type TemplateInterface interface {
	Register()
	Menu() map[string]*Menu
	Path() string
	Assets() string
	Static() map[string]string
	Name() string
}

type Template struct {
	Name        string
	Version     string
	Path        string
	Assets      string
	Templates   *template.Template
	StyleSheets []CSS
	Scripts     []JS
	Interface   *TemplateInterface
}

type CSS struct {
	Path    string
	Attribs []string
}

type JS struct {
	Path    string
	Attribs []string
}

func (t Template) RenderLayout(request *router.Request, layout string, data map[string]interface{}) []byte {
	buf := new(bytes.Buffer)
	data["App"] = config.App
	for key, menu := range (*t.Interface).Menu() {
		data[key] = menu.Render(request)
	}
	t.Templates.ExecuteTemplate(buf, layout+".html", data)
	return buf.Bytes()
}

func (t Template) CSS(path string, attribs ...string) {
	t.StyleSheets = append(t.StyleSheets, CSS{
		Path:    path,
		Attribs: attribs,
	})
}

func (t Template) JS(path string, attribs ...string) {
	t.Scripts = append(t.Scripts, JS{
		Path:    path,
		Attribs: attribs,
	})
}

func SetDefaultTemplate(ti TemplateInterface) (Template, error) {
	t, err := LoadTemplate(ti)
	Templates["_DEFAULT_"] = &t
	return t, err
}

func LoadTemplate(ti TemplateInterface) (Template, error) {
	t := Template{}
	dir := path.Dir(ti.Path())
	if !dir.Exist() {
		return t, fmt.Errorf("unable to find template at %s", ti.Name())
	}
	templateJson := dir.File("template.json")
	if !templateJson.Exist() {
		return t, fmt.Errorf("unable to find template.json at %s", ti.Path()+"/template.json")
	}
	data, err := templateJson.Content()
	if err != nil {
		return t, fmt.Errorf("unable to read template.json at %s", ti.Path()+"/template.json")
	}
	err = json.Unmarshal([]byte(data), &t)
	if err != nil {
		return t, fmt.Errorf("unable to parse template.json at %s", ti.Path()+"/template.json")
	}

	files, err := dir.Find("*.html")
	if err != nil {
		return t, fmt.Errorf("unable to parse template html files %s", err.Error())
	}

	t.Templates, err = template.ParseFiles(files...)
	if err != nil {
		return t, fmt.Errorf("unable to parse template html layouts. %s", err.Error())
	}

	t.Path = ti.Path()
	t.Assets = ti.Assets()
	t.Interface = &ti
	Router.Static(ti.Name()+"/assets", ti.Assets(), nil)
	if len(ti.Static()) > 0 {
		for route, path := range ti.Static() {
			Router.Static(route, ti.Path()+"/"+path, nil)
		}
	}
	Templates[ti.Name()] = &t
	return t, nil
}

func GetTemplate() *Template {
	return Templates["_DEFAULT_"]
}

func FindTemplate(t string) *Template {
	return Templates[t]
}
