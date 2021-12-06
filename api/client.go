package prismasdk2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"

	"github.com/aporeto-se/prisma-sdk-go-v2/types"
)

// Client is this client
type Client struct {
	api           string
	namespacePath string
	TokenProvider
	httpClient *http.Client
	namespace  *types.Namespace
	namespaces []*types.Namespace
	mutex      sync.Mutex
}

// NewClient returns new Client
func NewClient(ctx context.Context, config *Config) (*Client, error) {

	zap.L().Debug("entering NewClient")

	var errors *multierror.Error

	if config.API == "" {
		errors = multierror.Append(errors, fmt.Errorf("attribute API is required"))
	}

	if config.Namespace == "" {
		errors = multierror.Append(errors, fmt.Errorf("attribute Namespace is required"))
	}

	if config.TokenProvider == nil {
		errors = multierror.Append(errors, fmt.Errorf("interface TokenProvider is required"))
	}

	err := errors.ErrorOrNil()
	if err != nil {
		zap.L().Debug("returning NewClient with error(s)")
		return nil, err
	}

	httpClient := config.HTTPClient

	if httpClient == nil {
		httpClient = &http.Client{}
		zap.L().Debug("HTTPClient created new")
	} else {
		zap.L().Debug("HTTPClient set from config")
	}

	client := &Client{
		api:           config.API,
		namespacePath: config.Namespace,
		TokenProvider: config.TokenProvider,
		httpClient:    config.HTTPClient,
		namespace: &types.Namespace{
			Name:          basename(config.Namespace),
			NamespaceType: types.NamespaceTypeUndefined,
		},
	}

	err = client.SyncNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	zap.L().Debug("returning NewClient")
	return client, nil
}

func basename(s string) string {
	sp := strings.Split(s, "/")
	return sp[len(sp)-1]
}

func namespaceResToNamespace(namespace *namespaceRes) *types.Namespace {

	s := strings.Split(namespace.Name, "/")
	name := s[len(s)-1]

	namespaceType, _ := types.NamespaceTypeFromString(namespace.Type)
	ingressTrafficAction, _ := types.TrafficActionFromString(namespace.DefaultPUIncomingTrafficAction)
	egressTrafficAction, _ := types.TrafficActionFromString(namespace.DefaultPUOutgoingTrafficAction)

	return &types.Namespace{
		Name:                           name,
		NamespaceType:                  namespaceType,
		DefaultPUIncomingTrafficAction: ingressTrafficAction,
		DefaultPUOutgoingTrafficAction: egressTrafficAction,
		ID:                             namespace.ID,
		Annotations:                    namespace.Annotations,
	}
}

type namespaceRes struct {
	ID                 string `json:"ID"`
	JWTCertificateType string `json:"JWTCertificateType"`
	JWTCertificates    struct {
	} `json:"JWTCertificates"`
	SSHCAEnabled                   bool                `json:"SSHCAEnabled"`
	Annotations                    map[string][]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	AssociatedSSHCAID              string              `json:"associatedSSHCAID"`
	AssociatedTags                 []string            `json:"associatedTags"`
	CreateTime                     time.Time           `json:"createTime"`
	CustomZoning                   bool                `json:"customZoning"`
	DefaultEnforcerVersion         string              `json:"defaultEnforcerVersion"`
	DefaultPUIncomingTrafficAction string              `json:"defaultPUIncomingTrafficAction"`
	DefaultPUOutgoingTrafficAction string              `json:"defaultPUOutgoingTrafficAction"`
	Description                    string              `json:"description"`
	LocalCAEnabled                 bool                `json:"localCAEnabled"`
	Metadata                       []string            `json:"metadata"`
	Name                           string              `json:"name"`
	Namespace                      string              `json:"namespace"`
	NetworkAccessPolicyTags        []interface{}       `json:"networkAccessPolicyTags"`
	NormalizedTags                 []string            `json:"normalizedTags"`
	OrganizationalMetadata         []string            `json:"organizationalMetadata"`
	Protected                      bool                `json:"protected"`
	ServiceCertificateValidity     string              `json:"serviceCertificateValidity"`
	TagPrefixes                    []interface{}       `json:"tagPrefixes"`
	Type                           string              `json:"type"`
	UpdateTime                     time.Time           `json:"updateTime"`
	Zoning                         int                 `json:"zoning"`
}

