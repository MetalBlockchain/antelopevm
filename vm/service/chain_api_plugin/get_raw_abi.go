package chain_api_plugin

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/crypto"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
)

type GetRawAbiRequest struct {
	AccountName string `json:"account_name"`
}

type GetRawAbiResponse struct {
	AccountName string        `json:"account_name"`
	CodeHash    crypto.Sha256 `json:"code_hash"`
	AbiHash     crypto.Sha256 `json:"abi_hash"`
	Abi         string        `json:"abi"`
}

func init() {
	service.RegisterHandler("/v1/chain/get_raw_abi", service.Handler{
		Methods:     []string{http.MethodPost},
		HandlerFunc: GetRawAbi,
	})
}

func GetRawAbi(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body GetRawAbiRequest
		json.NewDecoder(c.Request.Body).Decode(&body)
		session := vm.GetState().CreateSession(false)
		defer session.Discard()
		acc, err := session.FindAccountByName(name.StringToName(body.AccountName))

		if err != nil {
			c.JSON(400, "account not found")
			return
		}

		if len(acc.Abi) == 0 {
			response := GetRawAbiResponse{
				AccountName: acc.Name.String(),
			}

			c.JSON(200, response)
			return
		}

		rawAbi := base64.StdEncoding.EncodeToString(acc.Abi)
		response := GetRawAbiResponse{
			AccountName: acc.Name.String(),
			CodeHash:    acc.CodeVersion,
			AbiHash:     *crypto.NewSha256String("bf13acab1b4bc2676ef6f0afcf1765ab5db3ffa1ac18453a628e1e65fe26e045"),
			Abi:         rawAbi,
		}
		data, err := json.Marshal(response)
		c.Writer.Write(data)
		c.Writer.Header().Set("Content-Length", strconv.Itoa(len(data)))
	}
}
