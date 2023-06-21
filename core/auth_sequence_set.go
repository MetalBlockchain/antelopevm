package core

import (
	"encoding/json"

	"github.com/MetalBlockchain/antelopevm/core/name"
)

//go:generate msgp
type AuthSequence struct {
	Account  name.AccountName
	Sequence uint64
}

type AuthSequenceSet struct {
	Data []AuthSequence
}

func NewAuthSequenceSet() AuthSequenceSet {
	return AuthSequenceSet{
		Data: make([]AuthSequence, 0),
	}
}

func (a *AuthSequenceSet) Set(account name.AccountName, sequence uint64) {
	for _, k := range a.Data {
		if k.Account == account {
			k.Sequence = sequence
			return
		}
	}

	a.Data = append(a.Data, AuthSequence{Account: account, Sequence: sequence})
}

func (a AuthSequenceSet) MarshalJSON() ([]byte, error) {
	var out [][]interface{}

	for _, v := range a.Data {
		out = append(out, []interface{}{v.Account, v.Sequence})
	}

	return json.Marshal(out)
}
