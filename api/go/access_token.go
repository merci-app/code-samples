package main

import (
	"fmt"
	"github.com/merci-app/code-samples/api/go/authorization"
)

var (
	accessToken = authorization.AccessToken{
		Username: "<USERNAME>",
		Password: "<PASSWORD>",
	}
)

func main() {
	resp, err := accessToken.Authenticate()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp, err)
}
