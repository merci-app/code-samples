package request

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const requestTimeout = 30 * time.Second

type Request struct {
	Username    string
	Password    string
	Environment string
	Context     context.Context

	Token       string
	ExpireToken time.Time

	Lock sync.Mutex
}

type RequestInterface interface {
	Authenticate() (*http.Response, []byte, error)
	Post(module BaseUrlModule, ctx context.Context, url string, request, response interface{}) (*http.Response, []byte, error)
	Put(module BaseUrlModule, ctx context.Context, url string, request, response interface{}) (*http.Response, []byte, error)
	PutWithCustomHeader(module BaseUrlModule, ctx context.Context, url string, request, response interface{}, headers map[string]interface{}) (*http.Response, []byte, error)
	Patch(module BaseUrlModule, ctx context.Context, url string, request, response interface{}) (*http.Response, []byte, error)
	PostBinary(module BaseUrlModule, ctx context.Context, url string, contentType string, request *bytes.Buffer, response interface{}) (*http.Response, []byte, error)
	Get(module BaseUrlModule, ctx context.Context, url string, response interface{}, headers ...[]string) (*http.Response, []byte, error)
	PostWithCustomHeader(module BaseUrlModule, ctx context.Context, url string, request, response interface{}, headers map[string]interface{}) (*http.Response, []byte, error)
	Delete(module BaseUrlModule, ctx context.Context, url string, request, response interface{}, headers map[string]interface{}) (*http.Response, []byte, error)
}

type BaseUrlModule string

const (
	BaseUrlModuleApi              BaseUrlModule = "api"
	BaseUrlModuleDataSetup        BaseUrlModule = "datasetup"
	BaseUrlModuleData             BaseUrlModule = "data"
	BaseUrlModuleVoucher          BaseUrlModule = "voucher"
	BaseUrlModulePayments         BaseUrlModule = "payments"
	BaseUrlModulePaymentsSlip     BaseUrlModule = "paymentslip"
	BaseUrlModuleCompanies        BaseUrlModule = "companies"
	BaseUrlModuleRegDocs          BaseUrlModule = "regdocs"
	BaseUrlModuleOneClick         BaseUrlModule = "oneclick"
	BaseUrlModuleAliasBank        BaseUrlModule = "aliasbank"
	BaseUrlModuleBankTransfersIn  BaseUrlModule = "banktransfersin"
	BaseUrlModuleBankTransfersOut BaseUrlModule = "banktransfersout"
	BaseUrlModuleBackOffice       BaseUrlModule = "backoffice"
	BaseUrlModuleCredit           BaseUrlModule = "credit"
	BaseUrlModulePixBaas          BaseUrlModule = "pix-baas"
	BaseUrlModulePix              BaseUrlModule = "pix"
	BaseUrlModuleLimits           BaseUrlModule = "limits"
	BaseUrlModuleIncomeReport     BaseUrlModule = "declarables"

	ProdEnvironment = "prod"
	HmlEnvironment  = "hml"
)

func (r *Request) baseUrl(module BaseUrlModule) string {
	switch r.Environment {
	case ProdEnvironment:
		return fmt.Sprintf("https://%s.caradhras.io", module)
	default:
		return fmt.Sprintf("https://%s.hml.caradhras.io", module)
	}
}

func (r *Request) Authenticate(response interface{}) (*http.Response, []byte, error) {

	url := "https://auth.hml.caradhras.io/oauth2/token?grant_type=client_credentials"
	if r.Environment == ProdEnvironment {
		url = "https://auth.caradhras.io/oauth2/token?grant_type=client_credentials"
	}

	basicAuth := base64.StdEncoding.EncodeToString([]byte(r.Username + ":" + r.Password))
	req := NewClient(r.Context)

	return req.Post(url).
		Timeout(requestTimeout).
		Set("Content-Type", "application/x-www-form-urlencoded").
		Set("Authorization", "Basic "+basicAuth).
		Do(response)
}

func (r *Request) Post(module BaseUrlModule, ctx context.Context, url string, request, response interface{}) (*http.Response, []byte, error) {

	oauth, oauthErr := r.tokenOAuth(ctx)
	if oauthErr != nil {
		return nil, nil, oauthErr
	}

	req := NewClient(ctx)
	uri := fmt.Sprintf("%s%s", r.baseUrl(module), url)

	return req.Post(uri).
		Timeout(requestTimeout).
		Set("Content-Type", "application/json").
		Set("Authorization", oauth).
		Send(request).
		Do(response)
}

