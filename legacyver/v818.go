package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/mapping"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

const (
	// ItemVersion818 ...
	ItemVersion818 = 261
	// BlockVersion818 ...
	BlockVersion818 int32 = (1 << 24) | (21 << 16) | (90 << 8)
)

var (
	//go:embed data/required_item_list_818.json
	requiredItemList818 []byte
	//go:embed data/block_states_818.nbt
	blockStateData818 []byte
	//go:embed data/dragonfly_latest_vanilla_items.nbt
	dragonflyLatestVanillaItemsRaw []byte

	itemMappingLatestPocketMine = mapping.NewItemMapping(requiredItemList818, ItemVersion818)
	itemMappingLatestDragonfly  = mapping.NewItemMapping(requiredItemList818, ItemVersion818)
	blockMappingLatest          = mapping.NewBlockMapping(blockStateData818)
)

func init() {
	var m map[string]struct {
		RuntimeID      int32          `nbt:"runtime_id"`
		ComponentBased bool           `nbt:"component_based"`
		Version        int32          `nbt:"version"`
		Data           map[string]any `nbt:"data,omitempty"`
	}
	if err := nbt.Unmarshal(dragonflyLatestVanillaItemsRaw, &m); err != nil {
		panic(err)
	}

	for identifier, v := range m {
		itemMappingLatestDragonfly.SetItemRuntimeID(identifier, v.RuntimeID)
	}
}

func itemMappingLatest(dragonflyMapping bool) mapping.Item {
	if dragonflyMapping {
		return itemMappingLatestDragonfly
	}
	return itemMappingLatestPocketMine
}