// NewClient returns a new client from the existing client. This is because a namespace can have children.
// When a child namespace is to be manipulated you must get the child directly from its parent
func (t *Client) NewClient(ctx context.Context, name string) (*Client, error) {

	zap.L().Debug("entering NewClient")

	namespace, err := t.GetNamespace(name)
	if err != nil {
		zap.L().Debug("returning NewClient with error(s)")
		return nil, err
	}

	client := &Client{
		api:           t.api,
		namespacePath: t.namespacePath + "/" + namespace.Name,
		TokenProvider: t.TokenProvider,
		httpClient:    t.httpClient,
		namespace:     namespace,
	}

	err = client.SyncNamespaces(ctx)
	if err != nil {
		zap.L().Debug("returning NewClient with error(s)")
		return nil, err
	}

	zap.L().Debug("returning NewClient")
	return client, nil
}

// SyncNamespaces fetches the child namespaces of a namespace. It is called automatically when the
// client is created. If a new namespace is created or deleted from the client this client will be
// aware. But if the change is made from outside of this client then we will not be aware. Calling
// this will cause us to (re)sync the namespaces.
func (t *Client) SyncNamespaces(ctx context.Context) error {

	zap.L().Debug("entering SyncNamespaces")

	token, err := t.Token(ctx)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", t.api+"/namespaces", nil)
	if err != nil {
		zap.L().Debug("returning SyncNamespaces with HTTP Request error(s)")
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	req.Header.Add("X-Namespace", t.namespacePath)
	req.Header.Add("Authorization", "Bearer "+token)

	req.Header.Add("X-Fields", "name")
	req.Header.Add("X-Fields", "ID")
	req.Header.Add("X-Fields", "defaultPUIncomingTrafficAction")
	req.Header.Add("X-Fields", "defaultPUOutgoingTrafficAction")
	req.Header.Add("X-Fields", "description")
	req.Header.Add("X-Fields", "annotations")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		zap.L().Debug("returning SyncNamespaces with HTTP Response error(s)")
		return err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Debug("returning SyncNamespaces with IO error(s)")
		return err
	}

	if resp.StatusCode != 200 {
		zap.L().Debug("returning SyncNamespaces with StatusCode error(s)")
		return types.NewAPIError(bytes)
	}

	var raw []*namespaceRes
	err = json.Unmarshal(bytes, &raw)

	if err != nil {
		zap.L().Debug("returning SyncNamespaces with JSON Unmarshal error(s)")
		return err
	}

	var namespaces []*types.Namespace

	for i := range raw {
		e := raw[i]
		namespaces = append(namespaces, namespaceResToNamespace(e))
	}

	t.namespaces = namespaces

	zap.L().Debug(fmt.Sprintf("received %d children for namespace %s", len(t.namespaces), t.namespacePath))

	t.mutex.Lock()
	defer t.mutex.Unlock()

	zap.L().Debug("returning SyncNamespaces")
	return nil
}

type namespaceDataReq struct {
	Group                          string              `json:"group"`
	DefaultPUIncomingTrafficAction string              `json:"defaultPUIncomingTrafficAction"`
	DefaultPUOutgoingTrafficAction string              `json:"defaultPUOutgoingTrafficAction"`
	Name                           string              `json:"name"`
	AssociatedTags                 []string            `yaml:"associatedTags"`
	Annotations                    map[string][]string `json:"annotations" yaml:"annotations"`
}

