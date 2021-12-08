package token

import (
	"net/http"

	"go.uber.org/zap"
)

// Config config
type Config struct {
	API             string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	HTTPClient      *http.Client
}

// NewConfig returns new Config
func NewConfig() *Config {
	return &Config{}
}

// SetAPI sets attribute and returns self
func (t *Config) SetAPI(api string) *Config {
	t.API = api
	return t
}

// SetAccessKeyID sets attribute and returns self
func (t *Config) SetAccessKeyID(accessKeyID string) *Config {
	t.AccessKeyID = accessKeyID
	return t
}

// SetSecretAccessKey sets attribute and returns self
func (t *Config) SetSecretAccessKey(secretAccessKey string) *Config {
	t.SecretAccessKey = secretAccessKey
	return t
}

// SetSessionToken sets attribute and returns self
func (t *Config) SetSessionToken(sessionToken string) *Config {
	t.SessionToken = sessionToken
	return t
}

// SetHTTPClient sets entity and returns self
func (t *Config) SetHTTPClient(httpClient *http.Client) *Config {
	t.HTTPClient = httpClient
	return t
}

// GetHTTPClient returns entity. If entity is nil entity will be initialized and returned.
func (t *Config) GetHTTPClient() *http.Client {

	if t.HTTPClient == nil {
		t.HTTPClient = &http.Client{}
		zap.L().Debug("HTTPClient created new")
	} else {
		zap.L().Debug("HTTPClient set from config")
	}

	return t.HTTPClient
}

// Build returns entity
func (t *Config) Build() (*Client, error) {
	return NewClient(t)
}
