package adminlte

import (
	"fmt"
	"github.com/iesreza/foundation/system"
)

var MainNav = system.Menu{
	SubMenu: []system.Menu{},
	ParentClass:"main-nav",
}

func Register(){
	_,err := system.SetDefaultTemplate(Template{})
	if err != nil {
		fmt.Println(err)
	}
}
type Template struct {

}

func (Template) Register() {

}

func (Template) Menu() map[string]*system.Menu {
	return map[string]*system.Menu{
		"MainNavigation":&MainNav,
	}
}

func (Template) Path() string {
	return "./httpdocs/templates/adminlte"
}

func (Template) Assets() string {
	return "./httpdocs/templates/adminlte/assets"
}

func (Template) Static() map[string]string {
	return map[string]string{
		"mystatic":"static/static.html",
	}
}

func (Template) Name() string {
	return "adminlte"
}

