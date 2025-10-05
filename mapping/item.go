package mapping

import (
	"encoding/base64"
	"encoding/json"
	"sync"

	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

type Item interface {
	// ItemRuntimeIDToName converts an item runtime ID to a string ID.
	ItemRuntimeIDToName(int32) (string, bool)
	// ItemNameToRuntimeID converts a string ID to an item runtime ID.
	ItemNameToRuntimeID(string) (int32, bool)
	// ItemRuntimeIDToVersion converts an item runtime ID to its version.
	ItemRuntimeIDToVersion(int32) (uint8, bool)
	// ItemRuntimeIDToData converts an item runtime ID to its data.
	ItemRuntimeIDToData(int32) (map[string]any, bool)
	RegisterEntry(string) int32
	RegisterEntryRID(string, int32, uint8, map[string]any)
	Air() int32
	ItemVersion() uint16
	ItemEntries() []ItemEntry
}

type ItemEntry struct {
	Name           string
	RuntimeID      int16
	ComponentBased bool
	Version        uint8
	Data           map[string]any
}

type DefaultItemMapping struct {
	mu sync.Mutex
	// itemRuntimeIDsToNames holds a map to translate item runtime IDs to string IDs.
	itemRuntimeIDsToNames map[int32]string
	// itemNamesToRuntimeIDs holds a map to translate item string IDs to runtime IDs.
	itemNamesToRuntimeIDs map[string]int32
	// itemRuntimeIDToVersion holds a map to translate item runtime IDs to versions.
	itemRuntimeIDToVersion map[int32]uint8
	// itemRuntimeIDToData holds a map to translate item runtime IDs to data.
	itemRuntimeIDToData map[int32]map[string]any
	// itemEntries holds a map to translate item string IDs to runtime IDs.
	itemEntries []ItemEntry

	airRID      int32
	itemVersion uint16
}

func NewItemMapping(requiredItemList []byte, itemVersion uint16) *DefaultItemMapping {
	itemRuntimeIDsToNames := make(map[int32]string)
	itemNamesToRuntimeIDs := make(map[string]int32)
	itemRuntimeIDToVersion := make(map[int32]uint8)
	itemRuntimeIDToData := make(map[int32]map[string]any)
	itemEntries := make([]ItemEntry, 0, 1600)
	var airRID *int32

	var m map[string]struct {
		RuntimeID      int16  `json:"runtime_id"`
		ComponentBased bool   `json:"component_based"`
		Version        *uint8 `json:"version"`
		Data           string `json:"component_nbt,omitempty"`
	}
	if err := json.Unmarshal(requiredItemList, &m); err != nil {
		panic(err)
	}

	for name, data := range m {
		rid := int32(data.RuntimeID)
		if name == "minecraft:air" {
			airRID = &rid
		}

		entry := ItemEntry{
			Name:           name,
			RuntimeID:      data.RuntimeID,
			ComponentBased: data.ComponentBased,
		}
		itemNamesToRuntimeIDs[name] = rid
		itemRuntimeIDsToNames[rid] = name
		if data.Version != nil {
			itemRuntimeIDToVersion[rid] = *data.Version
			entry.Version = *data.Version
		}
		if data.Data != "" {
			var nbtData map[string]any
			nbtBytes, err := base64.StdEncoding.DecodeString(data.Data)
			if err != nil {
				panic(err)
			}
			if err := nbt.Unmarshal(nbtBytes, &nbtData); err != nil {
				panic(err)
			}
			itemRuntimeIDToData[rid] = nbtData
			entry.Data = nbtData
		}
		itemEntries = append(itemEntries, entry)
	}

	if airRID == nil {
		panic("couldn't find air")
	}

	return &DefaultItemMapping{itemRuntimeIDsToNames: itemRuntimeIDsToNames, itemNamesToRuntimeIDs: itemNamesToRuntimeIDs, itemRuntimeIDToVersion: itemRuntimeIDToVersion, airRID: *airRID, itemVersion: itemVersion, itemRuntimeIDToData: itemRuntimeIDToData, itemEntries: itemEntries}
}

func (m *DefaultItemMapping) ItemRuntimeIDToName(runtimeID int32) (name string, found bool) {
	defer m.mu.Unlock()
	m.mu.Lock()
	name, ok := m.itemRuntimeIDsToNames[runtimeID]
	return name, ok
}

func (m *DefaultItemMapping) ItemNameToRuntimeID(name string) (runtimeID int32, found bool) {
	defer m.mu.Unlock()
	m.mu.Lock()
	rid, ok := m.itemNamesToRuntimeIDs[name]
	return rid, ok
}

func (m *DefaultItemMapping) ItemRuntimeIDToVersion(runtimeID int32) (version uint8, found bool) {
	defer m.mu.Unlock()
	m.mu.Lock()
	v, ok := m.itemRuntimeIDToVersion[runtimeID]
	return v, ok
}

func (m *DefaultItemMapping) ItemRuntimeIDToData(runtimeID int32) (data map[string]any, found bool) {
	defer m.mu.Unlock()
	m.mu.Lock()
	d, ok := m.itemRuntimeIDToData[runtimeID]
	return d, ok
}

func (m *DefaultItemMapping) RegisterEntry(name string) int32 {
	defer m.mu.Unlock()
	m.mu.Lock()
	if rid, ok := m.itemNamesToRuntimeIDs[name]; ok {
		return rid
	}
	nextRID := int32(len(m.itemRuntimeIDsToNames))
	m.itemNamesToRuntimeIDs[name] = nextRID
	m.itemRuntimeIDsToNames[nextRID] = name
	return nextRID
}

func (m *DefaultItemMapping) RegisterEntryRID(name string, rid int32, version uint8, data map[string]any) {
	defer m.mu.Unlock()
	m.mu.Lock()
	if _, ok := m.itemNamesToRuntimeIDs[name]; !ok {
		m.itemNamesToRuntimeIDs[name] = rid
		m.itemRuntimeIDsToNames[rid] = name
		m.itemRuntimeIDToVersion[rid] = version
		if data != nil {
			m.itemRuntimeIDToData[rid] = data
		}
	}
}

func (m *DefaultItemMapping) Air() int32 {
	defer m.mu.Unlock()
	m.mu.Lock()
	return m.airRID
}

func (m *DefaultItemMapping) ItemVersion() uint16 {
	defer m.mu.Unlock()
	m.mu.Lock()
	return m.itemVersion
}

func (m *DefaultItemMapping) ItemEntries() []ItemEntry {
	defer m.mu.Unlock()
	m.mu.Lock()
	entries := make([]ItemEntry, len(m.itemEntries))
	copy(entries, m.itemEntries)
	return entries
}