func (r *Request) Get(module BaseUrlModule, ctx context.Context, url string, response interface{}, headers ...[]string) (*http.Response, []byte, error) {

	oauth, oauthErr := r.tokenOAuth(ctx)
	if oauthErr != nil {
		return nil, nil, oauthErr
	}

	req := NewClient(ctx)
	uri := fmt.Sprintf("%s%s", r.baseUrl(module), url)
	req.DoNotUseDefaultHeaders()

	for _, header := range headers {
		req.Set(header[0], header[1])
	}

	return req.Get(uri).
		Timeout(requestTimeout).
		Set("Authorization", oauth).
		Do(response)
}

func (r *Request) Put(module BaseUrlModule, ctx context.Context, url string, request, response interface{}) (*http.Response, []byte, error) {

	oauth, oauthErr := r.tokenOAuth(ctx)
	if oauthErr != nil {
		return nil, nil, oauthErr
	}

	req := NewClient(ctx)
	uri := fmt.Sprintf("%s%s", r.baseUrl(module), url)

	return req.Put(uri).
		Timeout(requestTimeout).
		Set("Content-Type", "application/json").
		Set("Authorization", oauth).
		Send(request).
		Do(response)
}

func (r *Request) PutWithCustomHeader(module BaseUrlModule, ctx context.Context, url string, request, response interface{}, headers map[string]interface{}) (*http.Response, []byte, error) {

	oauth, oauthErr := r.tokenOAuth(ctx)
	if oauthErr != nil {
		return nil, nil, oauthErr
	}

	req := NewClient(ctx)
	uri := fmt.Sprintf("%s%s", r.baseUrl(module), url)

	put := req.Put(uri).
		Timeout(requestTimeout).
		Set("Content-Type", "application/json").
		Set("Authorization", oauth)

	for h, v := range headers {
		put.Set(h, fmt.Sprint(v))
	}

	return put.Send(request).Do(response)
}

func (r *Request) Patch(module BaseUrlModule, ctx context.Context, url string, request, response interface{}) (*http.Response, []byte, error) {

	oauth, oauthErr := r.tokenOAuth(ctx)
	if oauthErr != nil {
		return nil, nil, oauthErr
	}

	req := NewClient(ctx)
	uri := fmt.Sprintf("%s%s", r.baseUrl(module), url)

	return req.Patch(uri).
		Timeout(requestTimeout).
		Set("Content-Type", "application/json").
		Set("Authorization", oauth).
		Send(request).
		Do(response)
}

func (r *Request) PostBinary(module BaseUrlModule, ctx context.Context, url, contentType string, body *bytes.Buffer, response interface{}) (*http.Response, []byte, error) {

	oauth, oauthErr := r.tokenOAuth(ctx)
	if oauthErr != nil {
		return nil, nil, oauthErr
	}

	req := NewClient(ctx)
	uri := fmt.Sprintf("%s%s", r.baseUrl(module), url)
	rq, err := req.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		return nil, nil, err
	}

	return req.Request(rq).
		Timeout(requestTimeout).
		Set("Content-Type", contentType).
		Set("Authorization", oauth).
		Do(response)
}

func (r *Request) PostWithCustomHeader(module BaseUrlModule, ctx context.Context, url string, request, response interface{}, headers map[string]interface{}) (*http.Response, []byte, error) {

	oauth, oauthErr := r.tokenOAuth(ctx)
	if oauthErr != nil {
		return nil, nil, oauthErr
	}

	req := NewClient(ctx)
	uri := fmt.Sprintf("%s%s", r.baseUrl(module), url)

	post := req.Post(uri).
		Timeout(requestTimeout).
		Set("Content-Type", "application/json").
		Set("Authorization", oauth)

	for h, v := range headers {
		post.Set(h, fmt.Sprint(v))
	}

	return post.Send(request).Do(response)
}

func (caradhras *Request) Delete(module BaseUrlModule, ctx context.Context, url string, request, response interface{}, headers map[string]interface{}) (*http.Response, []byte, error) {

	oauth, oauthErr := caradhras.tokenOAuth(ctx)
	if oauthErr != nil {
		return nil, nil, oauthErr
	}

	r := NewClient(ctx)
	uri := fmt.Sprintf("%s%s", caradhras.baseUrl(module), url)
	delete := r.Delete(uri).
		Timeout(requestTimeout).
		Set("Content-Type", "application/json").
		Set("Authorization", oauth)

	for h, v := range headers {
		delete.Set(h, fmt.Sprint(v))
	}

	return delete.Send(request).Do(response)
}
