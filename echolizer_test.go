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

func TestAssertBodyEqual(t *testing.T) {
	type UserRequest struct {
		Name string `json:"name"`
	}

	server, requestRecorder := NewEcholizer(http.MethodPost, "/user", http.StatusCreated)

	expectedPayload := UserRequest{Name: "noname"}
	request := newJSONRequest(http.MethodPost, server.URL+"/user", expectedPayload)

	res, _ := http.DefaultClient.Do(request)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertBodyEqual(tester, expectedPayload))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertBodyEqual(tester, UserRequest{Name: "lel"}))
	assert.True(t, tester.Failed())

	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestAssertXMLBodyEqual(t *testing.T) {
	type UserRequest struct {
		Name string `xml:"name"`
	}

	server, requestRecorder := NewEcholizer(http.MethodPost, "/user", http.StatusCreated)

	expectedPayload := UserRequest{Name: "noname"}
	request := newXMLRequest(http.MethodPost, server.URL+"/user", expectedPayload)

	res, _ := http.DefaultClient.Do(request)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertBodyEqual(tester, expectedPayload))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertBodyEqual(tester, UserRequest{Name: "lel"}))
	assert.True(t, tester.Failed())

	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestAssertParamEqual(t *testing.T) {
	server, requestRecorder := NewEcholizer(http.MethodGet, "/user/:id", http.StatusOK)

	request := newJSONRequest(http.MethodGet, server.URL+"/user/123", http.NoBody)
	http.DefaultClient.Do(request)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertParamEqual(tester, "id", "123"))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertParamEqual(tester, "id", "321"))
	assert.True(t, tester.Failed())
}

func TestAssertQueryParamEqual(t *testing.T) {
	server, requestRecorder := NewEcholizer(http.MethodGet, "/user", http.StatusOK)

	request := newJSONRequest(http.MethodGet, server.URL+"/user?name=Joe", http.NoBody)
	http.DefaultClient.Do(request)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertQueryParamEqual(tester, "name", []string{"Joe"}))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertQueryParamEqual(tester, "name", []string{"Doe"}))
	assert.True(t, tester.Failed())
}

func TestAssertFormParamEqual(t *testing.T) {
	server, requestRecorder := NewEcholizer(http.MethodPost, "/user", http.StatusCreated)

	form := url.Values{}
	form.Add("name", "Joe")
	request := newFormRequest(http.MethodPost, server.URL+"/user", form)

	http.DefaultClient.Do(request)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertFormParamEqual(tester, "name", []string{"Joe"}))
	assert.False(t, tester.Failed())
}

func TestEcholizerResponse(t *testing.T) {
	type UserResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	expectedResponse := UserResponse{ID: 123, Name: "kalt"}

	server, _ := NewEcholizerWithResponse(http.MethodGet, "/user/:id", http.StatusOK, expectedResponse)

	request := newJSONRequest(http.MethodGet, server.URL+"/user/123", http.NoBody)
	res, _ := http.DefaultClient.Do(request)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	body, _ := ioutil.ReadAll(res.Body)

	actualResponse := UserResponse{}
	json.Unmarshal(body, &actualResponse)

	assert.Equal(t, expectedResponse, actualResponse)
}
