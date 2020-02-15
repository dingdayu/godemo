package restyhttp

import (
	"context"
	"io"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"demo/pkg/jaeger"
	"demo/pkg/log"
)

const operationName = "Request"

type Client struct {
	ctx    context.Context
	Client *resty.Client
	Header http.Header
	logger *log.Logger
}

func NewClient(ctx context.Context) *Client {
	span, ctx := jaeger.StartSpanFromContext(ctx, "NewHTTP")
	defer span.Finish()

	client := new(Client)
	client.ctx = ctx
	client.Client = resty.New()
	client.logger = log.New().Named("HTTP.GET")
	return client
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

	ext.HTTPMethod.Set(span, method)
	ext.HTTPUrl.Set(span, url)
	ctx = opentracing.ContextWithSpan(ctx, span)

	if c.Client == nil {
		c.Client = resty.New()
	}

	// set http client header
	if len(c.Header) > 0 {
		c.Client.Header = c.Header
	}
	jaeger.InjectSpanToHeader(span, c.Client.Header)

	req := c.Client.NewRequest().SetBody(body)
	tryResp, err := req.Execute(method, url)
	if err != nil {
		ext.Error.Set(span, true)
		span.SetTag("error", err)
		return
	}

	ext.HTTPStatusCode.Set(span, uint16(tryResp.StatusCode()))
	ext.PeerHostname.Set(span, req.RawRequest.Host)
	return tryResp.RawResponse, err
}
