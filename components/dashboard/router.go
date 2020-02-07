package dashboard

import (
	"github.com/iesreza/foundation/lib/form"
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
		form.Unmarshal(req, &data)

	})

	system.Router.Match("hello", "GET", func(req router.Request) {
		req.RemoveCookie("test")
		View{}.Hello(req)
	})
	system.Router.Static(component.Assets, component.Assets, nil)

}
