package mode

import (
	"strings"
)

type Mode string

const (
	Recurring Mode = "Recurring"
	Singular  Mode = "Singular"
	Unknown   Mode = "Unknown"
)

func Atos(a string) Mode {
	switch strings.ToLower(a) {
	case "recurring":
		return Recurring
	case "singular":
		return Singular
	default:
		return Unknown
	}
}
