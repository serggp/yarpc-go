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

package raw

import (
	"io/ioutil"

	"github.com/yarpc/yarpc-go/internal/meta"
	"github.com/yarpc/yarpc-go/transport"

	"golang.org/x/net/context"
)

// rawHandler adapts a Handler into a transport.Handler
type rawHandler struct {
	h Handler
}

func (r rawHandler) Handle(ctx context.Context, treq *transport.Request, rw transport.ResponseWriter) error {
	treq.Encoding = Encoding
	// TODO(abg): Should we fail requests if Rpc-Encoding does not match?

	reqBody, err := ioutil.ReadAll(treq.Body)
	if err != nil {
		return err
	}

	reqMeta := meta.FromTransportRequest(ctx, treq)
	resBody, resMeta, err := r.h(reqMeta, reqBody)
	if err != nil {
		return err
	}

	if resMeta != nil {
		_ = meta.ToTransportResponseWriter(resMeta, rw)
		// TODO(abg): Propagate response context
	}

	if _, err := rw.Write(resBody); err != nil {
		return err
	}

	return nil
}
