package request

import (
	"context"
	"errors"
	"net/http"
	"time"
)

func (r *Request) tokenOAuth(ctx context.Context) (string, error) {

	r.Lock.Lock()
	defer r.Lock.Unlock()

	token, err := r.getApiToken()

	// token expired
	if err != nil {
		t, expires, err := r.postCredentials(ctx)

		if err != nil {
			return "", err
		}

		r.setApiToken(t, expires)

		token = t
	}

	return token, nil
}

func (r *Request) postCredentials(ctx context.Context) (string, time.Time, error) {

	type OAuthResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	var credentials OAuthResponse

	resp, _, err := r.Authenticate(&credentials)
	if err != nil {
		return "", time.Time{}, errors.New("communication failed")
	}

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, errors.New("communication error")
	}

	timeToExpire := time.Now().Add(time.Duration(credentials.ExpiresIn-10) * time.Second)

	return credentials.AccessToken, timeToExpire, nil
}

func (r *Request) getApiToken() (string, error) {

	if time.Now().After(r.ExpireToken) {
		return "", errors.New("expired token")
	}

	return r.Token, nil
}

func (r *Request) setApiToken(token string, time time.Time) {

	r.Token = token
	r.ExpireToken = time
}

func (r *Request) expireApiToken() {

	r.Lock.Lock()
	defer r.Lock.Unlock()

	r.Token = ""
	r.ExpireToken = time.Time{}
}
