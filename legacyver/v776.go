package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/internal/chunk"
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion776 ...
	ItemVersion776 = 241
	// BlockVersion776 ...
	BlockVersion776 int32 = (1 << 24) | (21 << 16) | (60 << 8)
)

var (
	//go:embed data/required_item_list_776.json
	requiredItemList776 []byte
	//go:embed data/block_states_776.nbt
	blockStateData776 []byte

	itemMapping776  = mapping.NewItemMapping(requiredItemList776, ItemVersion776)
	blockMapping776 = mapping.NewBlockMapping(blockStateData776)
)

func New776() *Protocol {
	return &Protocol{
		ver:             "1.21.60",
		id:              proto.ID776,
		blockTranslator: NewBlockTranslator(blockMapping776, blockMappingLatest, chunk.NewNetworkPersistentEncoding(blockMapping776, BlockVersion776), chunk.NewBlockPaletteEncoding(blockMapping776, BlockVersion776), false),
		itemTranslator:  NewItemTranslator(itemMapping776, itemMappingLatest, blockMapping776, blockMappingLatest),
	}
}
