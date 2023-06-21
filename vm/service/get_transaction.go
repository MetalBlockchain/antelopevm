package service

import (
	"encoding/json"

	"github.com/MetalBlockchain/antelopevm/core"
)

type GetTransactionRequest struct {
	TransactionId string `json:"id"`
}

type TransactionReceipt core.TransactionReceipt

func (t *TransactionReceipt) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Status        core.TransactionStatus `json:"status"`
		CpuUsageUs    uint32                 `json:"cpu_usage_us"`
		NetUsageWords core.Vuint32           `json:"net_usage_words"`
		Transactions  []interface{}          `json:"trx"`
	}{
		CpuUsageUs:    t.TransactionReceiptHeader.CpuUsageUs,
		NetUsageWords: t.TransactionReceiptHeader.NetUsageWords,
		Status:        t.TransactionReceiptHeader.Status,
		Transactions: []interface{}{
			1,
			t.Transaction,
		},
	})
}

type TransactionMetaData struct {
	Receipt     TransactionReceipt     `json:"receipt"`
	Transaction core.SignedTransaction `json:"trx"`
}

type GetTransactionResponse struct {
	BlockNum              uint64                 `json:"block_num"`
	BlockTime             string                 `json:"block_time"`
	HeadBlockNum          uint32                 `json:"head_block_num"`
	Id                    core.TransactionIdType `json:"id"`
	Irreversible          bool                   `json:"irreversible"`
	LastIrreversibleBlock uint32                 `json:"last_irreversible_block"`
	TransactionNum        uint32                 `json:"transaction_num"`
	Traces                []core.ActionTrace     `json:"traces"`
	MetaData              TransactionMetaData    `json:"trx"`
}
