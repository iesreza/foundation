package main

import (
	"github.com/iesreza/foundation/components/dashboard"
	"github.com/iesreza/foundation/httpdocs/templates/adminlte"
	"github.com/iesreza/foundation/system"
	"github.com/iesreza/gutil/log"
	"time"
)


func main()  {

	log.Error("hi")
	adminlte.Register()
	dashboard.Register()
	system.Boot()


	for{
		time.Sleep(1*time.Minute)
	}
}
