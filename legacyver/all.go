package legacyver

import (
	"github.com/sandertv/gophertunnel/minecraft"
)

// All returns a slice of all legacy protocol versions that are supported. deleteDebugStick
// must be set to true if you're using Dragonfly. deleteDebugStick must be set to false if you're using
// PocketMine, Nukkit, or other server software that includes the debug stick in the item registry.
func All(deleteDebugStick bool) []minecraft.Protocol {
	return []minecraft.Protocol{
		New800(deleteDebugStick),
		New786(deleteDebugStick),
		New776(deleteDebugStick),
		New766(deleteDebugStick),
		New748(deleteDebugStick),
		New729(deleteDebugStick),
		New712(deleteDebugStick),
		New686(deleteDebugStick),
		New685(deleteDebugStick),
		New671(deleteDebugStick),
		New662(deleteDebugStick),
		New649(deleteDebugStick),
	}
}
