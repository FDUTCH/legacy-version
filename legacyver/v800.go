package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion800 ...
	ItemVersion800 = 251
	// BlockVersion800 ...
	BlockVersion800 int32 = (1 << 24) | (21 << 16) | (80 << 8)
)

var (
	//go:embed data/required_item_list_800.json
	requiredItemList800 []byte
	//go:embed data/block_states_800.nbt
	blockStateData800 []byte

	itemMappingLatest  = mapping.NewItemMapping(requiredItemList800, ItemVersion800)
	blockMappingLatest = mapping.NewBlockMapping(blockStateData800)
)
