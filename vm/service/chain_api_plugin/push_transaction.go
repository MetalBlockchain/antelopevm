package chain_api_plugin

import (
	"encoding/json"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
)

type PushTransactionResults struct {
	TransactionId string                `json:"transaction_id"`
	Processed     core.TransactionTrace `json:"processed"`
}

func init() {
	service.RegisterHandler("/v1/chain/send_transaction", PushTransaction)
	service.RegisterHandler("/v1/chain/push_transaction", PushTransaction)
}

func PushTransaction(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		var trx core.PackedTransaction

		if err := json.NewDecoder(c.Request.Body).Decode(&trx); err != nil {
			c.JSON(400, "failed to decode body")
			return
		}

		if ok := vm.GetMempool().Add(&trx); !ok {
			c.JSON(400, "could not submit trx")
			return
		}

		c.JSON(202, PushTransactionResults{
			TransactionId: trx.Id.String(),
			Processed: core.TransactionTrace{
				Hash: trx.Id,
				Receipt: core.TransactionReceipt{
					TransactionReceiptHeader: core.TransactionReceiptHeader{
						Status: core.TransactionStatusExecuted,
					},
				},
				ActionTraces: make([]core.ActionTrace, 0),
			},
		})
	}
}
