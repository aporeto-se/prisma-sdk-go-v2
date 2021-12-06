package token

/*

This implements the TokenProvider Interface and provides Prisma tokens using GCP
tokens. This implementation should run within a GCP environment where it can obtain
a GCP Service Account Token.

type TokenProvider interface {
	Token(context.Context) (string, error)
	AccountID(ctx context.Context) (string, error)
}

*/

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/compute/metadata"
	"go.uber.org/zap"

	"github.com/aporeto-se/prisma-sdk-go-v2/token/common"
	prisma_types "github.com/aporeto-se/prisma-sdk-go-v2/types"
)

var (
	identitySuffix = "instance/service-accounts/default/identity?audience=aporeto&format=full"
)

// Client the token client
type Client struct {
	api        string
	httpClient *http.Client

	token *common.PrismaToken
}

// NewClient returns new Client
func NewClient(config *Config) (*Client, error) {

	zap.L().Debug("entering NewClient(config)")

	if config.API == "" {
		return nil, fmt.Errorf("attribute API is required")
	}

	zap.L().Debug("returning NewClient(config)")
	return &Client{
		api:        config.API,
		httpClient: config.GetHTTPClient(),
	}, nil
}

type req struct {
	Realm    string `json:"realm"`
	Validity string `json:"validity"`
	Quota    int    `json:"quota"`
	Metadata struct {
		Token string `json:"token"`
	} `json:"metadata"`
}

func (t *Client) initToken(ctx context.Context) error {

	zap.L().Debug("entering initToken")

	if t.token != nil {
		zap.L().Debug("Token already exist")
		err := common.TokenExpired(t.token.Claims.Exp)
		if err != nil {
			zap.L().Debug("Token is expired, fetching a new one")
		} else {
			zap.L().Debug("returning initToken")
			return nil
		}
	} else {
		zap.L().Debug("Token does not exist; fetching")
	}

	cloudToken, err := metadata.Get(identitySuffix)
	if err != nil {
		zap.L().Debug("returning initToken with error(s)")
		return err
	}

	c := &req{
		Realm:    "GCPIdentityToken",
		Validity: "12h",
	}

	c.Metadata.Token = cloudToken

	jsonReq, _ := json.Marshal(c)

	req, err := http.NewRequestWithContext(ctx, "POST", t.api+"/issue", bytes.NewBuffer(jsonReq))
	if err != nil {
		zap.L().Debug("returning initToken with error(s)")
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		zap.L().Debug("returning initToken with error(s)")
		return err
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		zap.L().Debug("returning initToken with error(s)")
		return err
	}

	if resp.StatusCode != 200 {
		return prisma_types.NewAPIError(respBytes)
	}

	err = json.Unmarshal(respBytes, &t.token)
	if err != nil {
		zap.L().Debug("returning initToken with error(s)")
		return err
	}

	err = common.TokenExpired(t.token.Claims.Exp)
	if err != nil {
		zap.L().Debug("returning initToken with error(s)")
		return err
	}

	zap.L().Debug("returning initToken")
	return nil
}

// Token returns token string or error
func (t *Client) Token(ctx context.Context) (string, error) {

	zap.L().Debug("entering GetToken")

	err := t.initToken(ctx)
	if err != nil {
		zap.L().Debug("returning GetToken with error(s)")
		return "", err
	}

	zap.L().Debug("returning GetToken")
	return t.token.Token, nil
}

// AccountID returns Cloud Account ID or error
func (t *Client) AccountID(ctx context.Context) (string, error) {

	zap.L().Debug("entering AccountID")

	err := t.initToken(ctx)
	if err != nil {
		zap.L().Debug("returning AccountID with error(s)")
		return "", err
	}

	// We use the Data.Projectnumber attribute for GCP
	result := t.token.Claims.Data.Projectnumber

	if result == "" {
		zap.L().Debug("returning AccountID with error(s)")
		return "", fmt.Errorf("unable to get cloud account ID")
	}

	zap.L().Debug("returning AccountID")
	return result, nil
}
