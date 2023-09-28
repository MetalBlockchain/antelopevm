package chain_api_plugin

import (
	"encoding/json"
	"net/http"

	"github.com/MetalBlockchain/antelopevm/chain/asset"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
)

type GetCurrencyBalanceRequest struct {
	Code    name.AccountName `json:"code"`
	Account name.AccountName `json:"account"`
	Symbol  string           `json:"symbol"`
}

func init() {
	service.RegisterHandler("/v1/chain/get_currency_balance", service.Handler{
		Methods:     []string{http.MethodPost},
		HandlerFunc: GetCurrencyBalance,
	})
}

func GetCurrencyBalance(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body GetCurrencyBalanceRequest
		json.NewDecoder(c.Request.Body).Decode(&body)
		session := vm.GetState().CreateSession(false)
		defer session.Discard()
		results := make([]asset.Asset, 0)

		table, err := session.FindTableByCodeScopeTable(body.Code, body.Account, name.StringToName("accounts"))

		if err != nil {
			c.JSON(200, results)
			return
		}

		iterator := session.FindKeyValuesByScope(table.ID)
		defer iterator.Close()

		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
			item, err := iterator.Item()

			if err == nil {
				asset := asset.Asset{}

				if err := rlp.DecodeBytes(item.Value, &asset); err == nil {

					if len(body.Symbol) == 0 || body.Symbol == asset.Symbol.Symbol {
						results = append(results, asset)

						if len(body.Symbol) > 0 {
							break
						}
					}
				}
			}
		}

		c.JSON(200, results)
	}
}
