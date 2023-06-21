package service

import "github.com/MetalBlockchain/antelopevm/abi"

type GetAbiRequest struct {
	AccountName string `json:"account_name"`
}

type GetAbiResponse struct {
	AccountName string          `json:"account_name"`
	Abi         abi.ContractAbi `json:"abi,omitempty"`
}
