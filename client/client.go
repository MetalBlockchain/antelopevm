// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package client

import (
	"context"

	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils/rpc"
)

// Client defines timestampvm client operations.
type Client interface {
	// ProposeBlock submits data for a block
	ProposeBlock(ctx context.Context) (bool, error)

	// GetBlock fetches the contents of a block
	GetBlock(ctx context.Context, blockID *ids.ID) (bool, error)
}

// New creates a new client object.
func New(uri string) Client {
	req := rpc.NewEndpointRequester(uri)
	return &client{req: req}
}

type client struct {
	req rpc.EndpointRequester
}

func (cli *client) ProposeBlock(ctx context.Context) (bool, error) {
	return true, nil
}

func (cli *client) GetBlock(ctx context.Context, blockID *ids.ID) (bool, error) {
	return true, nil
}
