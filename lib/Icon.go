package lib

func Icon(icon string) string {
	if String(icon).EndsWith("fa-") {
		return "<i class=\"fa " + icon + "\"></i>"
	}
	return icon
}
