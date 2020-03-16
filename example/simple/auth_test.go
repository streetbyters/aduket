package auth

import (
	"net/http"
	"testing"

	"github.com/streetbyters/aduket"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	// Create a new aduket server to respond on /login with given ResponseRules
	// and to record incoming request
	server, requestRecorder := aduket.NewServer(
		http.MethodPost, "/login",
		aduket.StatusCode(http.StatusOK),
		aduket.JSONBody(Token{"12345"}),
	)
	defer server.Close()

	auth := AuthClient{authURL: server.URL}

	credentials := Credentials{Username: "streetbyters", Password: "aduket"}
	actualToken, err := auth.Login(credentials)
	assert.Nil(t, err)

	// Assert if request sent with correct body
	requestRecorder.AssertJSONBodyEqual(t, credentials)
	// Assert if request sent with correct headers
	expectedHeader := http.Header{}
	expectedHeader.Add("X-If-You-Read-This", "send-a-hadouken-back")
	requestRecorder.AssertHeaderContains(t, expectedHeader)

	assert.Equal(t, Token{"12345"}, actualToken)
}

func TestUnauthorizedLogin(t *testing.T) {
	// Create a new aduket server to response with 401 on /login
	server, _ := aduket.NewServer(
		http.MethodPost, "/login",
		aduket.StatusCode(http.StatusUnauthorized),
	)
	defer server.Close()

	auth := AuthClient{authURL: server.URL}

	_, err := auth.Login(Credentials{Username: "root", Password: "toor"})
	assert.Equal(t, err, NotAuthorizedError)
}

func TestAuthServerError(t *testing.T) {
	// Create a new aduket server to response with 500 on /login
	server, _ := aduket.NewServer(
		http.MethodPost, "/login",
		aduket.StatusCode(http.StatusInternalServerError),
	)
	defer server.Close()

	auth := AuthClient{authURL: server.URL}

	_, err := auth.Login(Credentials{Username: "streetbyters", Password: "aduket"})
	assert.Equal(t, err, InternalAuthServerError)
}
