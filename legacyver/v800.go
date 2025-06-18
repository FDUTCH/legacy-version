package legacyver

import (
	_ "embed"
	"github.com/akmalfairuz/legacy-version/internal/chunk"
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
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

	itemMapping800  = mapping.NewItemMapping(requiredItemList800, ItemVersion800)
	blockMapping800 = mapping.NewBlockMapping(blockStateData800)
)

func New800() *Protocol {
	return &Protocol{
		ver:             "1.21.80",
		id:              proto.ID800,
		blockTranslator: NewBlockTranslator(blockMapping800, blockMappingLatest, chunk.NewNetworkPersistentEncoding(blockMapping800, BlockVersion800), chunk.NewBlockPaletteEncoding(blockMapping800, BlockVersion800), false),
		itemTranslator:  NewItemTranslator(itemMapping800, itemMappingLatest, blockMapping800, blockMappingLatest),
	}
}
