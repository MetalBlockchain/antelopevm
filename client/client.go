// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package client

import (
	"context"
)

// Client defines timestampvm client operations.
type Client interface {
	// ProposeBlock submits data for a block
	PushTransaction(ctx context.Context) (bool, error)
}

// New creates a new client object.
func New(uri string) Client {
	return &client{uri: uri}
}

type client struct {
	uri string
}

func (cli *client) PushTransaction(ctx context.Context) (bool, error) {
	return true, nil
}
