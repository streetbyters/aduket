package aduket

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

func jsonMarshal(j interface{}) []byte {
	m, _ := json.Marshal(j)
	return m
}

func xmlMarshal(x interface{}) []byte {
	m, _ := xml.Marshal(x)
	return m
}

func TestServerWithResponseJSON(t *testing.T) {
	type UserResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	expectedUserResponse := UserResponse{ID: 123, Name: "kalt"}

	server, _ := NewServer(http.MethodGet, "/user/:id", StatusCode(http.StatusOK), JSONBody(expectedUserResponse))

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

	server, _ := NewServer(http.MethodGet, "/user/123", StatusCode(http.StatusOK), XMLBody(expectedUserResponse))

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

	server, _ := NewServer(http.MethodGet, "/hi", StatusCode(http.StatusOK), StringBody(expectedStringResponse))

	request := newJSONRequest(http.MethodGet, server.URL+"/hi", http.NoBody)
	res, _ := http.DefaultClient.Do(request)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	actualResponseBody, _ := ioutil.ReadAll(res.Body)

	assert.Equal(t, expectedStringResponse, string(actualResponseBody))
}

func TestServerWithByteResponse(t *testing.T) {
	expectedByteResponse := []byte{'S', 'T', 'R', 'E', 'E', 'T', ' ', 'B', 'Y', 'T', 'E', 'R', 'S'}

	server, _ := NewServer(http.MethodGet, "/hi", StatusCode(http.StatusOK), ByteBody(expectedByteResponse))

	request := newJSONRequest(http.MethodGet, server.URL+"/hi", http.NoBody)
	res, _ := http.DefaultClient.Do(request)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	actualResponseBody, _ := ioutil.ReadAll(res.Body)

	assert.Equal(t, expectedByteResponse, actualResponseBody)
}

func TestMultiRouteServer(t *testing.T) {
	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type Book struct {
		ISBN string `xml:"isbn"`
		Name string `xml:"name"`
	}

	server, _ := NewMultiRouteServer(map[Route][]ResponseRuleOption{
		{http.MethodGet, "/user"}: {
			StatusCode(http.StatusOK),
			JSONBody(User{ID: 123, Name: "kalt"}),
			Header(http.Header{"Content-Type": []string{"application/json"}}),
		},
		{http.MethodGet, "/book"}: {
			StatusCode(http.StatusTeapot),
			Header(http.Header{"Content-Type": []string{"application/xml"}}),
			XMLBody(Book{ISBN: "9780262510875", Name: "Structure and Interpretation of Computer Programs"}),
		},
	})

	multiRouteServerTests := []struct {
		request            *http.Request
		expectedStatusCode int
		expectedHeader     http.Header
		expectedBody       interface{}
	}{
		{
			request:            newJSONRequest(http.MethodGet, server.URL+"/user", http.NoBody),
			expectedStatusCode: http.StatusOK,
			expectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
			expectedBody:       jsonMarshal(User{ID: 123, Name: "kalt"}),
		},
		{
			request:            newXMLRequest(http.MethodGet, server.URL+"/book", http.NoBody),
			expectedStatusCode: http.StatusTeapot,
			expectedHeader:     http.Header{"Content-Type": []string{"application/xml"}},
			expectedBody:       xmlMarshal(Book{ISBN: "9780262510875", Name: "Structure and Interpretation of Computer Programs"}),
		},
	}

	for _, test := range multiRouteServerTests {
		response, err := http.DefaultClient.Do(test.request)
		assert.Nil(t, err)
		assert.Equal(t, test.expectedStatusCode, response.StatusCode)

		actualBody, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.Equal(t, test.expectedBody, actualBody)

		for key, value := range test.expectedHeader {
			actualValue, contains := response.Header[key]
			assert.True(t, contains)
			assert.Equal(t, value, actualValue)
		}
	}
	// requestRecorders[Route{http.MethodGet, "/book"}]
}
