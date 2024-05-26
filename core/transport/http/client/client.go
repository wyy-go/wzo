package client

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/wyy-go/wencoding"
	"github.com/wyy-go/wzo/core/errors"
	"github.com/wyy-go/wzo/core/registry"
	"golang.org/x/oauth2"
)

var noRequestBodyMethods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodHead:    {},
	http.MethodDelete:  {},
	http.MethodConnect: {},
	http.MethodOptions: {},
	http.MethodTrace:   {},
}

type Client struct {
	cc          *resty.Client
	serviceName string
	serviceAddr string
	registry    registry.Registry
	mws         []Middleware
	balancer    Balancer
	resolver    *Resolver

	codec *wencoding.Encoding
	// A TokenSource is anything that can return a token.
	tokenSource oauth2.TokenSource
	// validate request
	validate func(any) error
	// call option
	callOptions []CallOption
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		cc:       resty.New(),
		codec:    wencoding.New(),
		mws:      []Middleware{},
		balancer: &Wrr{},
	}
	for _, opt := range opts {
		opt.apply(c)
	}

	if c.registry != nil {
		c.resolver = newResolver(c.registry, c.balancer, &Target{
			Scheme:    "http",
			Authority: "",
			Endpoint:  c.serviceName,
		})
	}

	c.cc.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		if r.RawResponse != nil {
			body := r.RawResponse.Body
			defer body.Close()
			r.RawResponse.Body = io.NopCloser(bytes.NewBuffer(r.Body()))
		}
		return nil
	})
	return c
}

func (c *Client) Use(mws ...Middleware) *Client {
	c.mws = append(c.mws, mws...)
	return c
}

func (c *Client) Deref() *resty.Client { return c.cc }

func (c *Client) CallSetting(path string, cos ...CallOption) *CallSettings {
	cs := &CallSettings{
		contentType: "application/json",
		accept:      "application/json",
		Path:        path,
		header:      make(http.Header),
		noAuth:      false,
	}
	for _, co := range c.callOptions {
		co(cs)
	}
	for _, co := range cos {
		co(cs)
	}
	return cs
}

// Invoke2 the request
// NOTE: Do not use this function. use Execute instead.
func (c *Client) Invoke2(ctx context.Context, method, path string, in, out any, settings *CallSettings) error {
	if c.validate != nil {
		err := c.validate(in)
		if err != nil {
			return err
		}
	}
	ctx = WithValueCallOption(ctx, settings)
	r := c.cc.R().SetContext(ctx)
	if in != nil {
		reqBody, err := c.codec.Encode(settings.contentType, in)
		if err != nil {
			return err
		}
		r = r.SetBody(reqBody)
	}
	if !settings.noAuth {
		if c.tokenSource == nil {
			return errors.Parse("transport: token source should be not nil")
		}
		tk, err := c.tokenSource.Token()
		if err != nil {
			return err
		}
		r.SetHeader("Authorization", tk.Type()+" "+tk.AccessToken)
	}
	r.SetHeader("Content-Type", settings.contentType)
	r.SetHeader("Accept", settings.accept)
	for k, vs := range settings.header {
		for _, v := range vs {
			r.Header.Add(k, v)
		}
	}

	resp, err := r.Execute(method, c.cc.BaseURL+path)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return errors.Parse(string(resp.Body()))
	}
	defer resp.RawResponse.Body.Close()
	return c.codec.InboundForResponse(resp.RawResponse).NewDecoder(resp.RawResponse.Body).Decode(out)
}

