package request

import (
	"github.com/merci-app/code-samples/go/authorization"
	"github.com/merci-app/code-samples/go/client"
	"net/http"
	"time"
)

const requestTimeout = 30 * time.Second

type RequestInterface interface {
	Post(url string, request, response interface{}) (*http.Response, []byte, error)
	Get(url string, response interface{}) (*http.Response, []byte, error)
	Put(url string, request, response interface{}) (*http.Response, []byte, error)
	Delete(url string, request, response interface{}) (*http.Response, []byte, error)
}

type Request struct {
	accessToken authorization.Authorization
}

func NewRequest(accessToken authorization.Authorization) *Request {
	return &Request{
		accessToken: accessToken,
	}
}

func (r *Request) Post(url string, request, response interface{}) (*http.Response, []byte, error) {

	token, tokenErr := r.accessToken.Authenticate()
	if tokenErr != nil {
		return nil, nil, tokenErr
	}

	req := client.NewClient()
	return req.Post(url).
		Timeout(requestTimeout).
		Set("Content-Type", "application/json").
		Set("Authorization", token).
		Send(request).
		Do(response)
}

func (r *Request) Get(url string, response interface{}) (*http.Response, []byte, error) {

	token, tokenErr := r.accessToken.Authenticate()
	if tokenErr != nil {
		return nil, nil, tokenErr
	}

	req := client.NewClient()
	req.DoNotUseDefaultHeaders()
	return req.Get(url).
		Timeout(requestTimeout).
		Set("Authorization", token).
		Do(response)
}

func (r *Request) Put(url string, request, response interface{}) (*http.Response, []byte, error) {

	token, tokenErr := r.accessToken.Authenticate()
	if tokenErr != nil {
		return nil, nil, tokenErr
	}

	req := client.NewClient()
	return req.Put(url).
		Timeout(requestTimeout).
		Set("Content-Type", "application/json").
		Set("Authorization", token).
		Send(request).
		Do(response)
}

func (r *Request) Delete(url string, request, response interface{}) (*http.Response, []byte, error) {

	token, tokenErr := r.accessToken.Authenticate()
	if tokenErr != nil {
		return nil, nil, tokenErr
	}

	req := client.NewClient()
	return req.Delete(url).
		Timeout(requestTimeout).
		Set("Content-Type", "application/json").
		Set("Authorization", token).
		Send(request).
		Do(response)
}
