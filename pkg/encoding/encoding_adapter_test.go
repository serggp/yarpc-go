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

package encoding

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/api/encoding"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/yarpctest"
)

type testRequest struct {
	a, b int
}
type testResponse struct {
	sum string
}

func TestAdapterClientRoundTrip(t *testing.T) {
	const (
		procedure  = "test-procedure"
		routingKey = "test-routing-key"
	)

	// This fake outbound mocks a server and handler.
	outbound := yarpctest.NewFakeTransport().NewOutbound(nil, yarpctest.OutboundCallOverride(
		yarpctest.OutboundCallable(func(ctx context.Context, req *transport.Request) (*transport.Response, error) {
			// simulate inbound encoding
			ctx, inboundCall := encoding.NewInboundCall(context.Background())
			assert.NoError(t, inboundCall.ReadFromRequest(req))

			// simulate a handler, validate metadata
			call := yarpc.CallFromContext(ctx)
			assert.Equal(t, procedure+"-adapted", call.Procedure(), "metadata is missing from request")
			assert.Equal(t, routingKey, call.RoutingKey(), "metadata is missing from request")

			// return request body as response body
			return &transport.Response{
				Body: ioutil.NopCloser(req.Body),
			}, nil
		}),
	))

	client, err := NewAdapterClient(AdapterClientConfig{
		ClientConfig: &staticClientConfig{outbound: outbound},
	})
	require.NoError(t, err)

	tReq := &transport.Request{Procedure: procedure}
	request := &testRequest{a: 12, b: 8}
	adapter := &adapter{}

	result, err := client.Call(context.Background(), tReq, request, adapter, yarpc.WithRoutingKey(routingKey))
	require.NoError(t, err)

	response, ok := result.(*testResponse)
	require.True(t, ok, "expected '*testResponse', got %T", result)

	assert.Equal(t, strconv.Itoa(request.a+request.b), response.sum)
}

// adapter will sum both fields of *testRequest and put it into
// *testResponse.sum.
type adapter struct{}

func (a *adapter) ToRequest(ctx context.Context, tReq *transport.Request, request interface{}) (*transport.Request, error) {
	tReq.Procedure = tReq.Procedure + "-adapted"
	testReq, ok := request.(*testRequest)
	if !ok {
		return nil, fmt.Errorf("unexpected type: %T", testReq)
	}

	// sum both fields and marshaling it as a string
	tReq.Body = bytes.NewBufferString(strconv.Itoa(testReq.a + testReq.b))
	return tReq, nil
}

func (a *adapter) FromResponse(ctx context.Context, res *transport.Response) (interface{}, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// unmarshal string it into a *testResponse.
	return &testResponse{sum: string(body)}, nil
}

var _ transport.ClientConfig = (*staticClientConfig)(nil)

type staticClientConfig struct {
	outbound *yarpctest.FakeOutbound
}

func (*staticClientConfig) Caller() string {
	return ""
}

func (*staticClientConfig) Service() string {
	return ""
}

func (c *staticClientConfig) GetUnaryOutbound() transport.UnaryOutbound {
	return c.outbound
}

func (*staticClientConfig) GetOnewayOutbound() transport.OnewayOutbound {
	panic("unimplmented")
}
