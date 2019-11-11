package system

import "html/template"

var Components = map[string]Component{}


type Component interface {
	Register()
	Routers()
	Menu()
	Install()
	Uninstall()
	Update()
	ComputeHash()
	ViewPath() string
	AssetsPath() string
	GetTemplates() *template.Template
}
