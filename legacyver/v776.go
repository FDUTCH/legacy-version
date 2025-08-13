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
)

func New776(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList776, ItemVersion776)
	blockMapping := mapping.NewBlockMapping(blockStateData776)

	return &Protocol{
		ver:             "1.21.60",
		id:              proto.ID776,
		blockTranslator: NewBlockTranslator(blockMapping, blockMappingLatest, chunk.NewNetworkPersistentEncoding(blockMapping, BlockVersion776), chunk.NewBlockPaletteEncoding(blockMapping, BlockVersion776), false),
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockMapping, blockMappingLatest),
	}
}
