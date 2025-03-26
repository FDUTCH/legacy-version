package legacyver

import "github.com/sandertv/gophertunnel/minecraft"

// All returns a slice of all legacy protocol versions that are supported.
func All() []minecraft.Protocol {
	return []minecraft.Protocol{
		New776(),
		New766(),
		New748(),
		New729(),
		New712(),
		New686(),
		New685(),
		New671(),
	}
}
