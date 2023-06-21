package history_api_plugin

import "github.com/MetalBlockchain/antelopevm/core/name"

type GetKeyAccountsResponse struct {
	AccountNames []name.AccountName `json:"account_names"`
}