// Invoke the request
// NOTE: Do not use this function. use Execute instead.
func (c *Client) Invoke(ctx context.Context, method, path string, in, out any, settings *CallSettings) error {
	url := "http://"
	if c.serviceAddr != "" {
		url += c.serviceAddr
	} else {
		node, err := c.balancer.Select()
		if err != nil {
			return err
		}
		url += node.Addr
	}

	url += path

	h := func(ctx context.Context, in any) (any, error) {
		if c.validate != nil {
			err := c.validate(in)
			if err != nil {
				return nil, err
			}
		}
		ctx = WithValueCallOption(ctx, settings)
		r := c.cc.R().SetContext(ctx)
		if in != nil {
			reqBody, err := c.codec.Encode(settings.contentType, in)
			if err != nil {
				return nil, err
			}
			r = r.SetBody(reqBody)
		}
		if !settings.noAuth {
			if c.tokenSource == nil {
				return nil, errors.Parse("transport: token source should be not nil")
			}
			tk, err := c.tokenSource.Token()
			if err != nil {
				return nil, err
			}
			r.SetHeader("Authorization", tk.Type()+" "+tk.AccessToken)
		}
		r.SetHeader("Content-Type", settings.contentType)
		r.SetHeader("Accept", settings.accept)
		for k, vs := range settings.header {
			for _, v := range vs {
				r.Header.Add(k, v)
			}
		}

		resp, err := r.Execute(method, c.cc.BaseURL+path)
		if err != nil {
			return nil, err
		}
		if resp.IsError() {
			return nil, errors.Parse(string(resp.Body()))
		}
		defer resp.RawResponse.Body.Close()
		return nil, c.codec.InboundForResponse(resp.RawResponse).NewDecoder(resp.RawResponse.Body).Decode(out)
	}

	if len(c.mws) > 0 {
		Chain(c.mws...)(h)
	}

	_, err := h(ctx, in)

	return err

}

func hasRequestBody(method string) bool {
	_, ok := noRequestBodyMethods[method]
	return !ok
}

func (c *Client) Execute(ctx context.Context, method, path string, req, resp any, opts ...CallOption) error {
	var r any

	settings := c.CallSetting(path, opts...)
	hasBody := hasRequestBody(method)
	if hasBody {
		r = req
	}
	url := c.EncodeURL(settings.Path, req, !hasBody)
	return c.Invoke(ctx, method, url, r, &resp, settings)
}

// Get method does GET HTTP request. It's defined in section 4.3.1 of RFC7231.
func (c *Client) Get(ctx context.Context, path string, req, resp any, opts ...CallOption) error {
	return c.Execute(ctx, http.MethodGet, path, req, resp, opts...)
}

// Head method does HEAD HTTP request. It's defined in section 4.3.2 of RFC7231.
func (c *Client) Head(ctx context.Context, path string, req, resp any, opts ...CallOption) error {
	return c.Execute(ctx, http.MethodHead, path, req, resp, opts...)
}

// Post method does POST HTTP request. It's defined in section 4.3.3 of RFC7231.
func (c *Client) Post(ctx context.Context, path string, req, resp any, opts ...CallOption) error {
	return c.Execute(ctx, http.MethodPost, path, req, resp, opts...)
}

// Put method does PUT HTTP request. It's defined in section 4.3.4 of RFC7231.
func (c *Client) Put(ctx context.Context, path string, req, resp any, opts ...CallOption) error {
	return c.Execute(ctx, http.MethodPut, path, req, resp, opts...)
}

// Delete method does DELETE HTTP request. It's defined in section 4.3.5 of RFC7231.
func (c *Client) Delete(ctx context.Context, path string, req, resp any, opts ...CallOption) error {
	return c.Execute(ctx, http.MethodDelete, path, req, resp, opts...)
}

// Options method does OPTIONS HTTP request. It's defined in section 4.3.7 of RFC7231.
func (c *Client) Options(ctx context.Context, path string, req, resp any, opts ...CallOption) error {
	return c.Execute(ctx, http.MethodOptions, path, req, resp, opts...)
}

// Patch method does PATCH HTTP request. It's defined in section 2 of RFC5789.
func (c *Client) Patch(ctx context.Context, path string, req, resp any, opts ...CallOption) error {
	return c.Execute(ctx, http.MethodPatch, path, req, resp, opts...)
}

// EncodeURL encode msg to url path.
// pathTemplate is a template of url path like http://helloworld.dev/{name}/sub/{sub.name}.
func (c *Client) EncodeURL(pathTemplate string, msg any, needQuery bool) string {
	return c.codec.EncodeURL(pathTemplate, msg, needQuery)
}

// EncodeQuery encode v into “URL encoded” form
// ("bar=baz&foo=quux") sorted by key.
func (c *Client) EncodeQuery(v any) (string, error) {
	vv, err := c.codec.EncodeQuery(v)
	if err != nil {
		return "", err
	}
	return vv.Encode(), nil
}
