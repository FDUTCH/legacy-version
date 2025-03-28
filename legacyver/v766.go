package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/internal/chunk"
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion766 ...
	ItemVersion766 = 221
	// BlockVersion766 ...
	BlockVersion766 int32 = (1 << 24) | (21 << 16) | (50 << 8)
)

var (
	//go:embed data/required_item_list_766.json
	requiredItemList766 []byte
	//go:embed data/block_states_766.nbt
	blockStateData766 []byte

	itemMapping766  = mapping.NewItemMapping(requiredItemList766, ItemVersion766)
	blockMapping766 = mapping.NewBlockMapping(blockStateData766)
)

func New766() *Protocol {
	return &Protocol{
		ver:             "1.21.50",
		id:              proto.ID766,
		blockTranslator: NewBlockTranslator(blockMapping766, blockMappingLatest, chunk.NewNetworkPersistentEncoding(blockMapping766, BlockVersion766), chunk.NewBlockPaletteEncoding(blockMapping766, BlockVersion766), false),
		itemTranslator:  NewItemTranslator(itemMapping766, itemMappingLatest, blockMapping766, blockMappingLatest),
	}
}
