package main

import (
	"fmt"

)

func main() {
	accessToken := authorization.NewAuthorization("<USERNAME>", "<PASSWORD>")
	resp, err := accessToken.Authenticate()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp)
}
