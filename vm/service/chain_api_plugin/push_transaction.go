package chain_api_plugin

import (
	"github.com/MetalBlockchain/antelopevm/core"
)

type PushTransactionResults struct {
	TransactionId string                `json:"transaction_id"`
	Processed     core.TransactionTrace `json:"processed"`
}
