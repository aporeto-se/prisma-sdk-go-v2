package main

import (
	"context"
	"fmt"

	token "github.com/aporeto-se/prisma-sdk-go-v2/token/env"
)

func main() {

	err := example1()

	if err != nil {
		panic(err)
	}

}

func example1() error {

	tokenProvider, err := token.NewClient(token.NewConfig())
	if err != nil {
		return err
	}

	token, err := tokenProvider.Token(context.Background())
	if err != nil {
		return err
	}

	fmt.Println(token)

	return nil
}
