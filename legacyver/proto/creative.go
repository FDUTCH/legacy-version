package proto

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// CreativeItem represents a creative item present in the creative inventory.
type CreativeItem struct {
	// CreativeItemNetworkID is a unique ID for the creative item. It has to be unique for each creative item
	// sent to the client. An incrementing ID per creative item does the job.
	CreativeItemNetworkID uint32
	// Item is the item that should be added to the creative inventory.
	Item protocol.ItemStack
	// GroupIndex is the index of the group that the item should be placed in. It is the index of the group in
	// the CreativeContent packet previously sent to the client.
	GroupIndex uint32
}

// ToLatest ...
func (x *CreativeItem) ToLatest() protocol.CreativeItem {
	return protocol.CreativeItem{
		CreativeItemNetworkID: x.CreativeItemNetworkID,
		Item:                  x.Item,
		GroupIndex:            x.GroupIndex,
	}
}

// FromLatest ...
func (x *CreativeItem) FromLatest(y protocol.CreativeItem) CreativeItem {
	x.CreativeItemNetworkID = y.CreativeItemNetworkID
	x.Item = y.Item
	x.GroupIndex = y.GroupIndex
	return *x
}

// Marshal encodes/decodes a CreativeItem.
func (x *CreativeItem) Marshal(r protocol.IO) {
	r.Varuint32(&x.CreativeItemNetworkID)
	r.Item(&x.Item)
	if IsProtoGTE(r, ID776) {
		r.Varuint32(&x.GroupIndex)
	}
}
