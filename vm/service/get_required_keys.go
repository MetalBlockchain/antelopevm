package service

import (
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
)

type RequiredKeysRequest struct {
	Transaction   core.Transaction `json:"transaction"`
	AvailableKeys []ecc.PublicKey  `json:"available_keys"`
}

type RequiredKeysResponse struct {
	RequiredKeys []ecc.PublicKey `json:"required_keys"`
}
