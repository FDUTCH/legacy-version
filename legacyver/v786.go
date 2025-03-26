package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion786 ...
	ItemVersion786 = 241
	// BlockVersion786 ...
	BlockVersion786 int32 = (1 << 24) | (21 << 16) | (70 << 8)
)

var (
	//go:embed data/item_runtime_ids_786.nbt
	itemRuntimeIDData786 []byte
	//go:embed data/required_item_list_786.json
	requiredItemList786 []byte
	//go:embed data/block_states_786.nbt
	blockStateData786 []byte

	itemMappingLatest  = mapping.NewItemMapping(itemRuntimeIDData786, requiredItemList786, ItemVersion786, false)
	blockMappingLatest = mapping.NewBlockMapping(blockStateData786)
)
