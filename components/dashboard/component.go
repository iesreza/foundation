package dashboard

import (
	"fmt"
	"github.com/iesreza/foundation/httpdocs/templates/adminlte"
	"github.com/iesreza/foundation/system"
	"github.com/iesreza/gutil/path"
	"html/template"
)

type component struct {
	Path string
	Assets string
	Views string
	Templates *template.Template
}

func (component component) GetTemplates() *template.Template {
	return component.Templates
}

func (component component) ViewPath() string{
	return component.Views
}

func (component component) AssetsPath() string{
	return component.Assets
}

var Component = component{
	Path:   "components/dashboard/",
	Assets: "components/dashboard/assets/",
	Views:  "components/dashboard/views/",
}


func Register()  {
	Component.Register()

}

func (component *component) Register() {
	system.Components["dashboard"] = component
	files,err := path.Dir(component.Views).Find("*.html")
	if err != nil{
		system.Critical(fmt.Errorf("unable to parse template html files %s",err.Error()))
	}

	component.Templates, err = template.ParseFiles(files...)
	if err != nil{
		system.Critical(fmt.Errorf("unable to parse template html layouts. %s",err.Error()))
	}
}



func (component *component) Menu() {
	adminlte.MainNav.Push(
		system.Menu{
		Name:    "MainMenu",Title:   "Home",Permission:"",URL: "dashboard",Icon:"fa-home",Class:"home",
		SubMenu: []system.Menu{
			system.Menu{
				Name:    "MainMenu",Title:   "->1",Permission:"",URL: "dashboard",Icon:"fa-home",Class:"home",
			},
			system.Menu{
				Name:    "MainMenu",Title:   "->2",Permission:"",URL: "dashboard",Icon:"fa-home",Class:"home",
			},
			system.Menu{
				Name:    "MainMenu",Title:   "->3",Permission:"",URL: "dashboard",Icon:"fa-home",Class:"home",
			},
	}} ,
		system.Menu{
			Name:    "MainMenu2",Title:   "Home2",Permission:"",URL: "dashboard",Icon:"fa-home",Class:"home",
			SubMenu: []system.Menu{
				system.Menu{
					Name:    "MainMenu",Title:   "->1",Permission:"",URL: "dashboard",Icon:"fa-home",Class:"home",
				},
				system.Menu{
					Name:    "MainMenu",Title:   "->2",Permission:"",URL: "dashboard",Icon:"fa-home",Class:"home",
				},
				system.Menu{
					Name:    "MainMenu",Title:   "->3",Permission:"",URL: "dashboard",Icon:"fa-home",Class:"home",
				},
			}},
	)
}

func (component *component) Install() {
	panic("implement me")
}

func (component *component) Uninstall() {
	panic("implement me")
}

func (component *component) Update() {
	panic("implement me")
}

func (component *component) ComputeHash() {
	panic("implement me")
}




