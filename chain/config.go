package chain

import "github.com/MetalBlockchain/antelopevm/core"

type Config struct {
	SystemAccountName         core.AccountName
	MaxInlineActionDepth      uint16
	SetCodeRamBytesMultiplier uint32

	// Permissions
	ActiveName   core.PermissionName
	OwnerName    core.PermissionName
	EosioAnyName core.PermissionName
}

func GetDefaultConfig() *Config {
	return &Config{
		SystemAccountName:         core.StringToName("eosio"),
		MaxInlineActionDepth:      4,
		SetCodeRamBytesMultiplier: 10,

		ActiveName:   core.StringToName("active"),
		OwnerName:    core.StringToName("owner"),
		EosioAnyName: core.StringToName("eosio.any"),
	}
}
