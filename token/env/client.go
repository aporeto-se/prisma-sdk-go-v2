package token

/*

This implements the TokenProvider Interface and provides Prisma tokens using tokens
already set and existing as env vars.


*/

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/aporeto-se/prisma-sdk-go-v2/token/common"
)

// Client the token client
type Client struct {
	tokenString string
	token       *common.OAuthToken
}

// NewClient returns new Client
func NewClient(config *Config) (*Client, error) {

	zap.L().Debug("entering NewClient")

	if config.TokenString != "" {
		zap.L().Debug("token set in config")
		return newClient(config.TokenString)
	}

	zap.L().Debug("token not set in config, will attempt to find in env var")

	tokenString := getTokenStringFromEnv()

	if tokenString != "" {
		zap.L().Debug("using token from env vars")
		return newClient(tokenString)
	}

	zap.L().Debug("token not set in config or env vars")

	return nil, fmt.Errorf("token not set in config or found in env vars")
}

func newClient(tokenString string) (*Client, error) {

	tokenStringSlice := strings.Split(tokenString, ".")
	if len(tokenStringSlice) != 3 {
		return nil, fmt.Errorf("attribute tokenString is invalid")
	}

	tokenBytes, err := base64.RawURLEncoding.DecodeString(tokenStringSlice[1])
	if err != nil {
		return nil, err
	}

	var token *common.OAuthToken
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		return nil, err
	}

	return &Client{
		tokenString: tokenString,
		token:       token,
	}, nil

}

func getTokenStringFromEnv() string {

	tokenString := getTokenStringFromEnvVar(PrismaTokenEnv)
	if tokenString != "" {
		return tokenString
	}

	tokenString = getTokenStringFromEnvVar(ApoctlTokenEnv)
	if tokenString != "" {
		return tokenString
	}

	tokenString = getTokenStringFromEnvVar(EnforcerdTokenEnv)
	if tokenString != "" {
		return tokenString
	}

	zap.L().Debug(fmt.Sprintf("token not found in env var %s, %s, or %s", PrismaTokenEnv, ApoctlTokenEnv, EnforcerdTokenEnv))

	return ""
}

func getTokenStringFromEnvVar(v string) string {
	r := os.Getenv(v)
	if r != "" {
		zap.L().Debug(fmt.Sprintf("got token from env var %s", v))
	}
	return r
}

// AccountID returns Cloud Account ID or error
func (t *Client) AccountID(ctx context.Context) (string, error) {
	return "", fmt.Errorf("this implementation does not support this function")
}

// Token returns the token or an error
func (t *Client) Token(ctx context.Context) (string, error) {

	err := common.TokenExpired(t.token.Exp)

	if err != nil {
		return "", err
	}

	return t.tokenString, nil

}
