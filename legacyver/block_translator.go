package legacyver

import (
	"bytes"
	_ "unsafe"

	"github.com/akmalfairuz/legacy-version/mapping"
	"github.com/cespare/xxhash/v2"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type BlockTranslator interface {
	// DowngradeBlockPackets downgrades the input block packets to legacy block packets.
	DowngradeBlockPackets([]packet.Packet, *minecraft.Conn) (result []packet.Packet)
	// UpgradeBlockPackets upgrades the input block packets to the latest block packets.
	UpgradeBlockPackets([]packet.Packet, *minecraft.Conn) (result []packet.Packet)
	// DowngradeLevelChunk downgrades the given LevelChunk packet to a legacy format.
	DowngradeLevelChunk(*packet.LevelChunk) error
	// BlockMapping returns the block mapping used by this translator.
	BlockMapping() mapping.Block
}
type DefaultBlockTranslator struct {
	mapping              mapping.Block
	latest               mapping.Block
	dimensionDefinitions []protocol.DimensionDefinition
}

func NewBlockTranslator(mapping mapping.Block, latestMapping mapping.Block) *DefaultBlockTranslator {
	return &DefaultBlockTranslator{mapping: mapping, latest: latestMapping}
}

func (t *DefaultBlockTranslator) BlockMapping() mapping.Block {
	return t.mapping
}

func (t *DefaultBlockTranslator) DowngradeBlockPackets(pks []packet.Packet, conn *minecraft.Conn) (result []packet.Packet) {
	for _, pk := range pks {
		switch pk := pk.(type) {
		case *packet.LevelChunk:
			if !EnableChunkTranslation {
				break
			}
			if err := t.DowngradeLevelChunk(pk); err != nil {
				//fmt.Println(err)
				break
			}
		case *packet.SubChunk:
			if !EnableChunkTranslation {
				break
			}
			var (
				r   = dimensionRange(pk)
				c   *chunk.Chunk
				err error
			)

			for i := range len(pk.SubChunkEntries) {
				c, err = t.DowngradeSubChunkEntry(&pk.SubChunkEntries[i], c, r)
				if err != nil {
					panic(err)
				}
			}
		case *packet.ClientCacheMissResponse:
			c := chunk.New(t.latest.Air(), cube.Range{-64, 319})
			for i, blob := range pk.Blobs {
				buf := bytes.NewBuffer(blob.Payload)
				ind := byte(0)
				subChunk, err := decodeSubChunk(buf, c, &ind, chunk.NetworkEncoding)
				if err != nil {
					// Has a possibility to be a biome, ignore then
					continue
				}
				t.DowngradeSubChunk(subChunk)
				blob.Payload = append(chunk.EncodeSubChunk(c, chunk.NetworkEncoding, int(ind)), buf.Bytes()...)
				blob.Hash = xxhash.Sum64(blob.Payload)
				pk.Blobs[i] = blob
			}
		case *packet.UpdateSubChunkBlocks:
			for i, block := range pk.Blocks {
				block.BlockRuntimeID = t.DowngradeBlockRuntimeID(block.BlockRuntimeID)
				pk.Blocks[i] = block
			}
			for i, block := range pk.Extra {
				block.BlockRuntimeID = t.DowngradeBlockRuntimeID(block.BlockRuntimeID)
				pk.Extra[i] = block
			}
		case *packet.UpdateBlock:
			pk.NewBlockRuntimeID = t.DowngradeBlockRuntimeID(pk.NewBlockRuntimeID)
		case *packet.UpdateBlockSynced:
			pk.NewBlockRuntimeID = t.DowngradeBlockRuntimeID(pk.NewBlockRuntimeID)
		case *packet.InventoryTransaction:
			if transactionData, ok := pk.TransactionData.(*protocol.UseItemTransactionData); ok {
				transactionData.BlockRuntimeID = t.DowngradeBlockRuntimeID(transactionData.BlockRuntimeID)
				pk.TransactionData = transactionData
			}
		case *packet.LevelEvent:
			switch pk.EventType {
			case packet.LevelEventParticleLegacyEvent | 20: // terrain
				fallthrough
			case packet.LevelEventParticlesDestroyBlock:
				fallthrough
			case packet.LevelEventParticlesDestroyBlockNoSound:
				pk.EventData = int32(t.DowngradeBlockRuntimeID(uint32(pk.EventData)))
			case packet.LevelEventParticlesCrackBlock:
				face := pk.EventData >> 24
				rid := t.DowngradeBlockRuntimeID(uint32(pk.EventData & 0xffff))
				pk.EventData = int32(rid) | (face << 24)
			}
		case *packet.LevelSoundEvent:
			switch pk.SoundType {
			case packet.SoundEventBreak:
				fallthrough
			case packet.SoundEventPlace:
				fallthrough
			case packet.SoundEventHit:
				fallthrough
			case packet.SoundEventLand:
				fallthrough
			case packet.SoundEventItemUseOn:
				pk.ExtraData = int32(t.DowngradeBlockRuntimeID(uint32(pk.ExtraData)))
			}
		case *packet.AddActor:
			if pk.EntityType == "minecraft:falling_block" {
				pk.EntityMetadata = t.downgradeEntityMetadata(pk.EntityMetadata)
			}
		case *packet.SetActorData:
			//pk.EntityMetadata = t.downgradeEntityMetadata(pk.EntityMetadata)
		case *packet.StartGame:
			t.latest.Adjust(pk.Blocks)
			t.mapping.Adjust(pk.Blocks)
		case *packet.ResourcePackStack:
			var packs []protocol.StackResourcePack
			for _, pack := range pk.TexturePacks {
				if pack.UUID == "0fba4063-dba1-4281-9b89-ff9390653530" {
					continue
				}
				packs = append(packs, pack)
			}
			pk.TexturePacks = packs
		}
		result = append(result, pk)
	}
	return result
}

func (t *DefaultBlockTranslator) DowngradeLevelChunk(pk *packet.LevelChunk) error {
	count := int(pk.SubChunkCount)
	if count == protocol.SubChunkRequestModeLimitless || count == protocol.SubChunkRequestModeLimited {
		return nil
	}
	buf := bytes.NewBuffer(pk.RawPayload)
	c, err := chunk.NetworkDecodeBuffer(t.latest.Air(), buf, count, dimensionRange(pk))
	if err != nil {
		return err
	}

	panic("not supported yet")
	sub := c.Sub()
	for i := range sub {
		t.DowngradeSubChunk(sub[i])
	}
	var (
		data   = chunk.Encode(c, chunk.NetworkEncoding)
		blobs  = append(data.SubChunks, data.Biomes)
		hashes = make([]uint64, len(blobs))
		m      = make(map[uint64]struct{}, len(blobs))
	)
	for i, blob := range blobs {
		h := xxhash.Sum64(blob)
		hashes[i], m[h] = h, struct{}{}
	}
	return nil
}

func (t *DefaultBlockTranslator) DowngradeSubChunkEntry(entry *protocol.SubChunkEntry, c *chunk.Chunk, r cube.Range) (*chunk.Chunk, error) {
	if entry.Result != protocol.SubChunkResultSuccess || len(entry.RawPayload) == 0 {
		return c, nil
	}

	if entry.RawPayload[0] == 10 {
		entry.RawPayload[0] = 9
	}

	if c == nil {
		c = chunk.New(t.latest.Air(), r)
	}

	buf := bytes.NewBuffer(entry.RawPayload)
	var index byte
	sub, err := decodeSubChunk(buf, c, &index, chunk.NetworkEncoding)
	if err != nil {
		return c, err
	}
	t.DowngradeSubChunk(sub)
	serializedSubChunk := chunk.EncodeSubChunk(c, chunk.NetworkEncoding, int(index))
	entry.BlobHash = xxhash.Sum64(serializedSubChunk)
	entry.RawPayload = append(serializedSubChunk, buf.Bytes()...)
	return c, err
}

func (t *DefaultBlockTranslator) UpgradeBlockPackets(pks []packet.Packet, conn *minecraft.Conn) (result []packet.Packet) {
	for _, pk := range pks {
		switch pk := pk.(type) {
		case *packet.InventoryTransaction:
			if transactionData, ok := pk.TransactionData.(*protocol.UseItemTransactionData); ok {
				transactionData.BlockRuntimeID = t.UpgradeBlockRuntimeID(transactionData.BlockRuntimeID)
				pk.TransactionData = transactionData
			}
		case *packet.SetActorData:
			//pk.EntityMetadata = t.upgradeEntityMetadata(pk.EntityMetadata)
		}
		result = append(result, pk)
	}
	return result
}

func (t *DefaultBlockTranslator) DowngradeBlockRuntimeID(input uint32) uint32 {
	if t.latest == t.mapping {
		return input
	}
	state, ok := t.latest.RuntimeIDToState(input)
	if !ok {
		return t.mapping.InfoUpdate()
	}
	runtimeID, ok := t.mapping.StateToRuntimeID(state)
	if !ok {
		return t.mapping.InfoUpdate()
	}
	return runtimeID
}

func (t *DefaultBlockTranslator) DowngradeSubChunk(input *chunk.SubChunk) {
	if t.latest == t.mapping {
		return
	}
	for _, storage := range input.Layers() {
		storage.Palette().Replace(t.DowngradeBlockRuntimeID)
	}
}

func (t *DefaultBlockTranslator) downgradeEntityMetadata(metadata map[uint32]any) map[uint32]any {
	if t.latest == t.mapping {
		return metadata
	}
	if latestRID, ok := metadata[protocol.EntityDataKeyVariant]; ok {
		metadata[protocol.EntityDataKeyVariant] = int32(t.DowngradeBlockRuntimeID(uint32(latestRID.(int32))))
	}
	return metadata
}

func (t *DefaultBlockTranslator) UpgradeBlockRuntimeID(input uint32) uint32 {
	if t.latest == t.mapping {
		return input
	}
	state, ok := t.mapping.RuntimeIDToState(input)
	if !ok {
		return t.latest.InfoUpdate()
	}
	runtimeID, ok := t.latest.StateToRuntimeID(state)
	if !ok {
		return t.latest.InfoUpdate()
	}
	return runtimeID
}

func (t *DefaultBlockTranslator) upgradeEntityMetadata(metadata map[uint32]any) map[uint32]any {
	if t.latest == t.mapping {
		return metadata
	}
	if latestRID, ok := metadata[protocol.EntityDataKeyVariant]; ok {
		metadata[protocol.EntityDataKeyVariant] = int32(t.UpgradeBlockRuntimeID(uint32(latestRID.(int32))))
	}
	return metadata
}

func dimensionRange(p packet.Packet) cube.Range {
	var id int
	switch pk := p.(type) {
	case *packet.LevelChunk:
		id = int(pk.Dimension)
	case *packet.SubChunk:
		id = int(pk.Dimension)
	}
	dim, _ := world.DimensionByID(id)
	return dim.Range()
}

//go:linkname decodeSubChunk github.com/df-mc/dragonfly/server/world/chunk.decodeSubChunk
func decodeSubChunk(buf *bytes.Buffer, c *chunk.Chunk, index *byte, e chunk.Encoding) (*chunk.SubChunk, error)
