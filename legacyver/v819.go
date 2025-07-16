package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion819 ...
	ItemVersion819 = 271
	// BlockVersion819 ...
	BlockVersion819 int32 = (1 << 24) | (21 << 16) | (90 << 8)
)

var (
	//go:embed data/dragonfly_items.json
	dragonflyLatestItemList []byte
	//go:embed data/required_item_list_819.json
	requiredItemList819 []byte
	//go:embed data/block_states_819.nbt
	blockStateData819 []byte

	itemMappingLatestPocketMine = mapping.NewItemMapping(requiredItemList819, ItemVersion819)
	itemMappingLatestDragonfly  = mapping.NewItemMapping(dragonflyLatestItemList, ItemVersion819)
	blockMappingLatest          = mapping.NewBlockMapping(blockStateData819)
)

func itemMappingLatest(dragonflyMapping bool) mapping.Item {
	if dragonflyMapping {
		return itemMappingLatestDragonfly
	}
	return itemMappingLatestPocketMine
}
