package tracehttp

import (
	"context"
	"io"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const operationName = "Request"

type Client struct {
	ctx    context.Context
	Client *http.Client
	Header http.Header
}

func (c *Client) WithContext(ctx context.Context) *Client {
	c.ctx = ctx
	return c
}

// Example:
// resp, _ := tracehttp.Get(ctx, "http://www.baidu.com")
// defer resp.Body.Close()
// _, _ = ioutil.ReadAll(resp.Body)

// Get issues a GET to the specified URL. If the response is one of
// the following redirect codes, Get follows the redirect, up to a
// maximum of 10 redirects:
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
//    308 (Permanent Redirect)
//
// An error is returned if there were too many redirects or if there
// was an HTTP protocol error. A non-2xx response doesn't cause an
// error. Any returned error will be of type *url.Error. The url.Error
// value's Timeout method will report true if request timed out or was
// canceled.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.
//
// Get is a wrapper around DefaultClient.Get.
//
// To make a request with custom headers, use NewRequest and
// DefaultClient.Do.
func Get(ctx context.Context, url string) (resp *http.Response, err error) {
	var defaultClient = &Client{ctx: ctx}
	return defaultClient.Get(url)
}

// Get issues a GET to the specified URL. If the response is one of the
// following redirect codes, Get follows the redirect after calling the
// Client's CheckRedirect function:
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
//    308 (Permanent Redirect)
//
// An error is returned if the Client's CheckRedirect function fails
// or if there was an HTTP protocol error. A non-2xx response doesn't
// cause an error. Any returned error will be of type *url.Error. The
// url.Error value's Timeout method will report true if request timed
// out or was canceled.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.
//
// To make a request with custom headers, use NewRequest and Client.Do.
func (c *Client) Get(url string) (resp *http.Response, err error) {
	resp, err = c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Post issues a POST to the specified URL.
//
// Caller should close resp.Body when done reading from it.
//
// If the provided body is an io.Closer, it is closed after the
// request.
//
// Post is a wrapper around DefaultClient.Post.
//
// To set custom headers, use NewRequest and DefaultClient.Do.
//
// See the Client.Do method documentation for details on how redirects
// are handled.
func Post(ctx context.Context, url, contentType string, body io.Reader) (resp *http.Response, err error) {
	var defaultClient = &Client{ctx: ctx}
	return defaultClient.Post(url, contentType, body)
}

// Post issues a POST to the specified URL.
//
// Caller should close resp.Body when done reading from it.
//
// If the provided body is an io.Closer, it is closed after the
// request.
//
// To set custom headers, use NewRequest and Client.Do.
//
// See the Client.Do method documentation for details on how redirects
// are handled.
func (c *Client) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	resp, err = c.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	c.Header.Set("Content-Type", contentType)
	return resp, nil
}

// NewRequest
func (c *Client) NewRequest(method, url string, body io.Reader) (resp *http.Response, err error) {
	// 记录一个 span
	span, ctx := opentracing.StartSpanFromContext(c.ctx, operationName)
	defer span.Finish()

	if c.Client != nil {
		c.Client.Transport = &nethttp.Transport{}
	} else {
		c.Client = &http.Client{Transport: &nethttp.Transport{}}
	}

	ext.HTTPMethod.Set(span, method)
	ext.HTTPUrl.Set(span, url)

	ctx = opentracing.ContextWithSpan(ctx, span)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ext.Error.Set(span, true)
		span.SetTag("error", err)
		return
	}
	// set http client header
	if len(c.Header) > 0 {
		req.Header = c.Header
	}

	req = req.WithContext(ctx)
	// wrap the request in nethttp.TraceRequest
	req, ht := nethttp.TraceRequest(span.Tracer(), req)
	defer ht.Finish()

	resp, err = c.Client.Do(req)
	if err != nil {
		ext.Error.Set(span, true)
		span.SetTag("error", err)
		return
	}

	ext.HTTPStatusCode.Set(span, uint16(resp.StatusCode))
	ext.PeerHostname.Set(span, req.Host)
	return resp, err
}
