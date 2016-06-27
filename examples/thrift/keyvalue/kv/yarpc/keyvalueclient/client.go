// Code generated by thriftrw

// Copyright (c) 2016 Uber Technologies, Inc.
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

package keyvalueclient

import (
	"github.com/thriftrw/thriftrw-go/protocol"
	"github.com/thriftrw/thriftrw-go/wire"
	yarpc "github.com/yarpc/yarpc-go"
	"github.com/yarpc/yarpc-go/encoding/thrift"
	"github.com/yarpc/yarpc-go/examples/thrift/keyvalue/kv/service/keyvalue"
	"github.com/yarpc/yarpc-go/transport"
)

type Interface interface {
	GetValue(reqMeta yarpc.CallReqMeta, key *string) (string, yarpc.CallResMeta, error)
	SetValue(reqMeta yarpc.CallReqMeta, key *string, value *string) (yarpc.CallResMeta, error)
}

func New(c transport.Channel) Interface {
	return client{c: thrift.New(thrift.Config{Service: "KeyValue", Channel: c, Protocol: protocol.Binary})}
}

type client struct{ c thrift.Client }

func (c client) GetValue(reqMeta yarpc.CallReqMeta, key *string) (success string, resMeta yarpc.CallResMeta, err error) {
	args := keyvalue.GetValueHelper.Args(key)
	var w wire.Value
	w, err = args.ToWire()
	if err != nil {
		return
	}
	var body wire.Value
	body, resMeta, err = c.c.Call("getValue", reqMeta, w)
	if err != nil {
		return
	}
	var result keyvalue.GetValueResult
	if err = result.FromWire(body); err != nil {
		return
	}
	success, err = keyvalue.GetValueHelper.UnwrapResponse(&result)
	return
}

func (c client) SetValue(reqMeta yarpc.CallReqMeta, key *string, value *string) (resMeta yarpc.CallResMeta, err error) {
	args := keyvalue.SetValueHelper.Args(key, value)
	var w wire.Value
	w, err = args.ToWire()
	if err != nil {
		return
	}
	var body wire.Value
	body, resMeta, err = c.c.Call("setValue", reqMeta, w)
	if err != nil {
		return
	}
	var result keyvalue.SetValueResult
	if err = result.FromWire(body); err != nil {
		return
	}
	err = keyvalue.SetValueHelper.UnwrapResponse(&result)
	return
}
