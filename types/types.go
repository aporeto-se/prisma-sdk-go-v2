package types

import (
	"fmt"
)

const (
	// PrismaAPIVersion is the Prisma API Version that this struct implements
	PrismaAPIVersion = 0
)

// Namespace (Microsegmentation) namespaces define logical groups of resources. They have
// hierarchical relationships, like folders in a file system. You can propagate objects from parent
// to child, but never from child to child or child to parent. You should group your resources
// according to their users, which may be inside or outside of your organization. An optimal
// namespace scheme makes it easier to provide a good user experience and control access.
type Namespace struct {
	Name                           string              `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceType                  NamespaceType       `json:"namespaceType,omitempty" yaml:"namespaceType,omitempty"`
	DefaultPUIncomingTrafficAction TrafficAction       `json:"defaultPUIncomingTrafficAction,omitempty" yaml:"defaultPUIncomingTrafficAction,omitempty"`
	DefaultPUOutgoingTrafficAction TrafficAction       `json:"defaultPUOutgoingTrafficAction,omitempty" yaml:"defaultPUOutgoingTrafficAction,omitempty"`
	AssociatedTags                 []string            `json:"associatedTags,omitempty" yaml:"associatedTags,omitempty"`
	ID                             string              `json:"id,omitempty" yaml:"id,omitempty"`
	TagPrefixes                    []string            `json:"tagPrefixes,omitempty" yaml:"tagPrefixes,omitempty"`
	Annotations                    map[string][]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// NewNamespace returns new Namespace with specified name
func NewNamespace(name string) *Namespace {
	return &Namespace{
		Name: name,
	}
}

// SetNamespaceType sets NamespaceType and returns self
func (t *Namespace) SetNamespaceType(v NamespaceType) *Namespace {
	t.NamespaceType = v
	return t
}

// SetDefaultPUIncomingTrafficAction sets DefaultPUIncomingTrafficAction and returns self
func (t *Namespace) SetDefaultPUIncomingTrafficAction(v TrafficAction) *Namespace {
	t.DefaultPUIncomingTrafficAction = v
	return t
}

// SetDefaultPUOutgoingTrafficAction sets DefaultPUOutgoingTrafficAction and returns self
func (t *Namespace) SetDefaultPUOutgoingTrafficAction(v TrafficAction) *Namespace {
	t.DefaultPUOutgoingTrafficAction = v
	return t
}

// SetAssociatedTags sets AssociatedTags and returns self
func (t *Namespace) SetAssociatedTags(v []string) *Namespace {
	t.AssociatedTags = v
	return t
}

// AddAssociatedTags adds attribute and returns self
func (t *Namespace) AddAssociatedTags(v ...string) *Namespace {
	t.AssociatedTags = append(t.AssociatedTags, v...)
	return t
}

// SetAnnotations sets entity and returns self
func (t *Namespace) SetAnnotations(v map[string][]string) *Namespace {
	t.Annotations = v
	return t
}

// SetAnnotation sets entity key with specified value and returns self
func (t *Namespace) SetAnnotation(key string, value []string) *Namespace {

	if t.Annotations == nil {
		t.Annotations = make(map[string][]string)
	}
	t.Annotations[key] = value
	return t
}

// AddAnnotation adds specified value to entity with key and returns self
func (t *Namespace) AddAnnotation(key string, value []string) *Namespace {
	if t.Annotations == nil {
		t.Annotations = make(map[string][]string)
	}
	mapValue := t.Annotations[key]
	mapValue = append(mapValue, value...)
	t.Annotations[key] = mapValue
	return t
}

// SetID sets ID and returns self
func (t *Namespace) SetID(v string) *Namespace {
	t.ID = v
	return t
}

// SetTagPrefixes sets TagPrefixes and returns self
func (t *Namespace) SetTagPrefixes(v []string) *Namespace {
	t.TagPrefixes = v
	return t
}

// AddTagPrefixs adds TagPrefixe and returns self
func (t *Namespace) AddTagPrefixs(v ...string) *Namespace {
	t.TagPrefixes = append(t.TagPrefixes, v...)
	return t
}

// APIAuthorizationPolicy (API authorization) defines the operations a user can perform in a
// namespace: GET, POST, PUT, DELETE, PATCH, and/or HEAD. It is also possible to restrict the user
// to a subset of the APIs in the namespace by setting authorizedIdentities. An API authorization
// always propagates down to all the children of the current namespace.
type APIAuthorizationPolicy struct {
	Name                 string      `json:"name,omitempty" yaml:"name,omitempty"`
	Description          string      `json:"description,omitempty" yaml:"description,omitempty"`
	Protected            bool        `json:"protected" yaml:"protected"`
	Propagate            bool        `json:"propagate" yaml:"propagate"`
	AssociatedTags       []string    `json:"associatedTags,omitempty" yaml:"associatedTags,omitempty"`
	Annotations          interface{} `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Subject              [][]string  `json:"subject,omitempty" yaml:"subject,omitempty"`
	AuthorizedIdentities []string    `json:"authorizedIdentities,omitempty" yaml:"authorizedIdentities,omitempty"`
	AuthorizedNamespace  string      `json:"authorizedNamespace,omitempty" yaml:"authorizedNamespace,omitempty"`
}

