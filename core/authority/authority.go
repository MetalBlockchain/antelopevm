package authority

import (
	"encoding/json"

	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/hashicorp/go-set"
)

//go:generate msgp
type PermissionLevelSet = *set.HashSet[PermissionLevel, string]

func NewPermissionLevelSet(capacity int) PermissionLevelSet {
	return set.NewHashSet[PermissionLevel, string](capacity)
}

type PermissionLevel struct {
	Actor      name.AccountName    `serialize:"true" json:"actor"`
	Permission name.PermissionName `serialize:"true" json:"permission"`
}

func (p PermissionLevel) Hash() string {
	return crypto.Hash256(p).String()
}

func ComparePermissionLevel(first interface{}, second interface{}) int {
	if first.(PermissionLevel).Actor > second.(PermissionLevel).Actor {
		return 1
	} else if first.(PermissionLevel).Actor < second.(PermissionLevel).Actor {
		return -1
	}
	if first.(PermissionLevel).Permission > second.(PermissionLevel).Permission {
		return 1
	} else if first.(PermissionLevel).Permission < second.(PermissionLevel).Permission {
		return -1
	} else {
		return 0
	}
}

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

func (a Authority) MarshalJSON() ([]byte, error) {
	type Alias Authority

	b := struct {
		Alias
	}{
		Alias: (Alias)(a),
	}

	if b.Keys == nil {
		b.Keys = make([]KeyWeight, 0)
	}

	if b.Accounts == nil {
		b.Accounts = make([]PermissionLevelWeight, 0)
	}

	if b.Waits == nil {
		b.Waits = make([]WaitWeight, 0)
	}

	return json.Marshal(b)
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

func (a *Authority) GetBillableSize() uint64 {
	var accountsSize uint64 = uint64(len(a.Accounts)) * config.GetBillableSize("permission_level_weight")
	var waitsSize uint64 = uint64(len(a.Waits)) * config.GetBillableSize("wait_weight")
	var keysSize uint64 = 0

	for _, key := range a.Keys {
		keysSize += config.GetBillableSize("key_weight")
		keySize, err := rlp.EncodeSize(key.Key)

		if err != nil {
			panic(err)
		}

		keysSize += uint64(keySize)
	}

	return accountsSize + waitsSize + keysSize
}
