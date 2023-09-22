package chain_api_plugin

import (
	"net/http"

	"github.com/MetalBlockchain/antelopevm/core/protocol"
	"github.com/MetalBlockchain/antelopevm/vm/service"
	"github.com/gin-gonic/gin"
)

type GetActivatedProtocolFeaturesResults struct {
	ActivatedProtocolFeatures []protocol.ProtocolFeature `json:"activated_protocol_features"`
	More                      uint32                     `json:"more"`
}

func init() {
	service.RegisterHandler("/v1/chain/get_activated_protocol_features", service.Handler{
		Methods:     []string{http.MethodGet, http.MethodPost},
		HandlerFunc: GetActivatedProtocolFeatures,
	})
}

func GetActivatedProtocolFeatures(vm service.VM) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := GetActivatedProtocolFeaturesResults{
			ActivatedProtocolFeatures: make([]protocol.ProtocolFeature, 0),
			More:                      0,
		}
		c.JSON(200, result)
	}
}