// NewAPIAuthorizationPolicy returns a new APIAuthorizationPolicy with specified name
func NewAPIAuthorizationPolicy(name string) *APIAuthorizationPolicy {
	return &APIAuthorizationPolicy{
		Name: name,
	}
}

// SetDescription sets description and returns self
func (t *APIAuthorizationPolicy) SetDescription(description string) *APIAuthorizationPolicy {
	t.Description = description
	return t
}

// SetProtected sets protected and returns self
func (t *APIAuthorizationPolicy) SetProtected(protected bool) *APIAuthorizationPolicy {
	t.Protected = protected
	return t
}

// SetPropagate sets propagate and returns self
func (t *APIAuthorizationPolicy) SetPropagate(propagate bool) *APIAuthorizationPolicy {
	t.Propagate = propagate
	return t
}

// SetAssociatedTags sets associatedTags and returns self
func (t *APIAuthorizationPolicy) SetAssociatedTags(associatedTags []string) *APIAuthorizationPolicy {
	t.AssociatedTags = associatedTags
	return t
}

// AddAssociatedTag adds associatedTag and returns self
func (t *APIAuthorizationPolicy) AddAssociatedTag(associatedTag string) *APIAuthorizationPolicy {
	t.AssociatedTags = append(t.AssociatedTags, associatedTag)
	return t
}

// SetAnnotations sets annotations and returns self
func (t *APIAuthorizationPolicy) SetAnnotations(annotations interface{}) *APIAuthorizationPolicy {
	t.Annotations = annotations
	return t
}

// SetSubject sets subject and returns self
func (t *APIAuthorizationPolicy) SetSubject(subject [][]string) *APIAuthorizationPolicy {
	t.Subject = subject
	return t
}

// AddSubject adds subject and returns self
func (t *APIAuthorizationPolicy) AddSubject(subject ...string) *APIAuthorizationPolicy {
	t.Subject = append(t.Subject, subject)
	return t
}

// SetAuthorizedIdentities sets authorizedIdentities and returns self
func (t *APIAuthorizationPolicy) SetAuthorizedIdentities(authorizedIdentities []string) *APIAuthorizationPolicy {
	t.AuthorizedIdentities = authorizedIdentities
	return t
}

// AddAuthorizedIdentity sets authorizedIdentity and returns self
func (t *APIAuthorizationPolicy) AddAuthorizedIdentity(authorizedIdentity string) *APIAuthorizationPolicy {
	t.AuthorizedIdentities = append(t.AuthorizedIdentities, authorizedIdentity)
	return t
}

// SetAuthorizedNamespace sets authorizedNamespace and returns self
func (t *APIAuthorizationPolicy) SetAuthorizedNamespace(authorizedNamespace string) *APIAuthorizationPolicy {
	t.AuthorizedNamespace = authorizedNamespace
	return t
}

