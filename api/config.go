package prismasdk2

import (
	"context"
	"net/http"
)

// Config config
type Config struct {
	API           string
	Namespace     string
	HTTPClient    *http.Client
	TokenProvider TokenProvider
}

// NewConfig returns new Config
func NewConfig() *Config {
	return &Config{}
}

// SetAPI sets attribute and returns self
func (t *Config) SetAPI(v string) *Config {
	t.API = v
	return t
}

// SetNamespace sets attribute and returns self
func (t *Config) SetNamespace(v string) *Config {
	t.Namespace = v
	return t
}

// SetTokenProvider sets interface and returns self
func (t *Config) SetTokenProvider(v TokenProvider) *Config {
	t.TokenProvider = v
	return t
}

// SetHTTPClient sets entity and returns self
func (t *Config) SetHTTPClient(v *http.Client) *Config {
	t.HTTPClient = v
	return t
}

// Build returns entity
func (t *Config) Build(ctx context.Context) (*Client, error) {
	return NewClient(ctx, t)
}
