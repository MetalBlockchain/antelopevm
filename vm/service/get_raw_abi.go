package service

import (
	"github.com/MetalBlockchain/antelopevm/crypto"
)

type GetRawAbiRequest struct {
	AccountName string `json:"account_name"`
}

type GetRawAbiResponse struct {
	AccountName string        `json:"account_name"`
	CodeHash    crypto.Sha256 `json:"code_hash"`
	AbiHash     crypto.Sha256 `json:"abi_hash"`
	Abi         string        `json:"abi"`
}
