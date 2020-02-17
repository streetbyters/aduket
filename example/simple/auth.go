package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type AuthClient struct {
	authURL string
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

var NotAuthorizedError = errors.New("User not authorized")
var InternalAuthServerError = errors.New("Something went wrong in those lands")

func (a *AuthClient) Login(credentials Credentials) (Token, error) {
	credentialsJSON, _ := json.Marshal(credentials)

	req, _ := http.NewRequest(http.MethodPost, a.authURL+"/login", bytes.NewReader(credentialsJSON))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-If-You-Read-This", "send-a-hadouken-back")

	res, _ := http.DefaultClient.Do(req)

	if res.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)

		token := Token{}
		json.Unmarshal(body, &token)

		return token, nil
	}

	if res.StatusCode == http.StatusUnauthorized {
		return Token{}, NotAuthorizedError
	}

	return Token{}, InternalAuthServerError
}
