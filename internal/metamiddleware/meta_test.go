package metamiddleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/yarpc/api/middleware"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/yarpctest"
)

func TestMetaOutboundMidleware(t *testing.T) {
	const overrideName = "override"

	out := yarpctest.NewFakeTransport().NewOutbound(nil,
		yarpctest.OutboundName(overrideName),
		yarpctest.OutboundCallOverride(
			func(context.Context, *transport.Request) (*transport.Response, error) { return nil, nil },
		),
		yarpctest.OutboundCallStreamOverride(
			func(context.Context, *transport.StreamRequest) (*transport.ClientStream, error) { return nil, nil },
		),
		yarpctest.OutboundCallOnewayOverride(
			func(context.Context, *transport.Request) (transport.Ack, error) { return nil, nil },
		),
	)

	t.Run("unary", func(t *testing.T) {
		req := &transport.Request{Transport: "" /* not set */}

		outWithMiddleware := middleware.ApplyUnaryOutbound(out, New())
		_, err := outWithMiddleware.Call(context.Background(), req)
		require.NoError(t, err)

		assert.Equal(t, overrideName, string(req.Transport))
	})

	t.Run("oneway", func(t *testing.T) {
		req := &transport.Request{Transport: "" /* not set */}

		outWithMiddleware := middleware.ApplyOnewayOutbound(out, New())
		_, err := outWithMiddleware.CallOneway(context.Background(), req)
		require.NoError(t, err)

		assert.Equal(t, overrideName, string(req.Transport))
	})

	t.Run("stream", func(t *testing.T) {
		streamReq := &transport.StreamRequest{Meta: &transport.RequestMeta{Transport: "" /* not set */}}

		outWithMiddleware := middleware.ApplyStreamOutbound(out, New())
		_, err := outWithMiddleware.CallStream(context.Background(), streamReq)
		require.NoError(t, err)

		assert.Equal(t, overrideName, string(streamReq.Meta.Transport))
	})
}
