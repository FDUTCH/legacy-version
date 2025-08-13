package legacyver

import (
	_ "embed"

	"github.com/akmalfairuz/legacy-version/internal/chunk"
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion662 ...
	ItemVersion662 = 171
	// BlockVersion662 ...
	BlockVersion662 int32 = (1 << 24) | (20 << 16) | (70 << 8)
)

var (
	//go:embed data/required_item_list_662.json
	requiredItemList662 []byte
	//go:embed data/block_states_662.nbt
	blockStateData662 []byte
)

// New662 ...
func New662(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList662, ItemVersion662)
	blockMapping := mapping.NewBlockMapping(blockStateData662)

	return &Protocol{
		ver:             "1.20.70",
		id:              proto.ID662,
		blockTranslator: NewBlockTranslator(blockMapping, blockMappingLatest, chunk.NewNetworkPersistentEncoding(blockMapping, BlockVersion662), chunk.NewBlockPaletteEncoding(blockMapping, BlockVersion662), false),
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockMapping, blockMappingLatest),
	}
}
