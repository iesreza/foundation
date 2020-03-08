package admin

import (
	"github.com/iesreza/foundation/lib/request"
	"github.com/iesreza/foundation/system"
)

type View struct{}

func (v View) Hello(req request.Request) {
	view := system.View(&req)
	view.CSS(Component.Assets, "test.css")
	view.JS(Component.Assets, "js.css", "async=true", "defer=js")
	view.Call(&Component, "test").Write()
}
