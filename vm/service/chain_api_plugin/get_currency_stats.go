package chain_api_plugin

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/MetalBlockchain/antelopevm/chain/asset"
	"github.com/MetalBlockchain/antelopevm/chain/name"
	"github.com/MetalBlockchain/antelopevm/crypto/rlp"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
)

type GetCurrencyStatsRequest struct {
	Code   name.AccountName `json:"code"`
	Symbol string           `json:"symbol"`
}

type CurrencyStats struct {
	Supply    asset.Asset      `json:"supply"`
	MaxSupply asset.Asset      `json:"max_supply"`
	Issuer    name.AccountName `json:"issuer"`
}

type GetCurrencyStatsResponse map[string]CurrencyStats

func init() {
	service.RegisterHandler("/v1/chain/get_currency_stats", service.Handler{
		Methods:     []string{http.MethodPost},
		HandlerFunc: GetCurrencyStats,
	})
}

func GetCurrencyStats(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body GetCurrencyStatsRequest
		json.NewDecoder(c.Request.Body).Decode(&body)
		response := GetCurrencyStatsResponse{}
		session := vm.GetState().CreateSession(false)
		defer session.Discard()
		symbol, err := asset.StringToSymbol(0, strings.ToUpper(body.Symbol))

		if err != nil {
			c.JSON(400, service.NewError(400, "invalid symbol"))
			return
		}

		scope := symbol >> 8
		table, err := session.FindTableByCodeScopeTable(body.Code, name.Name(scope), name.StringToName("stat"))

		if err != nil {
			c.JSON(400, service.NewError(400, "could not find currency stats"))
			return
		}

		iterator := session.FindKeyValuesByScope(table.ID)
		defer iterator.Close()

		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
			obj, err := iterator.Item()

			if err != nil {
				continue
			}

			ds := rlp.NewDecoder(obj.Value)
			result := CurrencyStats{}

			if err := ds.Decode(&result.Supply); err != nil {
				continue
			}

			if err := ds.Decode(&result.MaxSupply); err != nil {
				continue
			}

			if err := ds.Decode(&result.Issuer); err != nil {
				continue
			}

			response[result.Supply.Symbol.Name()] = result
		}

		c.JSON(200, response)
	}
}
