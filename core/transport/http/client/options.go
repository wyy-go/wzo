package client

import (
	"github.com/wyy-go/wencoding"
	"github.com/wyy-go/wzo/core/registry"
	"golang.org/x/oauth2"
)

type Option interface {
	apply(c *Client)
}

type optionFunc func(c *Client)

type ClientOption func(*Client)

func (f optionFunc) apply(c *Client) {
	f(c)
}

func WithServiceName(name string) Option {
	return optionFunc(func(c *Client) {
		c.serviceName = name
	})
}

func WithServiceAddr(addr string) Option {
	return optionFunc(func(c *Client) {
		c.serviceAddr = addr
	})
}

func Registry(r registry.Registry) Option {
	return optionFunc(func(c *Client) {
		c.registry = r
	})
}

func WithEncoding(codec *wencoding.Encoding) Option {
	return optionFunc(func(c *Client) {
		c.codec = codec
	})
}

func WithTokenSource(t oauth2.TokenSource) Option {
	return optionFunc(func(c *Client) {
		c.tokenSource = t
	})
}

func WithValidate(f func(any) error) Option {
	return optionFunc(func(c *Client) {
		c.validate = f
	})
}

// WithCallOption WithCallOption(WithCoNoAuth)
func WithCallOption(co ...CallOption) Option {
	return optionFunc(func(c *Client) {
		c.callOptions = append(c.callOptions, co...)
	})
}