// CreateNamespace creates a new namespace and returns the new namespace. The primary difference between the
// input namespace and the returned namespace is the presence of the ID.
func (t *Client) CreateNamespace(ctx context.Context, namespace *types.Namespace) (*types.Namespace, error) {

	zap.L().Debug("entering CreateNamespace")

	existingNamespace := t.getNamespace(namespace.Name)
	if existingNamespace != nil {
		zap.L().Debug(fmt.Sprintf("returning CreateNamespace (%s already exist)", namespace.Name))
		return existingNamespace, nil
	}

	zap.L().Debug(fmt.Sprintf("Namespace %s does not exist; will be created", namespace.Name))

	token, err := t.Token(ctx)
	if err != nil {
		return nil, err
	}

	j, _ := json.Marshal(&namespaceDataReq{
		Name:                           namespace.Name,
		Group:                          string(namespace.NamespaceType),
		DefaultPUIncomingTrafficAction: string(namespace.DefaultPUIncomingTrafficAction),
		DefaultPUOutgoingTrafficAction: string(namespace.DefaultPUOutgoingTrafficAction),
		Annotations:                    namespace.Annotations,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", t.api+"/namespaces", bytes.NewBuffer(j))
	if err != nil {
		zap.L().Debug("returning CreateNamespace with HTTP Request error(s)")
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	req.Header.Add("X-Namespace", t.namespacePath)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("X-Fields", "name")
	req.Header.Add("X-Fields", "ID")
	req.Header.Add("X-Fields", "defaultPUIncomingTrafficAction")
	req.Header.Add("X-Fields", "defaultPUOutgoingTrafficAction")
	req.Header.Add("X-Fields", "description")
	req.Header.Add("X-Fields", "annotations")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		zap.L().Debug("returning CreateNamespace with HTTP Response error(s)")
		return nil, err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		zap.L().Debug("returning CreateNamespace with IO error(s)")
		return nil, err
	}

	if resp.StatusCode != 200 {
		zap.L().Debug("returning CreateNamespace with StatusCode error(s)")
		return nil, types.NewAPIError(bytes)
	}

	var raw *namespaceRes
	err = json.Unmarshal(bytes, &raw)

	if err != nil {
		zap.L().Debug("returning CreateNamespace with JSON Unmarshal error(s)")
		return nil, err
	}

	namespace = namespaceResToNamespace(raw)

	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.namespaces = append(t.namespaces, namespace)

	zap.L().Info(fmt.Sprintf("Namespace %s created with ID %s", namespace.Name, namespace.ID))

	zap.L().Debug("returning CreateNamespace")
	return namespace, nil
}

// DeleteNamespace deletes namespace by name. If the namespace is successfully deleted
// a nil error will be returned.
func (t *Client) DeleteNamespace(ctx context.Context, name string) error {

	zap.L().Debug("entering DeleteNamespace")

	token, err := t.Token(ctx)
	if err != nil {
		zap.L().Debug("returning DeleteNamespace with error(s)")
		return err
	}

	namespace, err := t.GetNamespace(name)
	if err != nil {
		zap.L().Debug("returning DeleteNamespace with error(s)")
		return err
	}

	if namespace.ID == "" {
		zap.L().Debug("returning DeleteNamespace with missing ID error")
		return fmt.Errorf("namespace is missing ID")
	}

	//TODO
	req, err := http.NewRequestWithContext(ctx, "DELETE", t.api+"/namespaces/"+namespace.ID, nil)
	if err != nil {
		zap.L().Debug("returning DeleteNamespace with HTTP Request error(s)")
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Namespace", t.namespacePath)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		zap.L().Debug("returning DeleteNamespace with HTTP Response error(s)")
		return err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		zap.L().Debug("returning DeleteNamespace with IO error(s)")
		return err
	}

	switch resp.StatusCode {
	case 200:
		break
	case 204:
		break

	default:
		zap.L().Debug("returning DeleteNamespace with error(s)")
		return types.NewAPIError(bytes)
	}

	zap.L().Info(fmt.Sprintf("Namespace %s deleted", name))

	zap.L().Debug("returning DeleteNamespace")
	return nil
}

// GetNamespacePath returns namespace path
func (t *Client) GetNamespacePath() string {
	return t.namespacePath
}

// GetNamespaces returns the child namespaces of the current client namespace
func (t *Client) GetNamespaces() []*types.Namespace {

	var result []*types.Namespace

	t.mutex.Lock()
	defer t.mutex.Unlock()

	return append(result, t.namespaces...)
}

// GetNamespace returns the child namespace of the current client by name
func (t *Client) GetNamespace(name string) (*types.Namespace, error) {

	zap.L().Debug("entering GetNamespace")

	result := t.getNamespace(name)

	if result != nil {
		zap.L().Debug("returning GetNamespace")
		return result, nil
	}

	zap.L().Debug("returning GetNamespace with error(s)")

	return nil, &types.APIError{
		Code:        404,
		Description: fmt.Sprintf("Namespace %s not found", name),
	}
}

func (t *Client) getNamespace(name string) *types.Namespace {

	t.mutex.Lock()
	defer t.mutex.Unlock()

	for _, v := range t.namespaces {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// HasNamespace returns true if the current client namespace has the child by name
func (t *Client) HasNamespace(name string) bool {

	zap.L().Debug("entering HasNamespace")

	t.mutex.Lock()
	defer t.mutex.Unlock()

	for _, v := range t.namespaces {
		if v.Name == name {
			zap.L().Debug("returning HasNamespace with true")
			return true
		}
	}

	zap.L().Debug("returning HasNamespace with false")
	return false
}

// ImportPrismaConfig imports config into the currently client namespace
func (t *Client) ImportPrismaConfig(ctx context.Context, prismaConfig *types.PrismaConfig) error {

	zap.L().Debug("entering ImportPrismaConfig")

	token, err := t.Token(ctx)
	if err != nil {
		zap.L().Debug("returning ImportPrismaConfig with error(s)")
		return err
	}

	if prismaConfig.Label == "" {
		zap.L().Debug("returning ImportPrismaConfig with error(s)")
		return fmt.Errorf("name is required")
	}

	if len(prismaConfig.Data.Apiauthorizationpolicies) > 0 {
		prismaConfig.Identities = append(prismaConfig.Identities, "apiauthorizationpolicy")
	}

	if len(prismaConfig.Data.Externalnetworks) > 0 {
		prismaConfig.Identities = append(prismaConfig.Identities, "externalnetwork")
	}

	if len(prismaConfig.Data.Networkrulesetpolicies) > 0 {
		prismaConfig.Identities = append(prismaConfig.Identities, "networkrulesetpolicy")
	}

	zap.L().Debug(fmt.Sprintf("ImportPrismaConfig: namespace=%s, label=%s : start", t.namespacePath, prismaConfig.Label))

	p := &types.PrismaConfigOuter{
		Data: prismaConfig,
	}

	j, err := json.Marshal(p)

	if err != nil {
		zap.L().Debug("returning ImportPrismaConfig with JSON Marshal error(s)")
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", t.api+"/import", bytes.NewBuffer(j))
	if err != nil {
		zap.L().Debug("returning ImportPrismaConfig with HTTP Request error(s)")
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Namespace", t.namespacePath)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		zap.L().Debug("returning ImportPrismaConfig with HTTP Response error(s)")
		return err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		zap.L().Debug("returning ImportPrismaConfig with IO error(s)")
		return err
	}

	switch resp.StatusCode {
	case 200:
		break
	case 204:
		break

	default:
		return types.NewAPIError(bytes)
	}

	zap.L().Info(fmt.Sprintf("PrismaConfig imported into namespace %s", t.namespacePath))

	zap.L().Debug("returning ImportPrismaConfig")
	return nil
}

// // AccountID returns Cloud Account ID or error
// func (t *Client) AccountID(ctx context.Context) (string, error) {

// 	result, err := t.TokenProvider.AccountID(ctx)
// 	if err != nil {
// 		zap.L().Debug("returning AccountID with error(s)")
// 		return "", err
// 	}

// 	return result, nil
// }
