package docker

import (
	"github.com/iesreza/foundation/system"
	"html/template"
)

type component struct {
	Path      string
	Assets    string
	Views     string
	Templates *template.Template
}

func (component component) Routers() {

}

func (component component) GetTemplates() *template.Template {
	return component.Templates
}

func (component component) ViewPath() string {
	return component.Views
}

func (component component) AssetsPath() string {
	return component.Assets
}

var Component = component{}

func Register() {
	Component.Register()
	RegisterCLI()
}

func (component *component) Register() {
	system.Components["docker"] = component

}

func (component *component) Menu() {

}

func (component *component) Install() {

}

func (component *component) Uninstall() {

}

func (component *component) Update() {

}

func (component *component) ComputeHash() {
	panic("implement me")
}
