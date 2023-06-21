package service

type GetActionsRequest struct {
	AccountName string `json:"account_name"`
}

type GetActionsResponse struct {
	Actions                  []string `json:"actions"`
	HeadBlockNum             int      `json:"head_block_num"`
	LastIrreversibleBlockNum int      `json:"last_irreversible_block"`
}
