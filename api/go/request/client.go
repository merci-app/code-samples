package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	contentTypeHeader = "Content-Type"
	acceptHeader      = "Accept"

	applicationJSON           = "application/json"
	applicationFormUrlEncoded = "application/x-www-form-urlencoded"
	applicationXml            = "application/xml"
	textXml                   = "text/xml"
)

type Client struct {
	httpClient *http.Client
	header     http.Header
	ctx        context.Context
	transport  *http.Transport
	req        *http.Request
	method     string
	url        string
	data       *bytes.Buffer
	params     map[string]string
}

type ProxiedRoundTrip struct {
	beforeRequestFn func(r *http.Request)
	afterRequestFn  func(r *http.Request, response *http.Response, err error)
	proxied         http.RoundTripper
}

func (p *ProxiedRoundTrip) BeforeRequest(fn func(r *http.Request)) {
	p.beforeRequestFn = fn
}

func (p *ProxiedRoundTrip) AfterRequest(fn func(req *http.Request, response *http.Response, err error)) {
	p.afterRequestFn = fn
}

func (p ProxiedRoundTrip) RoundTrip(req *http.Request) (res *http.Response, e error) {

	p.beforeRequestFn(req)

	res, e = p.proxied.RoundTrip(req)

	p.afterRequestFn(req, res, e)

	return
}

func NewClient(ctx context.Context) *Client {
	c := &Client{
		httpClient: &http.Client{},
		header:     http.Header{},
		ctx:        ctx,
	}

	c.header.Add(contentTypeHeader, applicationJSON)
	c.header.Add(acceptHeader, applicationJSON)

	return c
}

func (c *Client) Transport(tr *http.Transport) *Client {
	c.transport = tr
	return c
}

func (c *Client) Timeout(timeout time.Duration) *Client {
	c.httpClient.Timeout = timeout
	return c
}

func (c *Client) Set(param string, value string) *Client {
	c.header.Set(param, value)
	return c
}

func (c *Client) SetParams(params map[string]string) *Client {
	c.params = params
	return c
}

func (c *Client) Request(req *http.Request) *Client {
	c.req = req
	return c
}

func (c *Client) DoNotUseDefaultHeaders() *Client {
	c.header.Del(contentTypeHeader)
	c.header.Del(acceptHeader)
	return c
}

func (c *Client) Send(body interface{}) *Client {

	if body == nil {
		return c
	}

	buf := new(bytes.Buffer)

	switch strings.ToLower(c.header.Get(contentTypeHeader)) {

	case applicationJSON:

		json.NewEncoder(buf).Encode(body)

	case applicationFormUrlEncoded:

		switch body.(type) {
		case map[string]interface{}:
			data := url.Values{}
			for name, value := range body.(map[string]interface{}) {
				data.Set(name, fmt.Sprintf("%v", value))
			}
			buf = bytes.NewBuffer([]byte(data.Encode()))
		case string:
			buf = bytes.NewBuffer([]byte(body.(string)))
		}

	case applicationXml, textXml:

		switch body.(type) {
		case []byte:
			buf = bytes.NewBuffer(body.([]byte))
		default:
			buf = bytes.NewBuffer([]byte(fmt.Sprintf("%v", body)))
		}

	default:

		buf = bytes.NewBuffer([]byte(fmt.Sprintf("%v", body)))
	}

	c.data = buf

	return c
}

func (c *Client) NewRequest(method string, path string, body *bytes.Buffer) (req *http.Request, err error) {
	var content io.Reader
	if body != nil {
		content = body
	}

	req, err = http.NewRequest(method, path, content)
	if err != nil {
		return
	}

	req.URL.RawQuery = req.URL.Query().Encode()

	req.Header = c.header

	if c.ctx != nil {
		req = req.WithContext(c.ctx)
	}

	return
}

func (c *Client) Get(targetUrl string) *Client {
	c.method = http.MethodGet
	c.url = targetUrl
	return c
}

func (c *Client) Post(targetUrl string) *Client {
	c.method = http.MethodPost
	c.url = targetUrl
	return c
}

func (c *Client) Head(targetUrl string) *Client {
	c.method = http.MethodHead
	c.url = targetUrl
	return c
}

func (c *Client) Put(targetUrl string) *Client {
	c.method = http.MethodPut
	c.url = targetUrl
	return c
}

func (c *Client) Delete(targetUrl string) *Client {
	c.method = http.MethodDelete
	c.url = targetUrl
	return c
}

func (c *Client) Patch(targetUrl string) *Client {
	c.method = http.MethodPatch
	c.url = targetUrl
	return c
}

func (c *Client) Options(targetUrl string) *Client {
	c.method = http.MethodOptions
	c.url = targetUrl
	return c
}

func (c *Client) Do(v interface{}) (resp *http.Response, body []byte, err error) {
	return c.doRequest(v)
}

func (c *Client) doRequest(v interface{}) (resp *http.Response, body []byte, err error) {
	if c.transport != nil {
		c.httpClient.Transport = c.transport
	}

	if c.req == nil {
		c.req, err = c.NewRequest(c.method, c.url, c.data)
		if err != nil {
			return
		}
	}

	if c.params != nil {
		q := c.req.URL.Query()

		for k, v := range c.params {
			q.Add(k, v)
		}

		c.req.URL.RawQuery = q.Encode()
	}

	resp, err = c.httpClient.Do(c.req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if v != nil && len(body) > 0 {
		_ = json.NewDecoder(resp.Body).Decode(v)
	}

	return
}
