/*
This is an example for running on an AWS Lambda. The env var API and NAMESPACE must be
set. AWS will inject AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_SESSION_TOKEN.

This will print the child namespaces of the specified NAMESPACE for the given API

*/
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/go-multierror"

	prisma_api "github.com/aporeto-se/prisma-sdk-go-v2/api"
	token "github.com/aporeto-se/prisma-sdk-go-v2/token/aws/envvars"
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

	httpClient := &http.Client{}

	tokenprovider, err := token.NewConfig().
		SetAPI(api).
		SetAccessKeyID(accessKeyID).
		SetSecretAccessKey(secretAccessKey).
		SetSessionToken(sessionToken).Build()
	if err != nil {
		panic(err)
	}

	token, err := tokenprovider.Token(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(token)

	if namespace == "" {
		fmt.Println(fmt.Sprintf("env var %s is NOT set", NamespaceEnv))
		return
	}

	fmt.Println(fmt.Sprintf("env var %s IS set", NamespaceEnv))

	prismaClient, err := prisma_api.NewConfig().
		SetNamespace(namespace).
		SetAPI(api).
		SetTokenProvider(tokenprovider).
		SetHTTPClient(httpClient).Build(ctx)

	if err != nil {
		panic(err)
	}

	for _, ns := range prismaClient.GetNamespaces() {
		fmt.Println(ns.Name)
	}

}
