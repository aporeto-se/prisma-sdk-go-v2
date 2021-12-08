package token

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	prisma_types "github.com/aporeto-se/prisma-sdk-go-v2/types"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"

	"github.com/aporeto-se/prisma-sdk-go-v2/token/common"
)

type tokenRequest struct {
	Realm    string `json:"realm"`
	Validity string `json:"validity"`
	Quota    int    `json:"quota"`
	Metadata struct {
		AccessKeyID     string `json:"accessKeyID"`
		SecretAccessKey string `json:"secretAccessKey"`
		Token           string `json:"token"`
	} `json:"metadata"`
}

type awsStsToken struct {
	Code            string    `json:"code,omitempty" yaml:"code"`
	LastUpdated     time.Time `json:"lastUpdated,omitempty" yaml:"lastUpdated"`
	Type            string    `json:"type,omitempty" yaml:"type"`
	AccessKeyID     string    `json:"accessKeyId,omitempty" yaml:"accessKeyId"`
	SecretAccessKey string    `json:"secretAccessKey,omitempty" yaml:"secretAccessKey"`
	Token           string    `json:"token,omitempty" yaml:"token"`
	Expiration      time.Time `json:"expiration,omitempty" yaml:"expiration"`
}

// Client is the Client
type Client struct {
	api        string
	httpClient *http.Client

	token *common.PrismaToken
}

// NewClient returns a new client
func NewClient(config *Config) (*Client, error) {

	zap.L().Debug("entering NewClient")

	var errors *multierror.Error

	if config.API == "" {
		errors = multierror.Append(errors, fmt.Errorf("attribute API is required"))
	}

	err := errors.ErrorOrNil()
	if err != nil {
		zap.L().Debug("returning NewClient with error(s)")
		return nil, err
	}

	zap.L().Debug("returning NewClient")
	return &Client{
		api:        config.API,
		httpClient: config.GetHTTPClient(),
	}, nil
}

// Token returns token string or error
func (t *Client) Token(ctx context.Context) (string, error) {

	zap.L().Debug("entering Token")

	err := t.initToken(ctx)
	if err != nil {
		zap.L().Debug("returning Token with error(s)")
		return "", nil
	}

	zap.L().Debug("returning Token")
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

func (t *Client) initToken(ctx context.Context) error {

	zap.L().Debug("initToken enter")

	if t.token != nil {
		zap.L().Debug("existing token found in cache")
		err := common.TokenExpired(t.token.Claims.Exp)
		if err == nil {
			zap.L().Debug("token in cache is good")
			return nil
		}
		zap.L().Debug("token in cache is expired")
	} else {
		zap.L().Debug("no existing token in cache")
	}

	role, err := t.getAwsRole(ctx)
	if err != nil {
		zap.L().Debug("initToken returning with error(s)")
		return err
	}

	sessionToken, err := t.getAwsSessionToken(ctx)
	if err != nil {
		zap.L().Debug("initToken returning with error(s)")
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "http://169.254.169.254/latest/meta-data/iam/security-credentials/"+role, nil)
	if err != nil {
		zap.L().Debug("initToken returning with error(s)")
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-aws-ec2-metadata-token", sessionToken)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		zap.L().Debug("initToken returning with error(s)")
		return err
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Debug("initToken returning with error(s)")
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf(string(respBytes))
	}

	var awsStsToken *awsStsToken
	err = json.Unmarshal(respBytes, &awsStsToken)
	if err != nil {
		zap.L().Debug("initToken returning with error(s)")
		return err
	}

	c := &tokenRequest{
		Realm:    "AWSSecurityToken",
		Validity: "12h",
	}

	c.Metadata.AccessKeyID = awsStsToken.AccessKeyID
	c.Metadata.SecretAccessKey = awsStsToken.SecretAccessKey
	c.Metadata.Token = awsStsToken.Token

	jsonReq, err := json.Marshal(c)
	if err != nil {
		zap.L().Debug("initToken returning with error(s)")
		return err
	}

	req, err = http.NewRequestWithContext(ctx, "POST", t.api+"/issue", bytes.NewBuffer(jsonReq))
	if err != nil {
		zap.L().Debug("initToken returning with error(s)")
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err = t.httpClient.Do(req)
	if err != nil {
		zap.L().Debug("returning initToken with error(s)")
		return err
	}

	defer resp.Body.Close()

	respBytes, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		zap.L().Debug("initToken returning with error(s)")
		return err
	}

	if resp.StatusCode != 200 {
		zap.L().Debug("returning initToken with error(s)")
		return prisma_types.NewAPIError(respBytes)
	}

	err = json.Unmarshal(respBytes, &t.token)
	if err != nil {
		zap.L().Debug("initToken returning with error(s)")
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

func (t *Client) getAwsSessionToken(ctx context.Context) (string, error) {

	zap.L().Debug("getAwsSessionToken() enter")

	req, err := http.NewRequestWithContext(ctx, "PUT", "http://169.254.169.254/latest/api/token", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/text")
	req.Header.Add("Content-Type", "application/text")
	req.Header.Add("X-aws-ec2-metadata-token-ttl-seconds", "21600")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		zap.L().Error("Retrieving AWS Session Token: failed")
		return "", err
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		zap.L().Error("Retrieving AWS Session Token: failed")
		return "", err
	}

	if resp.StatusCode != 200 {
		zap.L().Error("Retrieving AWS Session Token: failed")
		return "", fmt.Errorf(string(respBytes))
	}

	zap.L().Debug("getAwsSessionToken() return")
	return string(respBytes), nil
}

func (t *Client) getAwsRole(ctx context.Context) (string, error) {

	zap.L().Debug("getAwsRole() enter")

	req, err := http.NewRequestWithContext(ctx, "GET", "http://169.254.169.254/latest/meta-data/iam/security-credentials/", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/text")
	req.Header.Add("Content-Type", "application/text")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf(string(respBytes))
	}

	zap.L().Debug("getAwsRole() return")
	return string(respBytes), nil
}
