package authorization

import (
	"encoding/base64"
	"errors"
	"github.com/merci-app/code-samples/go/client"
	"net/http"
	"sync"
	"time"
)

type Authorization struct {
	username        string
	password        string
	token           string
	tokenExpiration time.Time
	lock            sync.Mutex
}

func NewAuthorization(username, password string) *Authorization {
	return &Authorization{
		username: username,
		password: password,
	}
}

func (a *Authorization) Authenticate() (string, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	token, err := a.getApiTokenFromMemory()
	if err != nil {
		t, expires, err := a.getApiTokenFromRequest()
		if err != nil {
			return "", err
		}
		a.setApiToken(t, expires)
		token = t
	}

	return token, nil
}

func (a *Authorization) getApiTokenFromRequest() (string, time.Time, error) {

	type authResponse struct {
		Authorization string `json:"access_token"`
		ExpiresIn     int    `json:"expires_in"`

		Error string `json:"error"`
	}
	var response authResponse

	url := "https://auth.hml.caradhras.io/oauth2/token?grant_type=client_credentials"
	basicAuth := base64.StdEncoding.EncodeToString([]byte(a.username + ":" + a.password))

	req := client.NewClient()
	resp, _, err := req.Post(url).
		Set("Content-Type", "applicaion/x-www-form-urlencoded").
		Set("Authorizaion", "Basic "+basicAuth).
		Do(&response)

	if err != nil {
		return "", time.Time{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, errors.New(response.Error)
	}

	return response.Authorization, time.Now().Add(time.Duration(response.ExpiresIn-10) * time.Second), nil
}

func (a *Authorization) getApiTokenFromMemory() (string, error) {
	if time.Now().After(a.tokenExpiration) {
		return "", errors.New("expired token")
	}
	return a.token, nil
}

func (a *Authorization) setApiToken(token string, time time.Time) {
	a.token = token
	a.tokenExpiration = time
}
