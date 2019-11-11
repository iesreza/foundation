package system

import (
	"github.com/iesreza/foundation/lib/router"
	"strings"
)

var Router = router.GetInstance()
var Request router.Request
func SetupRouter()  {
	Router.Fallback = func(req router.Request) {
		View(&req).SetLayout("error").SetBody("Page not found error 404").Render().Write()
	}
	Router.Static("assets", config.App.Assets,nil)
	Router.Static("static", config.App.Static,nil)
}

func Route(route string) string  {
	if !strings.HasPrefix(route,"http") {
		return Request.Req().URL.Scheme + "://" + Request.Req().Host + "/" + strings.Trim(route, "/ ")
	}
	return route
}
