package legacyver

import (
	_ "embed"

	"github.com/akmalfairuz/legacy-version/internal/chunk"
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
	"github.com/akmalfairuz/legacy-version/mapping"
)

const (
	// ItemVersion748 ...
	ItemVersion748 = 221
	// BlockVersion748 ...
	BlockVersion748 int32 = (1 << 24) | (21 << 16) | (40 << 8)
)

var (
	//go:embed data/required_item_list_748.json
	requiredItemList748 []byte
	//go:embed data/block_states_748.nbt
	blockStateData748 []byte
)

func New748(dragonflyMapping bool) *Protocol {
	itemMapping := mapping.NewItemMapping(requiredItemList748, ItemVersion748)
	blockMapping := mapping.NewBlockMapping(blockStateData748)

	return &Protocol{
		ver:             "1.21.40",
		id:              proto.ID748,
		blockTranslator: NewBlockTranslator(blockMapping, blockMappingLatest, chunk.NewNetworkPersistentEncoding(blockMapping, BlockVersion748), chunk.NewBlockPaletteEncoding(blockMapping, BlockVersion748), false),
		itemTranslator:  NewItemTranslator(itemMapping, itemMappingLatest(dragonflyMapping), blockMapping, blockMappingLatest),
	}
}
