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

func assertHeaderContains(t *testing.T, expectedHeader, actualHeader http.Header) bool {
	for key, value := range expectedHeader {
		actualValue, contains := actualHeader[key]
		if !assert.True(t, contains) {
			return false
		}
		if !assert.Equal(t, value, actualValue) {
			return false
		}
	}
	return true
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Book struct {
	ISBN string `xml:"isbn"`
	Name string `xml:"name"`
}

//91.2
func TestServer(t *testing.T) {
	serverTests := []struct {
		method              string
		route               string
		responseRuleOptions []ResponseRuleOption
		request             *http.Request
		expectedStatusCode  int
		expectedHeader      http.Header
		expectedBody        []byte
	}{
		{
			method: http.MethodGet,
			route:  "/user",
			responseRuleOptions: []ResponseRuleOption{
				StatusCode(http.StatusOK),
				JSONBody(User{ID: 123, Name: "kalt"}),
				Header(http.Header{"Content-Type": []string{"application/json"}}),
			},
			expectedStatusCode: http.StatusOK,
			expectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
			expectedBody:       jsonMarshal(User{ID: 123, Name: "kalt"}),
		},
		{
			method: http.MethodPost,
			route:  "/user",
			responseRuleOptions: []ResponseRuleOption{
				StatusCode(http.StatusCreated),
				XMLBody(Book{ISBN: "223-123", Name: "n0 n4m3"}),
				Header(http.Header{"Content-Type": []string{"application/xml"}}),
			},
			expectedStatusCode: http.StatusCreated,
			expectedHeader:     http.Header{"Content-Type": []string{"application/xml"}},
			expectedBody:       xmlMarshal(Book{ISBN: "223-123", Name: "n0 n4m3"}),
		},
		{
			method: http.MethodGet,
			route:  "/hi",
			responseRuleOptions: []ResponseRuleOption{
				StatusCode(http.StatusOK),
				StringBody("Hello"),
			},
			expectedStatusCode: http.StatusOK,
			expectedHeader:     http.Header{},
			expectedBody:       []byte("Hello"),
		},
		{
			method: http.MethodGet,
			route:  "/community/best",
			responseRuleOptions: []ResponseRuleOption{
				StatusCode(http.StatusTeapot),
				ByteBody([]byte{'S', 'T', 'R', 'E', 'E', 'T', ' ', 'B', 'Y', 'T', 'E', 'R', 'S'}),
			},
			expectedStatusCode: http.StatusTeapot,
			expectedHeader:     http.Header{},
			expectedBody:       []byte{'S', 'T', 'R', 'E', 'E', 'T', ' ', 'B', 'Y', 'T', 'E', 'R', 'S'},
		},
	}

	for _, test := range serverTests {
		server, _ := NewServer(test.method, test.route, test.responseRuleOptions...)
		defer server.Close()

		request, err := http.NewRequest(test.method, server.URL+test.route, http.NoBody)
		assert.Nil(t, err)

		response, err := http.DefaultClient.Do(request)
		assert.Nil(t, err)

		assert.Equal(t, test.expectedStatusCode, response.StatusCode)

		assertHeaderContains(t, test.expectedHeader, response.Header)

		actualBody, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.Equal(t, test.expectedBody, actualBody)
	}
}

func TestMultiRouteServer(t *testing.T) {

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

		assertHeaderContains(t, test.expectedHeader, response.Header)
	}
}
