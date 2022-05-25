package main

import (
	"context"
	"fmt"
	"github.com/merci-app/code-samples/api/go/request"
	"net/http"
)

var (
	req = request.Request{
		Username:    "<USERNAME>",
		Password:    "<PASSWORD>",
		Environment: request.HmlEnvironment,
		Context:     context.Background(),
	}
)

func main() {

	resp, _, err := req.Authenticate(nil)
	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Body)
	}
	fmt.Println(resp)
}
