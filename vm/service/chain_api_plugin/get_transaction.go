package chain_api_plugin

import (
	"context"
	"encoding/json"

	"github.com/MetalBlockchain/antelopevm/abi"
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
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

func init() {
	service.RegisterHandler("/v1/chain/get_transaction", GetTransaction)
}

func GetTransaction(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body GetTransactionRequest
		json.NewDecoder(c.Request.Body).Decode(&body)

		hash := core.TransactionIdType(*crypto.NewSha256String(body.TransactionId))
		session := vm.GetState().CreateSession(false)
		defer session.Discard()
		trx, err := session.FindTransactionByHash(hash)

		if err != nil {
			c.JSON(400, service.NewError(400, "account not found"))
			return
		}

		lastAcceptedId, _ := vm.LastAccepted(context.Background())
		lastAccepted, _ := session.FindBlockByHash(core.BlockHash(lastAcceptedId))
		signedTrx, _ := trx.Receipt.Transaction.GetSignedTransaction()
		response := &GetTransactionResponse{
			BlockNum:              trx.BlockNum,
			BlockTime:             trx.BlockTime.String(),
			HeadBlockNum:          uint32(lastAccepted.Header.Index),
			Id:                    trx.Hash,
			Irreversible:          true,
			LastIrreversibleBlock: uint32(lastAccepted.Header.Index),
			TransactionNum:        0,
			Traces:                trx.ActionTraces,
			MetaData: TransactionMetaData{
				Receipt:     TransactionReceipt(trx.Receipt),
				Transaction: *signedTrx,
			},
		}

		for index, trace := range response.Traces {
			if acc, err := session.FindAccountByName(trace.Action.Account); err == nil {
				if len(acc.Abi) > 0 {
					if abi, err := abi.NewABI(acc.Abi); err == nil {
						if data, err := abi.DecodeAction(trace.Action.Name, trace.Action.Data); err == nil {
							parsedData := map[string]interface{}{}

							if err := json.Unmarshal(data, &parsedData); err == nil {
								response.Traces[index].Action.ParsedData = parsedData
							} else {
								log15.Error("err 2", "e", err)
							}
						} else {
							log15.Error("err", "e", err)
						}
					}
				}
			}
		}

		c.JSON(200, response)
	}
}
