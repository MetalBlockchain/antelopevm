package chain

import "github.com/MetalBlockchain/antelopevm/chain/types"

type Config struct {
	SystemAccountName         types.AccountName
	MaxInlineActionDepth      uint16
	SetCodeRamBytesMultiplier uint32

	// Permissions
	ActiveName   types.PermissionName
	OwnerName    types.PermissionName
	EosioAnyName types.PermissionName
}

func GetDefaultConfig() *Config {
	return &Config{
		SystemAccountName:         types.N("eosio"),
		MaxInlineActionDepth:      4,
		SetCodeRamBytesMultiplier: 10,

		ActiveName:   types.N("active"),
		OwnerName:    types.N("owner"),
		EosioAnyName: types.N("eosio.any"),
	}
}