// Externalnetwork (External Network) represents a random network or IP address that is not
// managed by Aporeto. External networks can be used in network policies to allow traffic from or
// to the declared network or IP, using the provided protocol and port (or range of ports). If you
// want to describe the internet (i.e., anywhere), use 0.0.0.0/0 as the address and 1-65000 for the
// ports. You must assign the external network one or more tags. These allow you to reference the
// external network from your network policies.
type Externalnetwork struct {
	Name           string      `json:"name,omitempty" yaml:"name,omitempty"`
	Description    string      `json:"description,omitempty" yaml:"description,omitempty"`
	Protected      bool        `json:"protected" yaml:"protected"`
	Propagate      bool        `json:"propagate" yaml:"propagate"`
	AssociatedTags []string    `json:"associatedTags,omitempty" yaml:"associatedTags,omitempty"`
	Annotations    interface{} `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Entries        []string    `json:"entries,omitempty" yaml:"entries,omitempty"`
}

// NewExternalnetwork returns a new Externalnetwork with specified name
func NewExternalnetwork(name string) *Externalnetwork {
	return &Externalnetwork{
		Name: name,
	}
}

// SetDescription sets description and returns self
func (t *Externalnetwork) SetDescription(description string) *Externalnetwork {
	t.Description = description
	return t
}

// SetProtected sets protected and returns self
func (t *Externalnetwork) SetProtected(protected bool) *Externalnetwork {
	t.Protected = protected
	return t
}

// SetPropagate sets propagate and returns self
func (t *Externalnetwork) SetPropagate(propagate bool) *Externalnetwork {
	t.Propagate = propagate
	return t
}

// SetAssociatedTags sets associatedTags and returns self
func (t *Externalnetwork) SetAssociatedTags(associatedTags []string) *Externalnetwork {
	t.AssociatedTags = associatedTags
	return t
}

// AddAssociatedTag adds associatedTag and returns self
func (t *Externalnetwork) AddAssociatedTag(associatedTag string) *Externalnetwork {
	t.AssociatedTags = append(t.AssociatedTags, associatedTag)
	return t
}

// SetAnnotations sets annotations and returns self
func (t *Externalnetwork) SetAnnotations(annotations interface{}) *Externalnetwork {
	t.Annotations = annotations
	return t
}

// SetEntries sets entries and returns self
func (t *Externalnetwork) SetEntries(entries []string) *Externalnetwork {
	t.Entries = entries
	return t
}

// AddEntry adds entry and returns self
func (t *Externalnetwork) AddEntry(entry string) *Externalnetwork {
	t.Entries = append(t.Entries, entry)
	return t
}

// Networkrulesetpolicy Prisma network rule set policy
type Networkrulesetpolicy struct {
	Description    string      `json:"description,omitempty" yaml:"description,omitempty"`
	Name           string      `json:"name,omitempty" yaml:"name,omitempty"`
	IncomingRules  []*Rule     `json:"incomingRules,omitempty" yaml:"incomingRules,omitempty"`
	OutgoingRules  []*Rule     `json:"outgoingRules,omitempty" yaml:"outgoingRules,omitempty"`
	Propagate      bool        `json:"propagate" yaml:"propagate"`
	Subject        [][]string  `json:"subject,omitempty" yaml:"subject,omitempty"`
	Protected      bool        `json:"protected" yaml:"protected"`
	AssociatedTags []string    `json:"associatedTags,omitempty" yaml:"associatedTags,omitempty"`
	Annotations    interface{} `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// NewNetworkrulesetpolicy returns a new Networkrulesetpolicy with specified name
func NewNetworkrulesetpolicy(name string) *Networkrulesetpolicy {
	return &Networkrulesetpolicy{
		Name: name,
	}
}

// SetDescription sets description and returns self
func (t *Networkrulesetpolicy) SetDescription(description string) *Networkrulesetpolicy {
	t.Description = description
	return t
}

// SetProtected sets protected and returns self
func (t *Networkrulesetpolicy) SetProtected(protected bool) *Networkrulesetpolicy {
	t.Protected = protected
	return t
}

// SetPropagate sets propagate and returns self
func (t *Networkrulesetpolicy) SetPropagate(propagate bool) *Networkrulesetpolicy {
	t.Propagate = propagate
	return t
}

// SetAssociatedTags sets associatedTags and returns self
func (t *Networkrulesetpolicy) SetAssociatedTags(associatedTags []string) *Networkrulesetpolicy {
	t.AssociatedTags = associatedTags
	return t
}

// AddAssociatedTag adds associatedTag and returns self
func (t *Networkrulesetpolicy) AddAssociatedTag(associatedTag string) *Networkrulesetpolicy {
	t.AssociatedTags = append(t.AssociatedTags, associatedTag)
	return t
}

// SetAnnotations sets annotations and returns self
func (t *Networkrulesetpolicy) SetAnnotations(annotations interface{}) *Networkrulesetpolicy {
	t.Annotations = annotations
	return t
}

// SetSubject sets subject and returns self
func (t *Networkrulesetpolicy) SetSubject(subject [][]string) *Networkrulesetpolicy {
	t.Subject = subject
	return t
}

// AddSubject adds subject and returns self
func (t *Networkrulesetpolicy) AddSubject(subject ...string) *Networkrulesetpolicy {
	t.Subject = append(t.Subject, subject)
	return t
}

// SetOutgoingRules sets outgoing rules and returns self
func (t *Networkrulesetpolicy) SetOutgoingRules(outgoingRules []*Rule) *Networkrulesetpolicy {
	t.OutgoingRules = outgoingRules
	return t
}

// AddOutgoingRule adds outgoing rule and returns self
func (t *Networkrulesetpolicy) AddOutgoingRule(rule *Rule) *Networkrulesetpolicy {
	t.OutgoingRules = append(t.OutgoingRules, rule)
	return t
}

// SetIncomingRules sets incoming rules and returns self
func (t *Networkrulesetpolicy) SetIncomingRules(incomingRules []*Rule) *Networkrulesetpolicy {
	t.IncomingRules = incomingRules
	return t
}

// AddIncomingRule adds incoming rule and returns self
func (t *Networkrulesetpolicy) AddIncomingRule(rule *Rule) *Networkrulesetpolicy {
	t.IncomingRules = append(t.IncomingRules, rule)
	return t
}

// Rule allows or denies a matching traffic type. A rule must be added to a Networkrulesetpolicy as an ingress or egress rule.
type Rule struct {
	Action             TrafficAction `json:"action,omitempty" yaml:"action,omitempty"`
	LogsDisabled       bool          `json:"logsDisabled" yaml:"logsDisabled"`
	Object             [][]string    `json:"object,omitempty" yaml:"object,omitempty"`
	ObservationEnabled bool          `json:"observationEnabled" yaml:"observationEnabled"`
	ProtocolPorts      []string      `json:"protocolPorts,omitempty" yaml:"protocolPorts,omitempty"`
}

// NewRule returns new rule
func NewRule() *Rule {
	return &Rule{}
}

// SetAction sets action and returns self
func (t *Rule) SetAction(trafficAction TrafficAction) *Rule {
	t.Action = trafficAction
	return t
}

// SetTrafficActionInherit sets action to TrafficActionAllow and returns self
func (t *Rule) SetTrafficActionInherit() *Rule {
	t.Action = TrafficActionInherit
	return t
}

// SetTrafficActionAllow sets action to TrafficActionAllow and returns self
func (t *Rule) SetTrafficActionAllow() *Rule {
	t.Action = TrafficActionAllow
	return t
}

// SetTrafficActionReject sets action to SetTrafficActionReject and returns self
func (t *Rule) SetTrafficActionReject() *Rule {
	t.Action = TrafficActionReject
	return t
}

// SetLogsDisabled sets logsDisabled and returns self
func (t *Rule) SetLogsDisabled(logsDisabled bool) *Rule {
	t.LogsDisabled = logsDisabled
	return t
}

// SetObject sets object and returns self
func (t *Rule) SetObject(object [][]string) *Rule {
	t.Object = object
	return t
}

// AddObject adds object and returns self
func (t *Rule) AddObject(object ...string) *Rule {
	t.Object = append(t.Object, object)
	return t
}

// SetObservationEnabled sets observationEnabled and returns self
func (t *Rule) SetObservationEnabled(observationEnabled bool) *Rule {
	t.ObservationEnabled = observationEnabled
	return t
}

// SetProtocolPorts sets protocolPorts and returns self
func (t *Rule) SetProtocolPorts(protocolPorts []string) *Rule {
	t.ProtocolPorts = protocolPorts
	return t
}

// AddProtocolPort adds protocol port and returns self
func (t *Rule) AddProtocolPort(protocolPort string) *Rule {
	t.ProtocolPorts = append(t.ProtocolPorts, protocolPort)
	return t
}

// AddTCPProtocolPort adds TCP protocol port and returns self
func (t *Rule) AddTCPProtocolPort(port int) *Rule {
	t.ProtocolPorts = append(t.ProtocolPorts, fmt.Sprintf("tcp/%d", port))
	return t
}

// AddUDPProtocolPort adds TCP protocol port and returns self
func (t *Rule) AddUDPProtocolPort(port int) *Rule {
	t.ProtocolPorts = append(t.ProtocolPorts, fmt.Sprintf("udp/%d", port))
	return t
}

// PrismaConfig represents the base configuration that is applied to the Prisma API via an http
// post call to api/import
type PrismaConfig struct {
	Label      string `json:"label,omitempty" yaml:"label,omitempty"`
	APIVersion int    `json:"apiVersion" yaml:"APIVersion"`
	Data       struct {
		Apiauthorizationpolicies []*APIAuthorizationPolicy `json:"apiauthorizationpolicies" yaml:"apiauthorizationpolicies"`
		Externalnetworks         []*Externalnetwork        `json:"externalnetworks,omitempty" yaml:"externalnetworks,omitempty"`
		Networkrulesetpolicies   []*Networkrulesetpolicy   `json:"networkrulesetpolicies,omitempty" yaml:"networkrulesetpolicies,omitempty"`
	} `json:"data,omitempty" yaml:"data,omitempty"`
	Identities []string `json:"identities,omitempty" yaml:"identities,omitempty"`
}

// NewPrismaConfig returns new PrismaConfig with specified name
func NewPrismaConfig(label string) *PrismaConfig {
	p := &PrismaConfig{
		APIVersion: PrismaAPIVersion,
		Label:      label,
	}
	return p
}

// SetApiauthorizationpolicies sets subject and returns self
func (t *PrismaConfig) SetApiauthorizationpolicies(apiauthorizationpolicies []*APIAuthorizationPolicy) *PrismaConfig {
	t.Data.Apiauthorizationpolicies = apiauthorizationpolicies
	return t
}

// AddApiauthorizationpolicy adds APIAuthorizationPolicy
func (t *PrismaConfig) AddApiauthorizationpolicy(apiAuthorizationPolicy *APIAuthorizationPolicy) *PrismaConfig {
	t.Data.Apiauthorizationpolicies = append(t.Data.Apiauthorizationpolicies, apiAuthorizationPolicy)
	return t
}

// SetExternalnetworks sets subject and returns self
func (t *PrismaConfig) SetExternalnetworks(externalnetworks []*Externalnetwork) *PrismaConfig {
	t.Data.Externalnetworks = externalnetworks
	return t
}

// AddExternalnetwork adds Externalnetwork
func (t *PrismaConfig) AddExternalnetwork(externalnetwork *Externalnetwork) *PrismaConfig {
	t.Data.Externalnetworks = append(t.Data.Externalnetworks, externalnetwork)
	return t
}

// SetNetworkrulesetpolicies sets subject and returns self
func (t *PrismaConfig) SetNetworkrulesetpolicies(networkrulesetpolicies []*Networkrulesetpolicy) *PrismaConfig {
	t.Data.Networkrulesetpolicies = networkrulesetpolicies
	return t
}

// AddNetworkrulesetpolicy adds Networkrulesetpolicy
func (t *PrismaConfig) AddNetworkrulesetpolicy(networkrulesetpolicy *Networkrulesetpolicy) *PrismaConfig {
	t.Data.Networkrulesetpolicies = append(t.Data.Networkrulesetpolicies, networkrulesetpolicy)
	return t
}

// PrismaConfigOuter is the config outer struct
type PrismaConfigOuter struct {
	Data *PrismaConfig `json:"data,omitempty" yaml:"data,omitempty"`
}
