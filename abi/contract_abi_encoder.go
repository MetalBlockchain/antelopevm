package abi

import (
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
)

func (a *ContractAbi) Encode() ([]byte, error) {
	bytes, err := rlp.EncodeToBytes(a)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
