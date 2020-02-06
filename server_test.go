package echolizer

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newStringRequest(method, url, body string) *http.Request {
	request, _ := http.NewRequest(method, url, strings.NewReader(body))
	return request
}

func newJSONRequest(method, url string, body interface{}) *http.Request {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")

	return request
}

func newXMLRequest(method, url string, body interface{}) *http.Request {
	requestBody, _ := xml.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/xml")

	return request
}

func newFormRequest(method, url string, form url.Values) *http.Request {
	request, _ := http.NewRequest(method, url, strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return request
}

func TestServerWithResponseJSON(t *testing.T) {
	type UserResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	expectedUserResponse := UserResponse{ID: 123, Name: "kalt"}

	server, _ := NewServer(http.MethodGet, "/user/:id", http.StatusOK, JSONResponse(expectedUserResponse))

	request := newJSONRequest(http.MethodGet, server.URL+"/user/123", http.NoBody)
	res, _ := http.DefaultClient.Do(request)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	body, _ := ioutil.ReadAll(res.Body)

	actualResponse := UserResponse{}
	json.Unmarshal(body, &actualResponse)

	assert.Equal(t, expectedUserResponse, actualResponse)
}

func TestServerWithXMLResponse(t *testing.T) {
	type UserResponse struct {
		Name string `xml:"name"`
	}

	expectedUserResponse := UserResponse{Name: "john"}

	server, _ := NewServer(http.MethodGet, "/user/123", http.StatusOK, XMLResponse(expectedUserResponse))

	request := newXMLRequest(http.MethodGet, server.URL+"/user/123", http.NoBody)
	res, _ := http.DefaultClient.Do(request)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	responseBody, _ := ioutil.ReadAll(res.Body)

	actualXMLResponse := UserResponse{}
	xml.Unmarshal(responseBody, &actualXMLResponse)

	assert.Equal(t, expectedUserResponse, actualXMLResponse)
}

func TestServerWithStringResponse(t *testing.T) {
	expectedStringResponse := "Hello"

	server, _ := NewServer(http.MethodGet, "/hi", http.StatusOK, StringResponse(expectedStringResponse))

	request := newJSONRequest(http.MethodGet, server.URL+"/hi", http.NoBody)
	res, _ := http.DefaultClient.Do(request)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	actualResponseBody, _ := ioutil.ReadAll(res.Body)

	assert.Equal(t, expectedStringResponse, string(actualResponseBody))
}

func TestServerWithByteResponse(t *testing.T) {
	expectedByteResponse := []byte{'S', 'T', 'R', 'E', 'E', 'T', ' ', 'B', 'Y', 'T', 'E', 'R', 'S'}

	server, _ := NewServer(http.MethodGet, "/hi", http.StatusOK, ByteResponse(expectedByteResponse))

	request := newJSONRequest(http.MethodGet, server.URL+"/hi", http.NoBody)
	res, _ := http.DefaultClient.Do(request)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	actualResponseBody, _ := ioutil.ReadAll(res.Body)

	assert.Equal(t, expectedByteResponse, actualResponseBody)
}
