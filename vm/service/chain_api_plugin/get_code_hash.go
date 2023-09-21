package chain_api_plugin

import (
	"encoding/json"
	"net/http"

	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
)

type GetCodeHashRequest struct {
	AccountName string `json:"account_name"`
}

type GetCodeHashResponse struct {
	AccountName string        `json:"account_name"`
	CodeHash    crypto.Sha256 `json:"code_hash"`
}

func NewGetCodeHashResponse(accountName string, codeHash crypto.Sha256) GetCodeHashResponse {
	return GetCodeHashResponse{
		AccountName: accountName,
		CodeHash:    codeHash,
	}
}

func init() {
	service.RegisterHandler("/v1/chain/get_code_hash", service.Handler{
		Methods:     []string{http.MethodPost},
		HandlerFunc: GetCodeHash,
	})
}

func GetCodeHash(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body GetCodeHashRequest
		json.NewDecoder(c.Request.Body).Decode(&body)
		session := vm.GetState().CreateSession(false)
		defer session.Discard()
		acc, err := session.FindAccountByName(name.StringToName(body.AccountName))

		if err != nil {
			c.JSON(400, "account not found")
			return
		}

		c.JSON(200, GetCodeHashResponse{AccountName: body.AccountName, CodeHash: acc.CodeVersion})
	}
}
