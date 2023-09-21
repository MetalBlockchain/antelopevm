package chain_api_plugin

import (
	"encoding/json"
	"net/http"

	"github.com/MetalBlockchain/antelopevm/core"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
)

type RequiredKeysRequest struct {
	Transaction   core.Transaction `json:"transaction"`
	AvailableKeys []ecc.PublicKey  `json:"available_keys"`
}

type RequiredKeysResponse struct {
	RequiredKeys []ecc.PublicKey `json:"required_keys"`
}

func init() {
	service.RegisterHandler("/v1/chain/get_required_keys", service.Handler{
		Methods:     []string{http.MethodPost},
		HandlerFunc: GetRequiredKeys,
	})
}

func GetRequiredKeys(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body RequiredKeysRequest
		json.NewDecoder(c.Request.Body).Decode(&body)
		session := vm.GetState().CreateSession(false)
		defer session.Discard()
		authorizationManager := vm.GetController().GetAuthorizationManager(session)
		keySet := ecc.NewPublicKeySetFromArray(body.AvailableKeys)
		data, err := authorizationManager.GetRequiredKeys(body.Transaction, keySet)

		if err != nil {
			c.JSON(400, err)
			return
		}

		c.JSON(200, RequiredKeysResponse{RequiredKeys: data})
	}
}
