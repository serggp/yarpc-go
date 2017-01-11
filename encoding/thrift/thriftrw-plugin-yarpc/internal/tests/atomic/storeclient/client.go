// Code generated by thriftrw-plugin-yarpc
// @generated

package storeclient

import (
	"context"
	"go.uber.org/thriftrw/wire"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/encoding/thrift"
	"go.uber.org/yarpc/encoding/thrift/thriftrw-plugin-yarpc/internal/tests/atomic"
	"go.uber.org/yarpc/encoding/thrift/thriftrw-plugin-yarpc/internal/tests/common/baseserviceclient"
	"go.uber.org/yarpc"
)

// Interface is a client for the Store service.
type Interface interface {
	baseserviceclient.Interface

	CompareAndSwap(
		ctx context.Context,
		Request *atomic.CompareAndSwap,
		opts ...yarpc.CallOption,
	) error

	Forget(
		ctx context.Context,
		Key *string,
		opts ...yarpc.CallOption,
	) (yarpc.Ack, error)

	Increment(
		ctx context.Context,
		Key *string,
		Value *int64,
		opts ...yarpc.CallOption,
	) error

	Integer(
		ctx context.Context,
		Key *string,
		opts ...yarpc.CallOption,
	) (int64, error)
}

// New builds a new client for the Store service.
//
// 	client := storeclient.New(dispatcher.ClientConfig("store"))
func New(c transport.ClientConfig, opts ...thrift.ClientOption) Interface {
	return client{
		c: thrift.New(thrift.Config{
			Service:      "Store",
			ClientConfig: c,
		}, opts...),
		Interface: baseserviceclient.New(c),
	}
}

func init() {
	yarpc.RegisterClientBuilder(func(c transport.ClientConfig) Interface {
		return New(c)
	})
}

type client struct {
	baseserviceclient.Interface

	c thrift.Client
}

func (c client) CompareAndSwap(
	ctx context.Context,
	_Request *atomic.CompareAndSwap,
	opts ...yarpc.CallOption,
) (err error) {

	args := atomic.Store_CompareAndSwap_Helper.Args(_Request)

	var body wire.Value
	body, err = c.c.Call(ctx, args, opts...)
	if err != nil {
		return
	}

	var result atomic.Store_CompareAndSwap_Result
	if err = result.FromWire(body); err != nil {
		return
	}

	err = atomic.Store_CompareAndSwap_Helper.UnwrapResponse(&result)
	return
}

func (c client) Forget(
	ctx context.Context,
	_Key *string,
	opts ...yarpc.CallOption,
) (yarpc.Ack, error) {
	args := atomic.Store_Forget_Helper.Args(_Key)
	return c.c.CallOneway(ctx, args, opts...)
}

func (c client) Increment(
	ctx context.Context,
	_Key *string,
	_Value *int64,
	opts ...yarpc.CallOption,
) (err error) {

	args := atomic.Store_Increment_Helper.Args(_Key, _Value)

	var body wire.Value
	body, err = c.c.Call(ctx, args, opts...)
	if err != nil {
		return
	}

	var result atomic.Store_Increment_Result
	if err = result.FromWire(body); err != nil {
		return
	}

	err = atomic.Store_Increment_Helper.UnwrapResponse(&result)
	return
}

func (c client) Integer(
	ctx context.Context,
	_Key *string,
	opts ...yarpc.CallOption,
) (success int64, err error) {

	args := atomic.Store_Integer_Helper.Args(_Key)

	var body wire.Value
	body, err = c.c.Call(ctx, args, opts...)
	if err != nil {
		return
	}

	var result atomic.Store_Integer_Result
	if err = result.FromWire(body); err != nil {
		return
	}

	success, err = atomic.Store_Integer_Helper.UnwrapResponse(&result)
	return
}
