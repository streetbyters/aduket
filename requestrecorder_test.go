package aduket

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertJSONBodyEqual(t *testing.T) {
	type UserRequest struct {
		Name string `json:"name"`
	}

	server, requestRecorder := NewServer(http.MethodPost, "/user", StatusCode(http.StatusCreated))
	expectedPayload := UserRequest{Name: "noname"}
	request := newJSONRequest(http.MethodPost, server.URL+"/user", expectedPayload)

	res, _ := http.DefaultClient.Do(request)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertJSONBodyEqual(tester, expectedPayload))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertJSONBodyEqual(tester, UserRequest{Name: "lel"}))
	assert.True(t, tester.Failed())

	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestAssertStringBodyEqual(t *testing.T) {

	server, requestRecorder := NewServer(http.MethodPost, "/user", StatusCode(http.StatusCreated))

	expectedPayload := "Hello"
	request := newStringRequest(http.MethodPost, server.URL+"/user", expectedPayload)
	res, _ := http.DefaultClient.Do(request)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertStringBodyEqual(tester, expectedPayload))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertStringBodyEqual(tester, "olleH"))
	assert.True(t, tester.Failed())

	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestAssertXMLBodyEqual(t *testing.T) {
	type UserRequest struct {
		Name string `xml:"name"`
	}

	server, requestRecorder := NewServer(http.MethodPost, "/user", StatusCode(http.StatusCreated))

	expectedPayload := UserRequest{Name: "noname"}
	request := newXMLRequest(http.MethodPost, server.URL+"/user", expectedPayload)
	res, _ := http.DefaultClient.Do(request)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertXMLBodyEqual(tester, expectedPayload))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertXMLBodyEqual(tester, UserRequest{Name: "lel"}))
	assert.True(t, tester.Failed())

	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestAssertParamEqual(t *testing.T) {
	server, requestRecorder := NewServer(http.MethodGet, "/user/:id", StatusCode(http.StatusOK))

	request := newJSONRequest(http.MethodGet, server.URL+"/user/123", http.NoBody)
	http.DefaultClient.Do(request)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertParamEqual(tester, "id", "123"))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertParamEqual(tester, "id", "321"))
	assert.True(t, tester.Failed())
}

func TestAssertQueryParamEqual(t *testing.T) {
	server, requestRecorder := NewServer(http.MethodGet, "/user", StatusCode(http.StatusOK))

	request := newJSONRequest(http.MethodGet, server.URL+"/user?name=Joe", http.NoBody)
	http.DefaultClient.Do(request)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertQueryParamEqual(tester, "name", []string{"Joe"}))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertQueryParamEqual(tester, "name", []string{"Doe"}))
	assert.True(t, tester.Failed())
}

func TestAssertFormParamEqual(t *testing.T) {
	server, requestRecorder := NewServer(http.MethodPost, "/user", StatusCode(http.StatusCreated))

	form := url.Values{}
	form.Add("name", "Joe")
	request := newFormRequest(http.MethodPost, server.URL+"/user", form)

	http.DefaultClient.Do(request)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertFormParamEqual(tester, "name", []string{"Joe"}))
	assert.False(t, tester.Failed())
}

func TestAssertHeaderEqual(t *testing.T) {
	type UserRequest struct {
		Name string `json:"name"`
	}

	server, requestRecorder := NewServer(http.MethodPost, "/user", StatusCode(http.StatusCreated))
	expectedPayload := UserRequest{Name: "noname"}
	request := newJSONRequest(http.MethodPost, server.URL+"/user", expectedPayload)

	request.Header.Add("Test", "123")

	http.DefaultClient.Do(request)

	tester := &testing.T{}

	assert.True(t, requestRecorder.AssertHeaderEqual(tester, http.Header{"Test": []string{"123"}}))
	assert.False(t, tester.Failed())

	assert.False(t, requestRecorder.AssertHeaderEqual(tester, http.Header{"Test": []string{"noo"}}))
	assert.False(t, requestRecorder.AssertHeaderEqual(tester, http.Header{"West": []string{"123"}}))
}
