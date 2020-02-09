package main

import (
	"github.com/iesreza/foundation/components/dashboard"
	"github.com/iesreza/foundation/httpdocs/templates/adminlte"
	"github.com/iesreza/foundation/lib/log"
	"github.com/iesreza/foundation/system"
)

func main() {
	system.PreBoot()

	//Register essential components (docker,language)
	system.Essentials()

	//Register user components
	adminlte.Register()
	dashboard.Register()

	system.Boot()
	log.Info("Foundation Has Booted")
	system.ListenCLI()

}
