package authorization

import (
	"encoding/base64"
	"errors"
	"github.com/merci-app/code-samples/api/go/client"
	"net/http"
	"sync"
	"time"
)

type AccessToken struct {
	Username    string
	Password    string
	Token       string
	ExpireToken time.Time
	Lock        sync.Mutex
}

func (at *AccessToken) Authenticate() (string, error) {
	at.Lock.Lock()
	defer at.Lock.Unlock()

	token, err := at.getApiTokenFromMemory()
	if err != nil {
		t, expires, err := at.getApiTokenFromRequest()
		if err != nil {
			return "", err
		}
		at.setApiToken(t, expires)
		token = t
	}

	return token, nil
}

func (at *AccessToken) getApiTokenFromRequest() (string, time.Time, error) {

	type oauthResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	var response oauthResponse

	url := "https://auth.hml.caradhras.io/oauth2/token?grant_type=client_credentials"
	basicAuth := base64.StdEncoding.EncodeToString([]byte(at.Username + ":" + at.Password))

	req := client.NewClient()
	resp, _, err := req.Post(url).
		Set("Content-Type", "application/x-www-form-urlencoded").
		Set("Authorization", "Basic "+basicAuth).
		Do(&response)

	if err != nil {
		return "", time.Time{}, errors.New("communication failed")
	}
	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, errors.New("communication error")
	}

	timeToExpire := time.Now().Add(time.Duration(response.ExpiresIn-10) * time.Second)

	return response.AccessToken, timeToExpire, nil
}

func (at *AccessToken) getApiTokenFromMemory() (string, error) {
	if time.Now().After(at.ExpireToken) {
		return "", errors.New("expired token")
	}
	return at.Token, nil
}

func (at *AccessToken) setApiToken(token string, time time.Time) {
	at.Token = token
	at.ExpireToken = time
}
