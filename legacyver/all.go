package legacyver

import (
	"github.com/sandertv/gophertunnel/minecraft"
)

// All returns a slice of all legacy protocol versions that are supported. dragonflyMapping
// must be set to true if you're using Dragonfly.
func All(dragonflyMapping bool) []minecraft.Protocol {
	return []minecraft.Protocol{
		New819(dragonflyMapping),
		New818(dragonflyMapping),
		New800(dragonflyMapping),
		New786(dragonflyMapping),
		New776(dragonflyMapping),
		New766(dragonflyMapping),
		New748(dragonflyMapping),
		New729(dragonflyMapping),
		New712(dragonflyMapping),
		New686(dragonflyMapping),
		New685(dragonflyMapping),
		New671(dragonflyMapping),
		New662(dragonflyMapping),
		New649(dragonflyMapping),
	}
}
