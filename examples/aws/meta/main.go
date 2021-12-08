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

	prisma_api "github.com/aporeto-se/prisma-sdk-go-v2/api"
	token "github.com/aporeto-se/prisma-sdk-go-v2/token/aws/meta"
)

const (

	// APIEnv enviroment variable
	APIEnv = "API"

	// NamespaceEnv enviroment variable
	NamespaceEnv = "NAMESPACE"
)

func main() {

	ctx := context.Background()

	api := os.Getenv(APIEnv)
	namespace := os.Getenv(NamespaceEnv)

	if api == "" {
		panic(fmt.Errorf("env var %s is required", APIEnv))
	}

	httpClient := &http.Client{}

	tokenprovider, err := token.NewConfig().
		SetAPI(api).
		SetHTTPClient(httpClient).
		Build()

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
