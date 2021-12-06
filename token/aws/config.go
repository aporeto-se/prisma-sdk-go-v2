package token

import (
	"net/http"

	"go.uber.org/zap"
)

// Config config
type Config struct {
	API             string
	Namespace       string
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

// SetNamespace sets attribute and returns self
func (t *Config) SetNamespace(namespace string) *Config {
	t.Namespace = namespace
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

// // PrismaClient returns PrismaClient directly
// func (t *Config) PrismaClient(ctx context.Context) (*prisma_api.Client, error) {

// 	tokenprovider, err := NewClient(t)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return prisma_api.NewConfig().
// 		SetNamespace(t.Namespace).
// 		SetTokenProvider(tokenprovider).
// 		SetHTTPClient(t.GetHTTPClient()).Build(ctx)
// }
