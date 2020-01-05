package main

import (
	"github.com/iesreza/foundation/components/dashboard"
	"github.com/iesreza/foundation/httpdocs/templates/adminlte"
	"github.com/iesreza/foundation/language"
	"github.com/iesreza/foundation/log"
	"github.com/iesreza/foundation/system"
)

func main() {
	system.PreBoot()

	adminlte.Register()
	dashboard.Register()
	language.Register()

	system.Boot()
	log.Info("Foundation Has Booted")
	system.ListenCLI()

}
