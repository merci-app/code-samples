package main

import (
	"fmt"
	"github.com/merci-app/code-samples/api/go/authorization"
)

var (
	accessToken = authorization.AccessToken{
		Username: "71nrr9g704auojjhii5pe4jel7",
		Password: "53871i0dmq0i08o96kb3csvf40dsn9i8thfd4saur2cpgg8bgmq",
	}
)

func main() {
	resp, err := accessToken.Authenticate()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp, err)
}
