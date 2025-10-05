package legacyver

import (
	_ "embed"

	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion844 ...
	ItemVersion844 = 241
	// BlockVersion844 ...
	BlockVersion844 int32 = (1 << 24) | (21 << 16) | (110 << 8)
)

var (
	//go:embed data/dragonfly_items.json
	dragonflyLatestItemList []byte
	//go:embed data/required_item_list_844.json
	requiredItemList844 []byte
	//go:embed data/block_states_844.nbt
	blockStateData844 []byte

	itemMappingLatestPocketMine = mapping.NewItemMapping(requiredItemList844, ItemVersion844)
	itemMappingLatestDragonfly  = mapping.NewItemMapping(dragonflyLatestItemList, ItemVersion844)
	blockMappingLatest          = mapping.NewBlockMapping(blockStateData844)
)

func itemMappingLatest(dragonflyMapping bool) mapping.Item {
	if dragonflyMapping {
		return itemMappingLatestDragonfly
	}
	return itemMappingLatestPocketMine
}
