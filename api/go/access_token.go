package main

import (
	"fmt"
	"github.com/merci-app/code-samples/api/go/authorization"
)

func main() {
	accessToken := authorization.AccessToken{
		Username: "<USERNAME>",
		Password: "<PASSWORD>",
	}

	resp, err := accessToken.Authenticate()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp)
}
