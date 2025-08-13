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
)

func New800(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList800, ItemVersion800)
	blockMapping := mapping.NewBlockMapping(blockStateData800)
	return &Protocol{
		ver:             "1.21.80",
		id:              proto.ID800,
		blockTranslator: NewBlockTranslator(blockMapping, blockMappingLatest, chunk.NewNetworkPersistentEncoding(blockMapping, BlockVersion800), chunk.NewBlockPaletteEncoding(blockMapping, BlockVersion800), false),
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockMapping, blockMappingLatest),
	}
}
