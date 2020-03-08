package system

import (
	"bytes"
	"github.com/iesreza/foundation/lib/request"
	"html/template"
	"strings"
)

type view struct {
	Layout         string
	Body           template.HTML
	Data           map[string]interface{}
	Rendered       []byte
	isRendered     bool
	isDataPrepared bool
	Template       *Template
	Request        *request.Request
	StyleSheets    []CSS
	Scripts        []JS
}

func (v view) CSS(dir string, path string, attribs ...string) {
	v.StyleSheets = append(v.StyleSheets, CSS{
		Path:    path,
		Attribs: attribs,
	})
}

func (v view) JS(dir string, path string, attribs ...string) {
	v.Scripts = append(v.Scripts, JS{
		Path:    path,
		Attribs: attribs,
	})
}

func (v *view) SetLayout(l string) *view {
	v.Layout = l
	return v
}

func (v *view) SetBody(body string) *view {
	v.Body = template.HTML(body)
	return v
}

func (v *view) SetData(data map[string]interface{}) *view {
	v.Data = data
	return v
}

func (v *view) SetTemplate(template string) *view {
	v.Template = FindTemplate(template)
	return v
}

func (v *view) AddVariable(data ...map[string]interface{}) *view {
	for _, item := range data {
		for key, val := range item {
			v.Data[key] = val
		}
	}
	return v
}

func (v *view) Render() *view {

	v.isRendered = true
	v.prepareData([]map[string]interface{}{})
	v.Data["Body"] = v.Body
	v.Rendered = v.Template.RenderLayout(v.Request, v.Layout, v.Data)
	return v
}

func (v *view) Write() *view {
	if !v.isRendered {
		v.Render()
	}
	v.Request.Write(v.Rendered)
	return v
}

func (v *view) Call(c Component, view string, args ...map[string]interface{}) *view {
	c.ViewPath()
	buf := new(bytes.Buffer)
	v.prepareData(args)
	c.GetTemplates().ExecuteTemplate(buf, view+".html", v.Data)
	v.Body = template.HTML(buf.String())
	return v
}

func (v view) prepareData(args []map[string]interface{}) {
	if v.isDataPrepared {
		return
	}
	v.Data["Request"] = *v.Request.Req()
	post := map[string]interface{}{}
	for key, value := range v.Request.Req().PostForm {
		post[key] = value
	}
	v.Data["POST"] = post
	v.Data["GET"] = v.Request.Req().URL.Query()
	v.Data["Parameter"] = v.Request.Parameters

	for _, item := range args {
		for key, val := range item {
			v.Data[key] = val
		}
	}

	v.Data["Scripts"] = ""
	v.Data["Stylesheets"] = ""
	s := ""
	for _, item := range v.Scripts {
		s += "<script language='javascript' " + strings.Join(item.Attribs, " ") + " src='" + item.Path + "></script>"
	}
	v.Data["Scripts"] = s
	s = ""
	for _, item := range v.StyleSheets {
		s += "<link rel='stylesheet' " + strings.Join(item.Attribs, " ") + " href='" + item.Path + " />"
	}
	v.Data["Stylesheets"] = s
	v.isDataPrepared = true
}

func View(req *request.Request) *view {
	return &view{
		Layout:   "index",
		Template: GetTemplate(),
		Data:     map[string]interface{}{},
		Request:  req,
	}
}

func Redirect() {

}

/*func InternalError(error string)  {
	View("error")
}
*/
