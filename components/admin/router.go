package admin

import (
	"encoding/base64"
	"fmt"
	"github.com/iesreza/foundation/lib/log"
	"github.com/iesreza/foundation/lib/request"
	"github.com/iesreza/foundation/system"
	"net/http"
)

var UnAuthorized = []byte(`{"Success":false,"Error":"Unauthorized"}`)

func (component component) Routers() {

	system.Router.Match("test", "POST", func(req request.Request) {

		data := struct {
			X string
			Y int
			Z bool
		}{}
		req.Unmarshal(&data)
		fmt.Println(req.Form.Get("mykey"))
		fmt.Println(data)
		if _, ok := req.Files["myfile"]; ok {
			fmt.Println(req.Files["myfile"][0].Move("d:/" + req.Files["myfile"][0].Name))
		}
		log.Notice(req.Query.Get("vvv"))
	})

	system.Router.Match("hello", "GET", func(req request.Request) {
		req.RemoveCookie("test")
		View{}.Hello(req)
	})

	system.Router.Group("admin", nil, func(handle *request.Route) {
		handle.Match("login", "POST", func(req request.Request) {
			controller := AuthController{}
			controller.login(req)
		})
		handle.Match("login", "GET", func(req request.Request) {
			req.WriteString("you are here")
		})

	}).Middleware(func(req request.Request) bool {

		if req.GetUser().IsGuest() && req.Path != "/admin/login" {
			if req.Req().Method == "GET" {
				url := req.Req().URL
				redirect := req.Path
				if url.RawQuery != "" {
					redirect += "?" + url.RawQuery
				}
				req.Redirect("/admin/login?redirect=" + base64.StdEncoding.EncodeToString([]byte(redirect)))
			} else {
				w := *req.Writer()
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(UnAuthorized)
				req.Terminate()
			}
			return true
		}

		return true
	})

	system.Router.Static(component.Assets, component.Assets, nil)

}
