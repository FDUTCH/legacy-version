package proto

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// SubChunkEntry contains the data of a sub-chunk entry relative to a center sub chunk position, used for the sub-chunk
// requesting system introduced in v1.18.10.
type SubChunkEntry struct {
	// Offset contains the offset between the sub-chunk position and the center position.
	Offset protocol.SubChunkOffset
	// Result is always one of the constants defined in the SubChunkResult constants.
	Result byte
	// RawPayload contains the serialized sub-chunk data.
	RawPayload []byte
	// HeightMapType is always one of the constants defined in the HeightMapData constants.
	HeightMapType byte
	// HeightMapData is the data for the height map.
	HeightMapData []int8
	// RenderHeightMapType is always one of the constants defined in the HeightMapData constants.
	RenderHeightMapType byte
	// RenderHeightMapData is the data for the render height map.
	RenderHeightMapData []int8
	// BlobHash is the hash of the blob.
	BlobHash uint64
}

// ToLatest ...
func (x *SubChunkEntry) ToLatest() protocol.SubChunkEntry {
	return protocol.SubChunkEntry{
		Offset:              x.Offset,
		Result:              x.Result,
		RawPayload:          x.RawPayload,
		HeightMapType:       x.HeightMapType,
		HeightMapData:       x.HeightMapData,
		RenderHeightMapType: x.RenderHeightMapType,
		RenderHeightMapData: x.RenderHeightMapData,
		BlobHash:            x.BlobHash,
	}
}

// FromLatest converts a protocol.SubChunkEntry to a SubChunkEntry. It is used to convert the latest
func (x *SubChunkEntry) FromLatest(y protocol.SubChunkEntry) SubChunkEntry {
	x.Offset = y.Offset
	x.Result = y.Result
	x.RawPayload = y.RawPayload
	x.HeightMapType = y.HeightMapType
	x.HeightMapData = y.HeightMapData
	x.RenderHeightMapType = y.RenderHeightMapType
	x.RenderHeightMapData = y.RenderHeightMapData
	x.BlobHash = y.BlobHash
	return *x
}

// Marshal encodes/decodes a SubChunkEntry assuming the blob cache is enabled.
func (x *SubChunkEntry) Marshal(r protocol.IO) {
	protocol.Single(r, &x.Offset)
	r.Uint8(&x.Result)
	if x.Result != protocol.SubChunkResultSuccessAllAir {
		r.ByteSlice(&x.RawPayload)
	}
	r.Uint8(&x.HeightMapType)
	if x.HeightMapType == protocol.HeightMapDataHasData {
		protocol.FuncSliceOfLen(r, 256, &x.HeightMapData, r.Int8)
	}
	if IsProtoGTE(r, ID818) {
		r.Uint8(&x.RenderHeightMapType)
		if x.RenderHeightMapType == protocol.HeightMapDataHasData {
			protocol.FuncSliceOfLen(r, 256, &x.RenderHeightMapData, r.Int8)
		}
	}
	r.Uint64(&x.BlobHash)
}

// SubChunkEntryNoCache encodes/decodes a SubChunkEntry assuming the blob cache is not enabled.
func SubChunkEntryNoCache(r protocol.IO, x *SubChunkEntry) {
	protocol.Single(r, &x.Offset)
	r.Uint8(&x.Result)
	r.ByteSlice(&x.RawPayload)
	r.Uint8(&x.HeightMapType)
	if x.HeightMapType == protocol.HeightMapDataHasData {
		protocol.FuncSliceOfLen(r, 256, &x.HeightMapData, r.Int8)
	}
	if IsProtoGTE(r, ID818) {
		r.Uint8(&x.RenderHeightMapType)
		if x.RenderHeightMapType == protocol.HeightMapDataHasData {
			protocol.FuncSliceOfLen(r, 256, &x.RenderHeightMapData, r.Int8)
		}
	}
}
