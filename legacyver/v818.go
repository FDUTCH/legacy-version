package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/mapping"
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

	itemMappingLatestWithDebugStick    = mapping.NewItemMapping(requiredItemList818, ItemVersion818, false)
	itemMappingLatestWithoutDebugStick = mapping.NewItemMapping(requiredItemList818, ItemVersion818, true)
	blockMappingLatest                 = mapping.NewBlockMapping(blockStateData818)
)

func itemMappingLatest(deleteDebugStick bool) mapping.Item {
	if deleteDebugStick {
		return itemMappingLatestWithoutDebugStick
	}
	return itemMappingLatestWithDebugStick
}
