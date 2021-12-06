/*
This is an example for running on an AWS Lambda. The env var API and NAMESPACE must be
set. AWS will inject AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_SESSION_TOKEN.

This will print the child namespaces of the specified NAMESPACE for the given API

*/
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"

	token "github.com/aporeto-se/prisma-sdk-go-v2/token/aws"
)

const (

	// APIEnv enviroment variable
	APIEnv = "API"

	// NamespaceEnv enviroment variable
	NamespaceEnv = "NAMESPACE"

	// AccessKeyIDEnv enviroment variable
	AccessKeyIDEnv = "AWS_ACCESS_KEY_ID"

	// SecretAccessKeyEnv enviroment variable
	SecretAccessKeyEnv = "AWS_SECRET_ACCESS_KEY"

	// SessionTokenEnv enviroment variable
	SessionTokenEnv = "AWS_SESSION_TOKEN"

	// AWSRegionEnv enviroment variable
	AWSRegionEnv = "AWS_REGION"
)

func main() {

	ctx := context.Background()

	var errors *multierror.Error

	api := os.Getenv(APIEnv)
	namespace := os.Getenv(NamespaceEnv)
	accessKeyID := os.Getenv(AccessKeyIDEnv)
	secretAccessKey := os.Getenv(SecretAccessKeyEnv)
	sessionToken := os.Getenv(SessionTokenEnv)

	if api == "" {
		errors = multierror.Append(errors, fmt.Errorf("env var %s is required", APIEnv))
	}

	if namespace == "" {
		errors = multierror.Append(errors, fmt.Errorf("env var %s is required", NamespaceEnv))
	}

	if accessKeyID == "" {
		errors = multierror.Append(errors, fmt.Errorf("env var %s is required", AccessKeyIDEnv))
	}

	if secretAccessKey == "" {
		errors = multierror.Append(errors, fmt.Errorf("env var %s is required", SecretAccessKeyEnv))
	}

	if sessionToken == "" {
		errors = multierror.Append(errors, fmt.Errorf("env var %s is required", SessionTokenEnv))
	}

	err := errors.ErrorOrNil()
	if err != nil {
		panic(err)
	}

	prismaClient, err := token.NewConfig().
		SetAPI(api).
		SetNamespace(namespace).
		SetAccessKeyID(accessKeyID).
		SetSecretAccessKey(secretAccessKey).
		SetSessionToken(sessionToken).PrismaClient(ctx)

	if err != nil {
		panic(err)
	}

	for _, ns := range prismaClient.GetNamespaces() {
		fmt.Println(ns.Name)
	}

}
