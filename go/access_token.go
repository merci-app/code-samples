package main

import (
	"fmt"
	"github.com/merci-app/code-samples/go/authorization"
	"os"
)

func main() {
	username := os.Getenv("dock_username")
	password := os.Getenv("dock_password")
	accessToken := authorization.NewAuthorization(username, password)

	resp, err := accessToken.Authenticate()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp)
}
