package token

/*
This implements the TokenProvider Interface and provides Prisma tokens using AWS
attributes secrets (AccessKeyID, SecretAccessKey and SessionToken).

This implementation should normally be used for AWS Lambda functions only. This is
because Lambda functions do not permit access to the AWS metadata API but instead
inject metadata into runtime enviornment variables. As these variables are not updated
they are only good for a short time. If this implementation is used for a longer period
these attributes would need to be updated. The general idea is that a Lambda function
is only going to run for 15 minutes max so we only need to provide the same token during
that period.

type TokenProvider interface {
	GetToken(context.Context) (*PrismaToken, error)
	GetAccountID(ctx context.Context) (string, error)
}

*/
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"

	"github.com/aporeto-se/prisma-sdk-go-v2/token/common"
	prisma_types "github.com/aporeto-se/prisma-sdk-go-v2/types"
)

type req struct {
	Realm    string `json:"realm"`
	Validity string `json:"validity"`
	Quota    int    `json:"quota"`
	Metadata struct {
		AccessKeyID     string `json:"accessKeyID"`
		SecretAccessKey string `json:"secretAccessKey"`
		Token           string `json:"token"`
	} `json:"metadata"`
}

// Client is the Client
type Client struct {
	api             string
	accessKeyID     string
	secretAccessKey string
	sessionToken    string
	httpClient      *http.Client

	token *common.PrismaToken
}

// NewClient returns a new client
func NewClient(config *Config) (*Client, error) {

	zap.L().Debug("entering NewClient")

	var errors *multierror.Error

	if config.API == "" {
		errors = multierror.Append(errors, fmt.Errorf("attribute API is required"))
	}

	if config.AccessKeyID == "" {
		errors = multierror.Append(errors, fmt.Errorf("attribute accessKeyID is required"))
	}

	if config.SecretAccessKey == "" {
		errors = multierror.Append(errors, fmt.Errorf("attribute secretAccessKey is required"))
	}

	if config.SessionToken == "" {
		errors = multierror.Append(errors, fmt.Errorf("attribute sessionToken is required"))
	}

	err := errors.ErrorOrNil()
	if err != nil {
		zap.L().Debug("returning NewClient with error(s)")
		return nil, err
	}

	zap.L().Debug("returning NewClient")
	return &Client{
		api:             config.API,
		accessKeyID:     config.AccessKeyID,
		secretAccessKey: config.SecretAccessKey,
		sessionToken:    config.SessionToken,
		httpClient:      config.GetHTTPClient(),
	}, nil
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

	c := &req{
		Realm:    "AWSSecurityToken",
		Validity: "12h",
	}

	c.Metadata.AccessKeyID = t.accessKeyID
	c.Metadata.SecretAccessKey = t.secretAccessKey
	c.Metadata.Token = t.sessionToken

	jsonReq, err := json.Marshal(c)
	if err != nil {
		zap.L().Debug("returning initToken with error(s)")
		return err
	}

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
		return err
	}

	if resp.StatusCode != 200 {
		zap.L().Debug("returning initToken with error(s)")
		return prisma_types.NewAPIError(respBytes)
	}

	err = json.Unmarshal(respBytes, &t.token)
	if err != nil {
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

	// We use the Data.Organization attribute for AWS
	result := t.token.Claims.Data.Organization

	if result == "" {
		zap.L().Debug("returning AccountID with error(s)")
		return "", fmt.Errorf("unable to get cloud account ID")
	}

	zap.L().Debug("returning AccountID")
	return result, nil
}
