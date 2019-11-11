package system

import "github.com/iesreza/foundation/lib"

func Icon(icon string) string  {
	if lib.String(icon).EndsWith("fa-"){
		return "<i class=\"fa "+icon+"\"></i>"
	}
	return icon
}