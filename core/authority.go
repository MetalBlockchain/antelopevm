package core

import "github.com/MetalBlockchain/antelopevm/crypto/ecc"

type WeightType uint16

type PermissionLevelWeight struct {
	Permission PermissionLevel `serialize:"true" json:"permission"`
	Weight     WeightType      `serialize:"true" json:"weight"`
}

type KeyWeight struct {
	Key    ecc.PublicKey `serialize:"true" json:"key"`
	Weight WeightType    `serialize:"true" json:"weight"`
}

type WaitWeight struct {
	WaitSec uint32     `serialize:"true" json:"wait_sec"`
	Weight  WeightType `serialize:"true" json:"weight"`
}

type Authority struct {
	Threshold uint32                  `serialize:"true" json:"threshold"`
	Keys      []KeyWeight             `serialize:"true" json:"keys"`
	Accounts  []PermissionLevelWeight `serialize:"true" json:"accounts"`
	Waits     []WaitWeight            `serialize:"true" json:"waits"`
}

func (auth *Authority) IsValid() bool {
	var totalWeight uint32 = 0

	if len(auth.Accounts)+len(auth.Keys)+len(auth.Waits) > 1<<16 {
		return false
	}

	if auth.Threshold == 0 {
		return false
	}

	for i, k := range auth.Keys {
		if i > 0 && !(ecc.ComparePubKey(auth.Keys[i-1].Key, k.Key) == -1) {
			return false
		}

		totalWeight += uint32(k.Weight)
	}

	for i, a := range auth.Accounts {
		if i > 0 && !(ComparePermissionLevel(auth.Accounts[i-1].Permission, a.Permission) == -1) {
			return false
		}

		totalWeight += uint32(a.Weight)
	}

	for i, w := range auth.Waits {
		if i > 0 && auth.Waits[i-1].WaitSec >= w.WaitSec {
			return false
		}

		totalWeight += uint32(w.Weight)
	}

	return totalWeight >= auth.Threshold
}
