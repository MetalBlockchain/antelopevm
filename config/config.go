package config

import (
	"github.com/MetalBlockchain/antelopevm/core/name"
)

type billableSize struct {
	overhead uint64
	value    uint64
}

var (
	SystemAccountName         name.AccountName = name.StringToName("eosio")
	NullAccountName           name.AccountName = name.StringToName("eosio.null")
	ProducersAccountName      name.AccountName = name.StringToName("eosio.prods")
	MaxInlineActionDepth      uint16           = 4
	MaxInlineActionSize       uint64           = 4096
	SetCodeRamBytesMultiplier uint32           = 10

	// Percentages
	Percent100 int = 10000
	Percent1   int = 100

	// Permissions
	MaxAuthDepth  uint16              = 6
	ActiveName    name.PermissionName = name.StringToName("active")
	OwnerName     name.PermissionName = name.StringToName("owner")
	EosioAnyName  name.PermissionName = name.StringToName("eosio.any")
	EosioCodeName name.PermissionName = name.StringToName("eosio.code")

	MajorityProducersPermissionName name.PermissionName = name.StringToName("prod.major")
	MinorityProducersPermissionName name.PermissionName = name.StringToName("prod.minor")

	BlockIntervalMs       int64  = 500
	BlockTimestampEpochMs int64  = 946684800000
	MaxBlockCpuUsage      uint32 = 200000

	MinTransactionCpuUsage uint32 = 100
	MaxTransactionCpuUsage uint32 = 150000

	BillableAlignment uint64                  = 16
	BillableSize      map[string]billableSize = map[string]billableSize{
		"key_value_object":        {overhead: 32 * 2, value: 32 + 8 + 4 + 32*2},
		"table_id_object":         {overhead: 32 * 2, value: 44 + 32*2},
		"shared_authority":        {value: (3 * uint64(FixedOverheadSharedVectorRamBytes)) + 4},
		"permission_object":       {overhead: 5 * 32, value: 64 + 64 + 5*32},
		"permission_level_weight": {value: 24},
		"key_weight":              {value: 8},
		"wait_weight":             {value: 16},
	}

	FixedOverheadSharedVectorRamBytes uint32 = 16       ///< overhead accounts for fixed portion of size of shared_vector field
	OverheadPerRowPerIndexRamBytes    uint32 = 32       ///< overhead accounts for basic tracking structures in a row per index
	OverheadPerAccountRamBytes        uint32 = 2 * 1024 ///< overhead accounts for basic account storage and pre-pays features like account recovery

	RateLimitingPrecision uint32 = 1000 * 1000

	MinNetUsageDeltaBetweenBaseAndMaxForTrx uint32 = 10 * 1024

	// Wasm parameters
	DefaultMaxWasmMutableGlobalBytes uint32 = 1024
	DefaultMaxWasmTableElements      uint32 = 1024
	DefaultMaxWasmSectionElements    uint32 = 8192
	DefaultMaxWasmLinearMemoryInit   uint32 = 64 * 1024
	DefaultMaxWasmFuncLocalBytes     uint32 = 8192
	DefaultMaxWasmNestedStructures   uint32 = 1024
	DefaultMaxWasmSymbolBytes        uint32 = 8192
	DefaultMaxWasmModuleBytes        uint32 = 20 * 1024 * 1024
	DefaultMaxWasmCodeBytes          uint32 = 20 * 1024 * 1024
	DefaultMaxWasmPages              uint32 = 528
	DefaultMaxWasmCallDepth          uint32 = 251

	// Producer parameters
	MaxProducers int = 125
)

func GetBillableSize(kind string) uint64 {
	return ((BillableSize[kind].value + BillableAlignment - 1) / BillableAlignment) * BillableAlignment
}
