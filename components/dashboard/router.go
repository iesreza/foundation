package dashboard

import (
	"fmt"
	"github.com/iesreza/foundation/lib/log"
	"github.com/iesreza/foundation/lib/router"
	"github.com/iesreza/foundation/system"
)

func (component component) Routers() {

	system.Router.Match("test", "POST", func(req router.Request) {

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

	system.Router.Match("hello", "GET", func(req router.Request) {
		req.RemoveCookie("test")
		View{}.Hello(req)
	})
	system.Router.Static(component.Assets, component.Assets, nil)

}
