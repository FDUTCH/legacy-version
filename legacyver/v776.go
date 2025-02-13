package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion776 ...
	ItemVersion776 = 241
	// BlockVersion776 ...
	BlockVersion776 int32 = (1 << 24) | (21 << 16) | (60 << 8)
)

var (
	//go:embed data/item_runtime_ids_776.nbt
	itemRuntimeIDData776 []byte
	//go:embed data/required_item_list_776.json
	requiredItemList776 []byte
	//go:embed data/block_states_776.nbt
	blockStateData776 []byte

	itemMappingLatest  = mapping.NewItemMapping(itemRuntimeIDData776, requiredItemList776, ItemVersion776, false)
	blockMappingLatest = mapping.NewBlockMapping(blockStateData776)
)
