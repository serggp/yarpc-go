// Copyright (c) 2019 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package encoding contains helper functionality for encoding implementations.
package encoding

import (
	"context"

	"go.uber.org/yarpc"
	"go.uber.org/yarpc/api/encoding"
	"go.uber.org/yarpc/api/transport"
)

// AdapterClient is a client encoding implementations can use to support wire
// encodings that differ from application used types.
type AdapterClient interface {
	Call(context.Context, *transport.Request, interface{}, encoding.Adapter, ...yarpc.CallOption) (interface{}, error)
}

type adapterClient struct {
	cc transport.ClientConfig
}

// AdapterClientConfig contains the configuration for the AdapterClient.
type AdapterClientConfig struct {
	ClientConfig transport.ClientConfig
}

// NewAdapterClient returns a new AdapterClient, suitable for an encoding
// implementation, for example Thrift.
//
// This only supports unary calls.
func NewAdapterClient(c AdapterClientConfig) (AdapterClient, error) {
	return &adapterClient{
		cc: c.ClientConfig,
	}, nil
}

func (c *adapterClient) Call(
	ctx context.Context,
	tReq *transport.Request,
	req interface{},
	adapter encoding.Adapter,
	opts ...yarpc.CallOption,
) (interface{}, error) {
	call := encoding.NewOutboundCall(FromOptions(opts)...)
	ctx, err := call.WriteToRequest(ctx, tReq)
	if err != nil {
		return nil, err
	}
	tReq, err = adapter.ToRequest(ctx, tReq, req)
	if err != nil {
		return nil, err
	}
	res, err := c.cc.GetUnaryOutbound().Call(ctx, tReq)
	if err != nil {
		return nil, err
	}
	ctx, err = call.ReadFromResponse(ctx, res)
	if err != nil {
		return nil, err
	}
	return adapter.FromResponse(ctx, res)
}
