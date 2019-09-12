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
	"context"

	"go.uber.org/yarpc/api/transport"
)

// AdapterProvider returns an Adapter for a procedure.
type AdapterProvider interface {
	Adapter(procedureName string) (_ Adapter, ok bool)
}

// Adapter is a procedure-scoped type used by encoding implementations to
// transparently change the encoding used on the wire for outbound requests.
//
// This is useful, for example, an aggressive Thrift to Protobuf migration;
// Thrift users wanting Protobuf on the wire would not need to change any code.
type Adapter interface {
	ToRequest(context.Context, *transport.Request, interface{}) (*transport.Request, error)
	FromResponse(context.Context, *transport.Response) (interface{}, error)
}

// NopAdapterProvider is a no-op implementation of AdapterProvider that always
// returns false.
var NopAdapterProvider AdapterProvider = nopAdapterProvider{}

type nopAdapterProvider struct{}

func (nopAdapterProvider) Adapter(string) (Adapter, bool) { return nil, false }
