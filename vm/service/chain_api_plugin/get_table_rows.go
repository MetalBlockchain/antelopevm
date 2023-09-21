package chain_api_plugin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MetalBlockchain/antelopevm/abi"
	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
)

type GetTableRowsRequest struct {
	Code  name.AccountName `json:"code"`
	Scope name.ScopeName   `json:"scope"`
	Table name.TableName   `json:"table"`
}

type TableRow struct {
	Data  map[string]interface{} `json:"data"`
	Payer string                 `json:"payer"`
}

type GetTableRowsResponse struct {
	More bool       `json:"more"`
	Rows []TableRow `json:"rows"`
}

func init() {
	service.RegisterHandler("/v1/chain/get_table_rows", service.Handler{
		Methods:     []string{http.MethodPost},
		HandlerFunc: GetTableRows,
	})
}

func GetTableRows(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body GetTableRowsRequest
		json.NewDecoder(c.Request.Body).Decode(&body)
		response := &GetTableRowsResponse{
			More: false,
			Rows: make([]TableRow, 0),
		}
		session := vm.GetState().CreateSession(false)
		defer session.Discard()
		acc, err := session.FindAccountByName(body.Code)

		if err != nil {
			c.AbortWithError(400, fmt.Errorf("account with name %s does not exist", body.Code))
			return
		}

		abi, err := abi.NewABI(acc.Abi)

		if err != nil {
			c.JSON(400, service.NewError(400, "failed to parse ABI"))
			return
		}

		tableDef := abi.TableForName(body.Table)

		if tableDef == nil {
			c.JSON(400, service.NewError(400, fmt.Sprintf("table %s is not specified in the ABI", body.Table)))
			return
		}

		table, err := session.FindTableByCodeScopeTable(body.Code, body.Scope, body.Table)

		if err != nil {
			c.JSON(400, service.NewError(400, "failed to find table"))
			return
		}

		keyValues := make([]*core.KeyValue, 0)
		keyValueIterator := session.FindKeyValuesByScope(table.ID)
		defer keyValueIterator.Close()

		for keyValueIterator.Rewind(); keyValueIterator.Valid(); keyValueIterator.Next() {
			if item, err := keyValueIterator.Item(); err == nil {
				keyValues = append(keyValues, item)
			}
		}

		for _, keyValue := range keyValues {
			structDef := abi.StructForName(tableDef.Type)

			if data, err := abi.DecodeStruct(structDef.Name, keyValue.Value); err == nil {
				parsedData := map[string]interface{}{}

				if err := json.Unmarshal(data, &parsedData); err == nil {
					response.Rows = append(response.Rows, TableRow{
						Data:  parsedData,
						Payer: keyValue.Payer.String(),
					})
				} else {
					log15.Error("failed to encode json", "err", err)
				}
			} else {
				log15.Error("failed to decode struct", "err", err)
			}
		}

		c.JSON(200, response)
	}
}
