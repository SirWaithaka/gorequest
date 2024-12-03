package client

import "net/http"

type (
	// Config is a configuration param type for Client
	Config struct {
		BaseURL    string
		HTTPClient *http.Client
	}
)

func (c *Config) MergeIn(other *Config) {
	if other == nil {
		return
	}

	if other.HTTPClient != nil {
		c.HTTPClient = other.HTTPClient
	}

}
