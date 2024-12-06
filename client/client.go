package client

import (
	"net/http"
)

type (
	// Client is an HTTP client.
	Client struct {
		config Config
	}
)

func New(cfg Config) *Client {

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}

	return &Client{config: cfg}
}
