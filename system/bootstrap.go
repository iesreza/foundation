package system

import (
	"github.com/iesreza/foundation/lib/router"
	"time"
)

func Boot()  {

	SetupDatabase()
	SetupRouter()
	router.SessionAge = time.Duration(config.App.SessionAge)*time.Second
	StartWebServer()

	for _,item := range Components{
		item.Menu()
		item.Routers()
	}
}
