package service

import (
	"github.com/MetalBlockchain/antelopevm/crypto"
)

type GetCodeHashRequest struct {
	AccountName string `json:"account_name"`
}

type GetCodeHashResponse struct {
	AccountName string        `json:"account_name"`
	CodeHash    crypto.Sha256 `json:"code_hash"`
}

func NewGetCodeHashResponse(accountName string, codeHash crypto.Sha256) GetCodeHashResponse {
	return GetCodeHashResponse{
		AccountName: accountName,
		CodeHash:    codeHash,
	}
}
