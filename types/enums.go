package types

import (
	"fmt"
	"strings"
)

// ================================================================================================

// TrafficAction is the action that will be taken for traffic (accept or reject)
type TrafficAction string

const (
	// TrafficActionUndefined is not a type it represents an invalid or undefined type
	TrafficActionUndefined TrafficAction = "Undefined"
	// TrafficActionInherit inherits behaviour of parenet namespace
	TrafficActionInherit TrafficAction = "Inherit"
	// TrafficActionAllow accepts traffic
	TrafficActionAllow TrafficAction = "Allow"
	// TrafficActionReject rejects traffic
	TrafficActionReject TrafficAction = "Reject"
)

// TrafficActionFromString returns TrafficAction Type from string
func TrafficActionFromString(s string) (TrafficAction, error) {

	switch strings.ToUpper(s) {

	case "ALLOW":
		return TrafficActionAllow, nil
	case "REJECT":
		return TrafficActionReject, nil
	}

	return TrafficActionUndefined, fmt.Errorf("String %s is not a valid TrafficAction type", s)
}

// ================================================================================================

// NamespaceType is the type of namespace
type NamespaceType string

const (
	// NamespaceTypeUndefined is an Invalid NamespaceType. It can be used as a placeholder.
	NamespaceTypeUndefined NamespaceType = "Undefined"
	// NamespaceTypeCloudAccount is a Cloud Acccount
	NamespaceTypeCloudAccount NamespaceType = "CloudAccount"
	// NamespaceTypeGroup is either a logical grouping of vms or a Kubernetes cluster
	NamespaceTypeGroup NamespaceType = "Group"
	// NamespaceTypeKubernetes is a Kubernetes namespace. This should not be confused with the cluster.
	// This is a child of a Kubernets cluster.
	NamespaceTypeKubernetes NamespaceType = "Kubernetes"
)

// NamespaceTypeFromString returns type NamespaceType from string. If the string is not a valid type an error will returned.
func NamespaceTypeFromString(s string) (NamespaceType, error) {

	switch strings.ToUpper(s) {

	case "GROUP":
		return NamespaceTypeGroup, nil
	case "KUBERNETES":
		return NamespaceTypeKubernetes, nil
	}

	return NamespaceTypeUndefined, fmt.Errorf("String %s is not a valid NamespaceType type", s)
}
