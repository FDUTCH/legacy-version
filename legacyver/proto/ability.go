package proto

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// AbilityData represents various data about the abilities of a player, such as ability layers or permissions.
type AbilityData struct {
	// EntityUniqueID is a unique identifier of the player. It appears it is not required to fill this field
	// out with a correct value. Simply writing 0 seems to work.
	EntityUniqueID int64
	// PlayerPermissions is the permission level of the player as it shows up in the player list built up using
	// the PlayerList packet.
	PlayerPermissions byte
	// CommandPermissions is a set of permissions that specify what commands a player is allowed to execute.
	CommandPermissions byte
	// Layers contains all ability layers and their potential values. This should at least have one entry, being the
	// base layer.
	Layers []AbilityLayer
}

// Marshal encodes/decodes an AbilityData.
func (x *AbilityData) Marshal(r protocol.IO) {
	r.Int64(&x.EntityUniqueID)
	r.Uint8(&x.PlayerPermissions)
	r.Uint8(&x.CommandPermissions)
	protocol.SliceUint8Length(r, &x.Layers)
}

// ToLatest ...
func (x *AbilityData) ToLatest() protocol.AbilityData {
	layers := make([]protocol.AbilityLayer, len(x.Layers))
	for i, layer := range x.Layers {
		layers[i] = layer.ToLatest()
	}
	return protocol.AbilityData{
		EntityUniqueID:     x.EntityUniqueID,
		PlayerPermissions:  x.PlayerPermissions,
		CommandPermissions: x.CommandPermissions,
		Layers:             layers,
	}
}

// FromLatest ...
func (x *AbilityData) FromLatest(y protocol.AbilityData) AbilityData {
	layers := make([]AbilityLayer, len(y.Layers))
	for i, layer := range y.Layers {
		layers[i] = (&AbilityLayer{}).FromLatest(layer)
	}
	x.EntityUniqueID = y.EntityUniqueID
	x.PlayerPermissions = y.PlayerPermissions
	x.CommandPermissions = y.CommandPermissions
	x.Layers = layers
	return *x
}

// AbilityLayer represents the abilities of a specific layer, such as the base layer or the spectator layer.
type AbilityLayer struct {
	// Type represents the type of the layer. This is one of the AbilityLayerType constants defined above.
	Type uint16
	// Abilities is a set of abilities that are enabled for the layer. This is one of the Ability constants defined
	// above.
	Abilities uint32
	// Values is a set of values that are associated with the enabled abilities, representing the values of the
	// abilities.
	Values uint32
	// FlySpeed is the default horizontal fly speed of the layer.
	FlySpeed float32
	// VerticalFlySpeed is the default vertical fly speed of the layer.
	VerticalFlySpeed float32
	// WalkSpeed is the default walk speed of the layer.
	WalkSpeed float32
}

// Marshal encodes/decodes an AbilityLayer.
func (x *AbilityLayer) Marshal(r protocol.IO) {
	r.Uint16(&x.Type)
	r.Uint32(&x.Abilities)
	r.Uint32(&x.Values)
	r.Float32(&x.FlySpeed)
	if IsProtoGTE(r, ID776) {
		r.Float32(&x.VerticalFlySpeed)
	}
	r.Float32(&x.WalkSpeed)
}

// ToLatest ...
func (x *AbilityLayer) ToLatest() protocol.AbilityLayer {
	return protocol.AbilityLayer{
		Type:             x.Type,
		Abilities:        x.Abilities,
		Values:           x.Values,
		FlySpeed:         x.FlySpeed,
		VerticalFlySpeed: x.VerticalFlySpeed,
		WalkSpeed:        x.WalkSpeed,
	}
}

// FromLatest ...
func (x *AbilityLayer) FromLatest(y protocol.AbilityLayer) AbilityLayer {
	x.Type = y.Type
	x.Abilities = y.Abilities
	x.Values = y.Values
	x.FlySpeed = y.FlySpeed
	x.VerticalFlySpeed = y.VerticalFlySpeed
	x.WalkSpeed = y.WalkSpeed
	return *x
}
