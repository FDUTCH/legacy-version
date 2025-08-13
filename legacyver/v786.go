package legacyver

import (
	_ "embed"

	"github.com/akmalfairuz/legacy-version/internal/chunk"
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion786 ...
	ItemVersion786 = 241
	// BlockVersion786 ...
	BlockVersion786 int32 = (1 << 24) | (21 << 16) | (70 << 8)
)

var (
	//go:embed data/required_item_list_786.json
	requiredItemList786 []byte
	//go:embed data/block_states_786.nbt
	blockStateData786 []byte
)

func New786(dragonflyMapping bool) *Protocol {
	itemMapping786 := mapping.NewItemMapping(requiredItemList786, ItemVersion786)
	blockMapping786 := mapping.NewBlockMapping(blockStateData786)
	return &Protocol{
		ver:             "1.21.70",
		id:              proto.ID786,
		blockTranslator: NewBlockTranslator(blockMapping786, blockMappingLatest, chunk.NewNetworkPersistentEncoding(blockMapping786, BlockVersion786), chunk.NewBlockPaletteEncoding(blockMapping786, BlockVersion786), false),
		itemTranslator:  NewItemTranslator(itemMapping786, itemMappingLatest(dragonflyMapping), blockMapping786, blockMappingLatest),
	}
}
