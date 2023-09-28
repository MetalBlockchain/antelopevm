package chain_api_plugin

import (
	"encoding/json"
	"net/http"

	"github.com/MetalBlockchain/antelopevm/chain/transaction"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
)

type PushTransactionResults struct {
	TransactionId string                       `json:"transaction_id"`
	Processed     transaction.TransactionTrace `json:"processed"`
}

func init() {
	service.RegisterHandler("/v1/chain/send_transaction", service.Handler{
		Methods:     []string{http.MethodPost},
		HandlerFunc: PushTransaction,
	})
	service.RegisterHandler("/v1/chain/push_transaction", service.Handler{
		Methods:     []string{http.MethodPost},
		HandlerFunc: PushTransaction,
	})
}

func PushTransaction(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		var trx transaction.PackedTransaction

		if err := json.NewDecoder(c.Request.Body).Decode(&trx); err != nil {
			c.JSON(400, "failed to decode body")
			return
		}

		if ok := vm.GetMempool().Add(&trx); !ok {
			c.JSON(400, "could not submit trx")
			return
		}

		id, _ := trx.ID()

		c.JSON(202, PushTransactionResults{
			TransactionId: id.String(),
			Processed: transaction.TransactionTrace{
				Hash: *id,
				Receipt: transaction.TransactionReceipt{
					TransactionReceiptHeader: transaction.TransactionReceiptHeader{
						Status: transaction.TransactionStatusExecuted,
					},
				},
				ActionTraces: make([]transaction.ActionTrace, 0),
			},
		})
	}
}
