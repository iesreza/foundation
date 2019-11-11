package system

import (
	"html/template"
	"strconv"
	"strings"
)

type Menu struct {
	SubMenu    []Menu
	Name        string
	Title       string
	Permission  string
	URL         string
	Icon        string
	Class       string
	ID          string
	ParentClass string
}

func (m *Menu)Push(menu... Menu){
	for _,item := range menu{
		m.SubMenu = append(m.SubMenu,item)
	}

}

func (m Menu)Render(attribs... string)  template.HTML {
	html := "<ul class=\""+m.ParentClass+"\" "+strings.Join(attribs," ")+">\n"

	for _,item := range m.SubMenu{

		t,_:= recursiveMenuRender(&item,1)

			html += "\t" + t + "\n"

	}

	html += "</ul>\n"
	return template.HTML(html)
}

func recursiveMenuRender(m *Menu,depth int) (string,bool)  {
	if !GetUser().HasPerm(m.Permission){
		return "",false
	}

	//Menu creation event fire
	for _,item := range onMenuRenderCallbacks{
		item(m)
	}
	html := "<li"
	if m.ID != ""{
		html += " id=\""+m.ID+"\""
	}
	if m.Class != ""{
		html += " class=\""+m.Class+"\""
	}

	hasChild := false
	temp := ""
	if len(m.SubMenu) > 0{
		for i:=0; i < depth+1; i++ {
			temp += "\t"
		}
		temp += "<ul class=\"child depth-"+strconv.Itoa(depth+1)+" "+m.ParentClass+"\">\n"
		for _,item := range m.SubMenu{
			t,p := recursiveMenuRender(&item,depth+1)
			if p{
				for i:=0; i < depth+2; i++ {
					temp += "\t"
				}
				temp += t
				hasChild = true
			}
		}
		for i:=0; i < depth+1; i++ {
			temp += "\t"
		}
		temp += "</ul>\n"

	}



	html += ">"
	if hasChild {
		html += "<a href=\"#\">"
	}else{
		html += "<a href=\""+m.URL+"\">"
	}
	if m.Icon != ""{
		html += Icon(m.Icon)
	}
	html += m.Title
	html += "</a>"

	if hasChild {
		html += "\n"+temp
	}
	for i:=0; i < depth; i++ {
		html += "\t"
	}
	html += "</li>\n"

	return html,true
}

