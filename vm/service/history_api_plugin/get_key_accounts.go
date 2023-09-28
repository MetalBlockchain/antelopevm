package history_api_plugin

import "github.com/MetalBlockchain/antelopevm/chain/name"

type GetKeyAccountsResponse struct {
	AccountNames []name.AccountName `json:"account_names"`
}
